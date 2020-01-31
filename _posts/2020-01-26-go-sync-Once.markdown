---
layout: post
title: Golang package sync 剖析(一)：sync.Once
date: 2020-01-26 19:27:00 +0800
categories: golang sync sync.Once
---

## 前言

Go语言在设计上对同步（Synchronization，数据同步和线程同步）提供大量的支持，比如 goroutine和channel同步原语，库层面有

- sync：提供基本的同步原语（比如Mutex、RWMutex、Locker）和 工具类（Once、WaitGroup、Cond、Pool、Map）
- sync/atomic：提供变量的原子操作（基于硬件指令 compare-and-swap）

注意：当我说“类”时，是指 Go 里的 struct（单身狗要有面向“对象”编程的觉悟）。

Go语言里对同步的支持主要有五类应用场景：

1. 资源独占：当多个线程依赖同一份资源（比如数据），需要同时读/写同一个内存地址时，runtime需要保证只有一个修改这份数据，并且保证该修改对其他线程可见。锁和变量的原子操作为此而设计
2. 生产者-消费者：在生产者-消费者模型中，消费者依赖生产者产出数据。 channel（管道） 为此而设计
3. 懒加载：一个资源，当且仅当第一次执行一个操作时，该操作执行过程中其他的同类操作都会被阻塞，直到该操作完成。sync.Once为此而设计
4. fork-join：一个任务首先创建出N个子任务，N个子任务全部执行完成以后，主任务搜集结果，执行后续操作。sync.WaitGroup 为此而设计
5. 条件变量：条件变量是一个同步原语，可以同时阻塞多个线程，直到另一个线程 1) 修改了条件; 2)通知一个（或所有）等待的线程。sync.Cond 为此而设计

注意：这里当我说"线程"时，了解Go的同学可以自动映射到 "goroutine"(协程)。

关于 1和2，通过[官方文档](https://golang.org/pkg/sync/)了解其用法和实现。本系列的主角是 sync 下的工工具类，从 sync.Once 开始。内容分两部分：

1. sync.Once 用法
2. sync.Once 实现

## sync.Once 用法

在多数情况下，sync.Once 被用于控制变量的初始化，这个变量的读写通常遵循单例模式，满足这三个条件：
1. 当且仅当第一次读某个变量时，进行初始化（写操作）
2. 变量被初始化过程中，所有读都被阻塞（读操作；当变量初始化完成后，读操作继续进行
3. 变量仅初始化一次，初始化完成后驻留在内存里

在 net 库里，系统的网络配置就是存放在一个变量里，代码如下：

```golang
package net

var (
  // guards init of confVal via initConfVal
	confOnce sync.Once 
	confVal  = &conf{goos: runtime.GOOS}
)

// systemConf returns the machine's network configuration.
func systemConf() *conf {
	confOnce.Do(initConfVal)
	return confVal
}

func initConfVal() {
	dnsMode, debugLevel := goDebugNetDNS()
  confVal.dnsDebugLevel = debugLevel
  // 省略部分代码...
}
```

上面这段代码里，`confVal` 存放数据， `confOnce` 控制读写，两个都是 package-level 单例变量。由于 Go 里变量被初始化为默认值，`confOnce` 可以被立即使用，我们重点关注`confOnce.Do`。首先看成员函数 `Do` 的定义：

```golang
func (o *Once) Do(f func())
```

`Do` 接收一个函数作为参数，该函数不接受任务参数，不返回任何参数。具体做什么由使用方决定，错误处理也由使用方控制。

`once.Sync` 可用于任何符合 "exactly once" 语义的场景，比如：

1. 初始化 rpc/http client
2. open/close 文件
3. close channel
4. 线程池初始化

Go语言中，文件被重复关闭会报error，而 channel 被重复关闭报 panic，`once.Sync` 可以保证这类事情不发生，但是不能保证其他业务层面的错误。下面这个例子给出了一种错误处理的方式，供大家参考：

```golang
// source: os/exec/exec.go
package exec

type closeOnce struct {
	*os.File

	once sync.Once
	err  error
}

func (c *closeOnce) Close() error {
	c.once.Do(c.close)
	return c.err
}

func (c *closeOnce) close() {
	c.err = c.File.Close()
}

```

## sync.Once 实现

sync.Once 类通过一个锁变量和原子变量保障 `exactly once`语义，直接撸下源码（为了便于阅读，做了简化处理）：

```golang
package sync

import "sync/atomic"

type Once struct {
	done uint32
	m    Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
      defer atomic.StoreUint32(&o.done, 1)
      f()
    }
	}
}
```

这里 `done` 是一个状态位，用于判断变量是否初始化完成，其有效值是：

- 0: 函数 f 尚未执行或执行中，Once对象创建时 `done`默认值就是0
- 1: 函数 f 已经执行结束，保证 `f` 不会被再次执行

而 `m Mutex` 用于控制临界区的进入，保证同一时间点最多有一个 `f`在执行。

`done` 在 `m.Lock()` 前后的两次校验都是必要的。


## 发散一下

在 Scala 里，有一个关键词 `lazy`，实现了 sync.Once 同样的功能。具体实现上，早期版本使用了 volatile 修饰状态变量 `done`，使用 `synchronized` 替代 `m Mutex`；后来，也改成了基于CAS的方式。

使用体验上，显然 `lazy` 更香！

## References

1. [Golang: sync.Once](https://golang.org/pkg/sync/#Once)
2. [Synchronization(Computer Science)](https://en.wikipedia.org/wiki/Synchronization_(computer_science))
3. [SIP-20 - Improved Lazy Vals Initialization](http://scalajp.github.io/scala.github.com/sips/pending/improved-lazy-val-initialization.html)
