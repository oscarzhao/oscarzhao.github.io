---
layout: post
title:  "Go语言并发模型：像Unix Pipe那样使用channel"
date:   2016-08-12 11:59:00 +0800
categories: Golang
---

## 简介

Go语言的并发原语允许开发者以类似于 Unix Pipe 的方式构建数据流水线 (data pipelines)，数据流水线能够高效地利用 I/O和多核 CPU 的优势。

本文要讲的就是一些使用流水线的一些例子，流水线的错误处理也是本文的重点。

## 阅读建议

数据流水线充分利用了多核特性，代码层面是基于 channel 类型 和 go 关键字。

如果你对这两个概念不太了解，建议先阅读之前发布的两篇文章：Go 语言内存模型(上/下)。

## 什么是 "流水线" (pipeline)?

对于"流水线"这个概念，Go语言中并没有正式的定义，它只是很多种并发方式的一种。这里我给出一个非官方的定义：一条流水线是 是由多个阶段组成的，相邻的两个阶段由 channel 进行连接；
每个阶段是由一组在同一个函数中启动的 goroutine 组成。在每个阶段，这些 goroutine 会执行下面三个操作：

1. 通过 inbound channels 从上游接收数据
2. 对接收到的数据执行一些操作，通常会生成新的数据
3. 将新生成的数据通过 outbound channels 发送给下游

除了第一个和最后一个阶段，每个阶段都可以有任意个 inbound 和 outbound channel。
显然，第一个阶段只有 outbound channel，而最后一个阶段只有 inbound channels。
我们通常称第一个阶段为"生产者"或"源头"，称最后一个阶段为"消费者"或"接收者"。

首先，我们通过一个简单的例子来演示这个概念和其中的技巧。后面我们会更出一个真实世界的例子。

## 第一个例子：求平方数

假设我们有一个流水线，它由三个阶段组成。

第一阶段是 gen 函数，它能够将一组整数转换为channel，channel 可以将数字发送出去。
gen 函数首先启动一个 goroutine，该goroutine 发送数字到 channel，当数字发送完时关闭channel。
代码如下：

``` golang
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}
```

第二阶段是 sq 函数，它从 channel 接收一个整数，然后返回 一个channel，返回的channel可以发送 接收到整数的平方。
当它的 inbound channel 关闭，并且把所有数字均发送到下游时，会关闭 outbound channel。代码如下：

``` golang
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}
```

main 函数 用于设置流水线并运行最后一个阶段。最后一个阶段会从第二阶段接收数字，并逐个打印出来，直到来自于上游的 inbound channel关闭。代码如下：

``` golang
func main() {
    // Set up the pipeline.
    c := gen(2, 3)
    out := sq(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}
```

由于 sq 函数的 inbound channel 和 outbound channel 类型一样，所以组合任意个 sq 函数。比如像下面这样使用：

``` golang
func main() {
    // Set up the pipeline and consume the output.
    for n := range sq(sq(gen(2, 3))) {
        fmt.Println(n) // 16 then 81
    }
}
```

如果我们稍微修改一下 gen 函数，便可以模拟 haskell的惰性求值。有兴趣的读者可以自己折腾一下。

## 第二个例子：扇入和扇出

扇出：同一个 channel 可以被多个函数读取数据，直到channel关闭。
这种机制允许将工作负载分发到一组worker，以便更好地并行使用 CPU 和 I/O。

扇入：多个 channel 的数据可以被同一个函数读取和处理，直到所有 channel都关闭。 
A function can read from multiple inputs and proceed until all are closed by multiplexing the input channels onto a single channel that's closed when all the inputs are closed. This is called fan-in.

我们修改一下上个例子中的流水线，这里我们运行两个 sq 实例，它们从同一个 channel 读取数据。
这里我们引入一个新函数 merge 对结果进行"扇入"操作：

``` golang
func main() {
    in := gen(2, 3)

    // 启动两个 sq 实例，即两个goroutines处理 channel "in" 的数据
    c1 := sq(in)
    c2 := sq(in)

    // merge 函数将 channel c1 和 c2 合并到一起，这段代码会消费 merge 的结果
    for n := range merge(c1, c2) {
        fmt.Println(n) // 打印 4 9, 或 9 4
    }
}
```

merge 函数 将多个 channel 转换为一个 channel，它为每一个 inbound channel 启动一个 goroutine，用于将数据
拷贝到 outbound channel。所有 `output` goroutine 被创建以后，merge 启动一个额外的 goroutine， 这个goroutine会
等待所有 inbound channel 上的发送操作结束以后，关闭 outbound channel。

对已经关闭的channel 执行发送操作(<-ch)会导致异常，所以我们必须保证所有的发送操作都在关闭channel之前结束。 
[sync.WaitGroup](http://golang.org/pkg/sync/#WaitGroup "sync.WaitGroup") 提供了一种组织同步的方式。
merge 函数的实现见下面代码 (注意 wg 变量)：

``` golang
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // 为每一个输入channel cs 创建一个 goroutine output
    // output 将数据从 c 拷贝到 out，直到 c 关闭，然后 调用 wg.Done
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // 启动一个 goroutine，用于所有 output goroutine结束时，关闭 out 
    // 该goroutine 必须在 wg.Add 之后启动
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

## Stopping short

在使用流水线函数时，有一个固定的模式：

1. 在一个阶段，当所有发送操作结束以后，关闭 outbound channel
2. 在一个阶段，一直从 inbount channel 接收数据，直到 inbound channel 全部关闭

在这种模式下，每一个接收阶段都可以写成 `range` 循环的方式，
从而保证所有数据都被成功发送到下游后，goroutine能够立即退出。

在现实中，阶段斌不总是接收所有的 inbound 数据。有时候是设计如此：接收者可能只需要数据的一个子集就可以继续执行。
更常见的情况是：由于前一个阶段返回一个错误，导致该阶段提前退出。
这两种情况下，接收者都不应该继续等待后面的值被传送过来。

我们期望的结果是：`当后一个阶段不需要时，前面的阶段能够停止生产数据。`

在我们的例子中，如果一个阶段不能消费所有的 inbound 数据，试图发送这些数据的 goroutine 会永久阻塞。看下面这段代码片段：

``` golang 
    // 只消费 out 的第一个数据
    out := merge(c1, c2)
    fmt.Println(<-out) // 4 or 9
    return
    // 由于我们不再接收 out 的第二个数据
    // 其中一个 goroutine output 将会在发送时被阻塞
}
```

显然这里存在资源泄漏。一方面goroutine 消耗内存和运行时资源，另一方面goroutine 栈中的堆引用会阻止 gc 执行回收操作。 
既然goroutine 不能被回收，那么他们必须自己退出。

我们重新整理一下流水线中的不同阶段，保证在下游阶段接收数据失败时，上游阶段也能够正常退出。
一个方式是使用带有缓冲的管道作为 outbound channel。缓存可以存储固定个数的数据。
如果缓存没有用完，那么发送操作会立即返回。看下面这段代码示例：

``` golang
c := make(chan int, 2) // buffer size 2
c <- 1  // succeeds immediately
c <- 2  // succeeds immediately
c <- 3  // blocks until another goroutine does <-c and receives 1
```

如果在创建 channel 时就知道要发送的值的个数，使用buffer就能够简化代码。
仍然使用求平方数的例子，我们对 gen 函数进行重写。我们将这组整型数拷贝到一个
缓冲 channel中，从而避免创建一个新的 goroutine：

``` golang
func gen(nums ...int) <-chan int {
    out := make(chan int, len(nums))
    for _, n := range nums {
        out <- n
    }
    close(out)
    return out
}
```

回到 流水线中被阻塞的 goroutine，我们考虑让 merge 函数返回一个缓冲管道：

``` golang
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int, 1) // enough space for the unread inputs
    // ... the rest is unchanged ...
```
尽管这种方法解决了这个程序中阻塞 goroutine的问题，但是从长远来看，它并不是好办法。
缓存大小选择为1 是建立在两个前提之上：

1. 我们已经知道 merge 函数有两个 inbound channel
2. 我们已经知道下游阶段会消耗多少个值

这段代码很脆弱。如果我们在传入一个值给 gen 函数，或者下游阶段读取的值变少，goroutine
会再次被阻塞。

为了从根本上解决这个问题，我们需要提供一种机制，让下游阶段能够告知上游发送者停止接收的消息。
下面我们看下这种机制。

## Explicit cancellation
to be continued...

原作者 Sameer Ajmani，翻译 Oscar

下期预告：[Go语言并发模型：像Unix Pipe 那样使用channel](https://blog.golang.org/pipelines "pipelines")

### 相关链接：
原文链接：https://blog.golang.org/pipelines

reflect 包：https://golang.org/pkg/reflect/

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")