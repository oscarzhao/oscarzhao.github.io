---
layout: post
title: Golang package sync 剖析(三)：sync.Cond
date: 2020-01-31 16:55:00 +0800
categories: Golang sync sync.Cond
---

# 一、前言

```text
Go语言在设计上对同步（Synchronization，数据同步和线程同步）提供大量的支持，比如 goroutine和channel同步原语，库层面有

- sync：提供基本的同步原语（比如Mutex、RWMutex、Locker）和 工具类（Once、WaitGroup、Cond、Pool、Map）
- sync/atomic：提供变量的原子操作（基于硬件指令 compare-and-swap）
```

-- 引用自[《Golang package sync 剖析(二)： sync.WaitGroup》](2020-01-27-go-sync-WaitGroup.markdown)

上一期中，我们介绍了如何使用 `sync.WaitGroup` 提高程序的并行度。本期文章我们介绍 `package sync` 下的另一个工具类：`sync.Cond`。

`sync.Cond` 对标 同步原语“条件变量”，它可以阻塞一个，或同时阻塞多个线程，直到另一个线程 1) 修改了条件变量; 2)通知一个（或所有）等待的线程。

*注：Go语言里没有线程，只有更轻量级的协程。本文中，“线程”均代指“协程”（goroutine）。*

相对于 `sync.Once` 和 `sync.WaitGroup`， `sync.Cond` 比较难以理解，使用门槛也很高，在 Google 上搜一下，排名前10结果中有这样几个：

![搜索结果](../assets/2020-01-31/sync_Cond_issue.jpeg)

非常神奇的是：一篇名为 “如何正确使用sync.Cond” 的帖子竟然有 16k 的浏览量！

![use sync.Cond](../assets/2020-01-31/sync_Cond_issue_detail.jpeg)

究竟是条件变量这个概念难以理解，还是 sync.Cond 的设计太反人类，我们一探究竟。

# 二、sync.Cond 怎么用

开篇我们就提到了条件变量的应用场景，我们回顾一下：

```
sync.Cond 对标 同步原语“条件变量”，它可以阻塞一个，或同时阻塞多个线程，直到另一个线程 1) 修改了**共享变量**; 2)通知该**条件变量**。
```

首先，我们把概念搞清楚，条件变量的作用是控制多个线程对一个共享变量的读写。我们有三类主体：

1. 共享变量：条件变量控制多个线程对该变量的读写；
2. 等待线程：被条件变量阻塞的线程，有一个或多个；
3. 更新线程：更新共享变量，并唤起一个或多个等待线程。

其次，我们看看 sync.Cond 的说明书：

```golang
// 创建一个 sync.Cond 对象
func NewCond(l Locker) *Cond

// 阻塞当前线程，并等待条件触发
func (c *Cond) Wait()

// 唤醒所有等待线程
func (c *Cond) Broadcast()

// 唤起一个等待线程
// 没有等待线程也不会报错
func (c *Cond) Signal()
```

大家看完这段代码，脑子里第一个问题大概是：`NewCond` 要一把锁是干嘛用的？为了便于理解，我们以 kubernetes 源码里 FIFO 队列为例，一步一步说 sync.Cond 的用法：

```golang
type FIFO struct {
  // lock 控制对象读写
  lock sync.RWMutex
  // 阻塞Pop操作，Add成功后激活被阻塞线程
  cond sync.Cond
  // items 存储数据
  items map[string]interface{}
  // queue 存储key
  queue []string
  // keyFunc是hash函数
  keyFunc KeyFunc

  // 维护items和queue同步
  populated bool
  initialPopulationCount int

  // 队列状态：是否已经关闭
  closed     bool
  closedLock sync.Mutex
}
```

首先，这是一个 FIFO 队列，问题又来了：go 内置的 channel 不香吗？还真的是不够香。

![真香警告](../assets/2020-01-31/really_xiang_warning.jpeg)

FIFO 具备一些额外的特性：

1. 支持自定义处理函数，并保障每个元素只被处理一次（exactly once）；
2. 支持元素去重，版本更新，并只处理最新版本，而不是每次更新都处理一次；
3. 支持元素删除，删除的元素不进行处理；
4. 支持 list 所有元素。

FIFO 的成员函数有：

```golang
// 从队头取一个元素，没有则会被阻塞
Pop(PopProcessFunc) (interface{}, error)
// 向队尾加一个元素，如果已经存在，则不做任何操作
Add(obj interface{}) error
AddIfNotPresent(interface{}) error
// 更新元素
Update(obj interface{}) error
// 删除元素
Delete(obj interface{}) error
// 关闭队列
Close()
// 读取所有元素
List() []interface{}
// 读取所有 key
ListKeys() []string
// 通过元素读取元素（通过 keyFunc 映射到同样的 key）
Get(obj interface{}) (item interface{}, exists bool, err error)
// 通过key读取元素
GetByKey(key string) (item interface{}, exists bool, err error)
// 用传入的数组替换队列内容
Replace([]interface{}, string) error
// 同步items和queue
Resync() error
// items和queue是否同步
HasSynced() bool
```

回到本文的主题 sync.Cond, 在上面这个例子中

- 一个 FIFO 实例就是一个**共享变量**；
- 调用 Pop 的线程是**等待线程**；
- 调用 Add 的线程是**更新线程**；

`lock sync.RWMutex` 用于控制对共享变量的并发访问，本质上是控制对 `queue` 和 `items` 两个字段的并发访问。
由于条件变量 `cond sync.Cond` 在实现 `Wait` 时，把锁操作也包含进去了，所以初始化时需要传入一个锁变量。在使用时，是这样的：

```golang
// 初始化一个 FIFO
func NewFIFO(keyFunc KeyFunc) *FIFO {
  // lock 和 cond 均是默认值
  f := &FIFO{
    items:   map[string]interface{}{},
    queue:   []string{},
    keyFunc: keyFunc,
  }
  // 将 lock 共享给 cond
  f.cond.L = &f.lock
  return f
}

// Pop 操作
func (f *FIFO) Pop(process PopProcessFunc) (interface{}, error) {
  // 锁住共享变量
  f.lock.Lock()
  defer f.lock.Unlock()
  for {
    for len(f.queue) == 0 {
      // 队列已关闭
      if f.IsClosed() {
        return nil, ErrFIFOClosed
      }
      // 队列为空，等待数据
      f.cond.Wait()
    }
    
    // 此处省略一段代码...
    // 从 items 和 queue 删除元素
  }
}

// Add 操作
func (f *FIFO) Add(obj interface{}) error {
  id, err := f.keyFunc(obj)
  if err != nil {
    return KeyError{obj, err}
  }
  
  // 锁住共享变量
  f.lock.Lock()
  defer f.lock.Unlock()
  
  // 此处省略一段代码 ...
  // 添加元素到 items 和 queue

  // 通知等待线程
  f.cond.Broadcast()
  return nil
}
```

上面的代码中，**等待线程**做的是：
1. 给共享变量加锁
2. 有数据，就返回数据；没有数据就调用 `Wait` 等数据

**更新线程**做的是：
1. 给共享变量加锁
2. 写入数据，调用 Broadcast

看起来很简单，Ok? 但是你品一品，你细品，发现事情没那么简单。

![事情没那么简单](../assets/2020-01-31/zhoudongyu_easy_thing.jpg)

**等待线程** 加锁以后，**更新线程** 要更新**共享变量**，怎么会取到锁呢？

我们先看看官方文档对 Wait 的解释：

```
Wait atomically unlocks c.L and suspends execution of the calling goroutine. After later resuming execution, Wait locks c.L before returning.
```

大概意思是： `Wait` 首先会解锁 c.L，然后阻塞当前的协程；后续协程被 Broadcast/Signal 唤醒以后，再对 c.L 加锁，然后 return。

所以，`cond sync.Cond` 的初始化需要一把锁，并且和 FIFO 实例用同一把锁。

# 三、sync.Cond 实现

如果不考虑 runtime 如何实现阻塞和激活，`sync.Cond` 本身的实现逻辑还是比较简单的。我们看下源码（删减版）：

```golang
type Cond struct {
  noCopy noCopy

  // 共享变量被访问前，必须取到锁 L
  L Locker

  notify  notifyList
  checker copyChecker
}

// Wait 
func (c *Cond) Wait() {
  // 给当前协程分配一张船票
  t := runtime_notifyListAdd(&c.notify)
  // 解锁
  c.L.Unlock()
  // 暂定当前协程的执行，等通知
  runtime_notifyListWait(&c.notify, t)
  // 加锁
  c.L.Lock()
}

// Signal 唤醒被 c 阻塞的一个协程（如果有）
func (c *Cond) Signal() {
  runtime_notifyListNotifyOne(&c.notify)
}

// Broadcast 唤醒所有被 c 阻塞的协程
func (c *Cond) Broadcast() {
  runtime_notifyListNotifyAll(&c.notify)
}
```

这里着重说下 runtime_* 函数的功能：

1. `runtime_notifyListAdd` 将当前线程添加到通知列表，以能够接收通知；
2. `runtime_notifyListWait` 将当前协程休眠，接收到通知以后才会被唤醒；
3. `runtime_notifyListNotifyOne` 发送通知，唤醒 `notify` 列表里一个协程
4. `runtime_notifyListNotifyAll` 发送通知，唤醒 `notify` 列表里所有协程


# 四、总结

`sync.Cond` 是Go语言对条件变量的一个实现方式，但不是唯一的方式。本质上，`sync.Once` 和 channel 也是条件变量的实现。

1. `sync.Once` 里锁和原子操作用于控制**共享变量**的读写；
2. channel 通过 `close(ch)` 可以通知其他协程读取数据；

但 `sync.Once` 和 channel 有一个明显的缺点是：它们都只能保证**第一次**满足条件变量，而 sync.Cond 可以提供持续的保障。

由于 `sync.Cond` 的复杂性（我认为是 godoc 写的太差了），且应用场景相对较少，其出现频次低于 `sync.Once` 和 `sync.WaitGroup`。不过在合适的应用场景出现时，它就会展示出自己的不可替代性。

# References

1. [C++ std::condition_variable](https://en.cppreference.com/w/cpp/thread/condition_variable)
2. [kubernetes FIFO queue](https://github.com/kubernetes/client-go/blob/master/tools/cache/fifo.go)

