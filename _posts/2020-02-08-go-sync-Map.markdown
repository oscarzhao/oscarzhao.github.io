---
layout: post
title: Golang package sync 剖析(四)：sync.Map
date: 2020-02-08 12:44:00 +0800
categories: Golang sync sync.Map
---

# 一、前言

```text
Go语言在设计上对同步（Synchronization，数据同步和线程同步）提供大量的支持，比如 goroutine和channel同步原语，库层面有

- sync：提供基本的同步原语（比如Mutex、RWMutex、Locker）和 工具类（Once、WaitGroup、Cond、Pool、Map）
- sync/atomic：提供变量的原子操作（基于硬件指令 compare-and-swap）
```

-- 引用自[《Golang package sync 剖析(三)： sync.Cond》](2020-01-27-go-sync-Cond.markdown)

上一期中，我们介绍了 Go 语言对条件变量的实现 `sync.Cond`。本期文章我们介绍 `package sync` 下的另一个工具类：`sync.Map`。

我们先看一个场景：

```
面试官：看你用过 Go，大概用过多久，感觉到哪个段位？
小明：用了一年半左右，算是精通
面试官：用过 sync.Map 吗？
小明：没用过，大致了解过
面试官：Go 内置的 map 和 sync.Map 有什么区别？
小明（假装会）：内置的map 不支持并发访问， sync.Map 支持并发访问
面试官：给内置的 map加上读写锁不一样能实现吗？两个有啥区别？
小明：...

小明卒，享年27岁
```

提到 sync.Map，我们首先想到的是 go 内置的 `map[KeyType]ValueType`（简称 map）。内置 map 是基于 hashtable 实现的，具有以下几个特点：

1. Get/Set 的时间复杂度都是 O(1)
2. 支持范型，使用起来非常方便
3. 不支持并发访问

如果要支持并发访问，通常通过加锁实现，示例代码如下：

```golang
type Any interface{}

type ConcurrentMap struct {
  sync.RWMutex
  data map[Any]Any
}

func NewConcurrentMap(capacity int) *ConcurrentMap {
  return &ConcurrentMap{
    data: make(map[Any]Any, capacity),
  }
}

func (m *ConcurrentMap) Get(key Any) (Any, bool) {
  m.RLock()
  defer m.RUnlock()
  val, ok := m.data[key]
  return val, ok
}

func (m *ConcurrentMap) Set(key, val Any) {
  m.Lock()
  defer m.Unlock()
  m.data[key] = val
}
```

如果把加锁的map做成通用类，由于Go不支持范型，要引入interface，可读性变差。所以通常情况下，我们在 package 级别声明两个变量，直接拿来用：

```golang
// 假设需要一个 string-> string 的map
var (
  rwlock sync.RWMutex
  data map[string]string
)
```

那么问题来了：既然自己给 map 加锁很简单，为什么还要造轮子，搞出来一个 sync.Map？

# 二、sync.Map 有啥优势？

我们扒一扒 [godoc](https://golang.org/pkg/sync/#Map)，原文是英文，下面是我的翻译：

```text
Map 类似于 Go内置的 map[interface{}]interface{}，但是支持多协程的并发访问，而不需要额外的锁或同步机制。它的三个成员函数 Load/Store/Delete 的时间复杂度(均摊时间)都是 O(1)。

Map 类型有特定的应用场景。由于内置 map 在类型安全上更好，也更方便对map里的数据进行定制化处理，大多数时候加一把锁就足够了。

Map 类型针对下面两个应用场景做了优化：
1）特定的 entry<key, value> 只被写入一次，但是会被读很多次，比如只增不减的cache；
2）多个 goroutine 的并发读、写、覆盖三种操作的 key 没有交集；
在上面两个应用场景下，相对于内置map+锁防范，Map 里抢占锁的行为会少很多。

Map 类型的默认值是空，可以被立即使用。Map实例在被第一次使用以后，不能被拷贝。
```

前两段话意思是 sync.Map 功能上是一个支持并发访问的 map，但是有特定的应用场景，大多数情况下你用不到。第三段意思是：sync.Map 针对**map粒度有并发读写但key粒度极少并发读写**的场景做了专门优化。

具体点说，在上面的 ConcurrentMap 类里，每次读写都会锁住整个map对象，而 sync.Map 的读写在 cache hit 时是lock-free的，cache miss 时才需要锁住整个map对象。下面是一个例子来说明 cache hit 和 cache miss：

```golang
package main

import (
  "fmt"
  "sync"
)

func main() {
  var cache sync.Map

  cache.Load("key") // cache miss
  cache.Store("key", "value") // cache miss
  cache.Load("key") // cache hit
  cache.Store("key", "value_v1") // cache hit

  // 删除key，先软删除
  // 触发特定条件后会硬删除
  cache.Delete("key")
  cache.Load("key") // cache hit，但读不到值
  cache.Store("key", "value_v2") // cache hit，但仍需上锁
}
```

需要注意的是：虽然 godoc 给的优化场景是**写一次，读多次**，但并不是绝对的。在对一个 key **写多次，读多次**时，sync.Map 也能保证结果的正确性。

**准确来说，一个key的写入是 lock-free的，对它的高频read/update也是lock-free的；一旦触发 delete，就没法利用lock-free的性能优势了。**

顺便提一句 lock-free。多个线程访问同一块内存时，数据同步是必须的，处理的方式无非是管道、锁、原子操作。通常情况下，提到lock-free，第一个想到的 Compare-And-Swap 就是基于原子操作实现的。锁是编程语言基于操作系统的信号量实现的，可能导致线程被休眠和唤醒，而且它包含加锁和解锁两个动作；原子操作是基于硬件指令，不会导致线程休眠，一步到位。当我们说 lock-free 时，是指实现一个数据结构时，用原子操作代替锁保证同步机制。

回顾下刚才的问题：

*既然自己给 map 加锁很简单，为什么还要造轮子，搞出来一个 sync.Map？*

大家心里应该有一些概念了，但重要的事情值得变着法子重复，sync.Map 的应用场景是：

1. 低频 create
2. 高频 read/update
3. 很少/不 delete

*注：create/read/update/delete 是从逻辑上对数据操作的定义。*

# 三、sync.Map 怎么用

通过上面的介绍，大家应该直到什么使用用 sync.Map，什么时候用普通的加锁map了。这个环节介绍下 sync.Map 的使用方法，首先看 它有哪些成员函数：

```golang
// Load 等价于 m[key]
func (m *Map) Load(key interface{}) (value interface{}, ok bool)
// Store 等价于 m[key]=value
func (m *Map) Store(key, value interface{})
// LoadOrStore 读不到key时，会写入 <key,value>；然后返回对应的value
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
// Delete 等价于 delete(m, key)
func (m *Map) Delete(key interface{})
// Range 等价于 for k, v := range m
func (m *Map) Range(f func(key, value interface{}) bool)
```

Map 的用法非常符合直觉，唯二需要说的是 Load 和 Range。

由于 Load 返回的是 interface{}，我们需要进行显式类型转换。如果 key 不存在，Load 会返回 nil，类型转换会 panic。比如，下面这段代码中 val2.(string) 会导致panic：

```golang
package main

import (
  "fmt"
  "sync"
)

func main() {
  var m sync.Map
  m.Store("key", "value")
  val1, ok1 := m.Load("key")
  fmt.Printf("val1=%v, ok1=%v\n", val1, ok1)
  fmt.Printf("type safe val1 = %v\n", val1.(string))
  val2, ok2 := m.Load("key_not_exist")
  fmt.Printf("val2=%v, ok2=%v\n", val2, ok2)
  fmt.Printf("type safe val2 = %v\n", val2.(string))
}
```

关于 **Range**。考虑这个场景：在 Range 执行过程中，如果另一个线程执行了 Store/Delete，执行 Range 的线程能立即感知到变化吗？

要回答这个问题，我们复用下上面的 create/read/update/delete 概念，而不是粗粒度的 Load/Store/Delete 等成员函数，结论是：

1. create：Range 不能感知到新增的 key
2. read: read 操作不会修改 Map，这里不讨论
3. update：Range 能感知到已有 key 对应的 value 的变化
4. delete：Range 能感知到对应 key 的 value 的标记删除，遍历时会跳过

由于Range 的执行过程是 1) 加载索引表 2) 遍历索引表，所以 key 集合在第一步就确定下来了，开始遍历后无法感知到新增的key；对于已有的key，遍历过程中通过原子操作取出对应的value，所以能够感知<遍历开始，当前时刻>这段时间内发生的update和delete。

# 四、sync.Map 怎么实现

![前方高能](../assets/2020-02-08/front_high_energy.jpg)

小明买了张复活卡，网上搜到了一篇标题叫 "Golang package sync 剖析(四)：sync.Map" 的文章，回顾了下面试官的问题，把前三节看完了。看到 lock-free 数据结构，十分好奇，但是急于通关，第四小节没看就去面试了。

```
面试官：看你用过 Go，大概用过多久，感觉到哪个段位？
小明：用了一年半左右，比较熟悉，算不上精通
面试官：用过 sync.Map 吗？
小明：用过，它是基于 lock-free hashtable 实现的，比<内置map, 读写锁>快。对读操作做了很多优化，&%..#$#@..%#%%$@*
面试官：嗯，还不错。实现一个简单的 sync.Map，把 Load/Store 两个函数的核心逻辑写出来就行，可以不支持 Delete
小明：...

小明卒，享年28岁
```

sync.Map 和 sync.WaitGroup 类似，是典型的用起来简单，实现起来超级复杂的数据结构。事实上，sync.Map 的实现要更复杂一些。

我们先看 Map 的内部结构：

```golang
type Map struct {
  // read字段内部是一个 readOnly 对象
  // Load和Store都首先从这里加载数据
  // 如果没有，才去访问 Map.dirty
  read atomic.Value

  // mu 用于控制 dirty 变量的同步
  // 调用Store时，对于新key或已删除的key，都要写到dirty
  // 调用Load时，如果 miss 达到 len(dirty)，把dirty浅拷贝到 Map.read.m
  mu Mutex 
  dirty map[interface{}]*entry

  // misses 用来记录Load里读 Map.read字段时的cache miss数
  misses int
}

// readOnly 是一个只读结构
// m 是一个只读索引表，从 dirty 浅拷贝过来
// m 和 Map.dirty 里 的 *entry 指向同样的内存位置
type readOnly struct {
  m       map[interface{}]*entry
  amended bool // 是否需要同步
}

// entry 本身不存数据
// 而是把数据放到指针 p 里
type entry struct {
  p unsafe.Pointer // *interface{}
}
```

Map 之所以被称为 lock-free，是因为它的两级缓存。

### 一级缓存

一级缓存是 read 字段，每次进行 Load/Store/Delete 时，在命中一级缓存时，执行的操作是：
1. 加载索引表 readOnly.m，这一步是原子操作
2. 通过 key 读取对应的 *entry
3. 读取/更新 entry.p，也是原子操作

*注：调用 Delete 时，如果能命中一级缓存，entry.p 被置为 nil。*

### 二级缓存

二级缓存是 dirty 字段，依赖锁保证同步。只有在一级缓存cache miss 时，key才会穿透到二级缓存。

对于 Load，cache miss 且两级缓存不同步时执行的操作是：
1. 加锁
2. 取数
3. 记录cache miss 如果 misses == len(dirty) 则刷新一级缓存，**清空二级缓存**
4. 返回结果，解锁

对于 Store，cache miss时执行的操作是：
1. 加锁
2. 如果二级缓存cache hit，则写数据
3. 如果二级缓存 cache miss，
  - 如果两级缓存“已同步”(amended==false)，则需要将**一级缓存刷到二级缓存**上，并将一级缓存置为“不同步”(amended=true)，然后把新 entry 写入二级缓存
  - 如果两级缓存不同步(amended==true)，就把新 entry 写入二级缓存
4. 解锁


对于 Delete，cache miss时执行的操作是：
1. 加锁
2. double check 一级缓存
3.1. 如果两级缓存已经同步(amended==false)，则直接将key 从 dirty 删除
3.2. 如果两级缓存不同步(amended==false)，则执行软删除（将entry.p 置为nil）
4. 解锁

一级缓存的读写是原子操作，二级缓存的读写通过锁保证同步。在使用过程中，应当尽量避免一级缓存被穿透。

说完了 sync.Map 的实现机制，我们看看源码。

### Load

Load 的操作相对简单：
1. 从一级缓存 m.read 取数
2. 一级缓存cache miss，从二级缓存 m.dirty 取数
3. 如果 cache miss 积累到阈值，触发一级缓存的更新

```golang
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
  // 从一级缓存 m.read 取数
  read, _ := m.read.Load().(readOnly)
  e, ok := read.m[key]
  if !ok && read.amended {
    // 一级缓存被击穿，且两级缓存不同步
    m.mu.Lock()
    // 加锁后 double check 一级缓存
    // 保证<上次读取read, 加锁>期间新增的key能被感知到
    read, _ = m.read.Load().(readOnly)
    e, ok = read.m[key]
    if !ok && read.amended {
      // 一级缓存未命中&&两级缓存不同步时
      // 从二级缓存m.dirty取数
      e, ok = m.dirty[key]
      // missLocked 记录 cache miss
      // 如果 misses == len(m.dirty)
      // 则刷新一级缓存 m.read, 并清空二级缓存
      m.missLocked()
    }
    m.mu.Unlock()
  }
  // 两级缓存都没有这个数
  if !ok {
    return nil, false
  }
  return e.load()
}
```

从代码我们可以看到，如果两级缓存不同步且都被击穿，Load 的性能会比较差。

### Store

```golang
func (m *Map) Store(key, value interface{}) {
  // 检查一级缓存 m.read ，并更新
  read, _ := m.read.Load().(readOnly)
  if e, ok := read.m[key]; ok && e.tryStore(&value) {
    return
  }

  // 一级缓存未命中，或命中已删除的key时
  // fallback 到二级缓存 m.dirty
  m.mu.Lock()
  // 加锁后 double check 一级缓存
  // 保证<上次读取read, 加锁>期间新增的key能被感知到
  read, _ = m.read.Load().(readOnly)
  if e, ok := read.m[key]; ok {
    if e.unexpungeLocked() {
      // 如果 key 被软删除，更新到二级缓存 m.dirty
      m.dirty[key] = e
    }
    e.storeLocked(&value)
  } else if e, ok := m.dirty[key]; ok {
    // 一级缓存 m.read 被穿透
    // 二级缓存 m.dirty 被命中
    e.storeLocked(&value)
  } else {
    // 一级缓存 m.read 被穿透
    // 二级缓存 m.dirty 也被穿透
    if !read.amended {
      // 两级缓存已同步，说明二级缓存 m.dirty 是 nil
      // 将一级缓存 m.read 的key 全部刷到二级缓存 m.dirty 里
      // 并将缓存同步状态置为"不同步"(amended=true)
      m.dirtyLocked()
      m.read.Store(readOnly{m: read.m, amended: true})
    }
    // 将新的 entry 写入二级缓存
    m.dirty[key] = newEntry(value)
  }
  m.mu.Unlock()
}
```

Store 一个 entry<key, value> 时，如果命中一级缓存，则更新一级缓存，原子操作；
否则只更新二级缓存，此时新增的 key 只存在于二级缓存里，只有 Load 里一级缓存被击穿的次数足够多时，该key 才会被刷到一级缓存里。

### Delete 

```golang
func (m *Map) Delete(key interface{}) {
  // 检查一级缓存 m.read ，并更新
  read, _ := m.read.Load().(readOnly)
  e, ok := read.m[key]
  if !ok && read.amended {
    // 一级缓存未命中，且两级缓存不同步
    // fallback 到二级缓存 m.dirty
    m.mu.Lock()
    // 加锁后 double check 一级缓存
    // 保证<上次读取read, 加锁>期间新增的key能被感知到
    read, _ = m.read.Load().(readOnly)
    e, ok = read.m[key]
    if !ok && read.amended {
      // 一级缓存确实没有，直接从二级缓存删
      delete(m.dirty, key)
    }
    m.mu.Unlock()
  }
  if ok {
    // 命中一级缓存，执行软删除
    e.delete()
  }
}
```

Delete 一个 key 时，如果命中一级缓存，则更新一级缓存，将 entry.p 置为 nil，是原子操作；
如果命中二级缓存，那么直接从索引 m.dirty 删除。

我们稍微细品一品，如果只命中二级缓存，逻辑非常简单。但是如果命中一级缓存，会有两个问题：

1. Load 和 Range 时，如何知道一个key被删除了？
2. 如果被删除的 key 个数占比非常高，会非常浪费内存，如何进行清理？

先说问题1，key 命中一级缓存，进行 Delete 时，我们可以从一级缓存拿到该 key 对应的 *entry，方法是 
**m.read.Load().(readOnly).m[key]**，我们拿到 *entry 以后，通过原子操作更新 entry.p，代码如下：

```golang
// 为了便于理解
// readOnly 和 entry 也放在这儿
type readOnly struct {
  m       map[interface{}]*entry
  amended bool // 是否需要同步
}
type entry struct {
  p unsafe.Pointer // *interface{}
}

// expunged 是一个占位符指针，用来标记entry已删除
var expunged = unsafe.Pointer(new(interface{}))

// 执行软删除
func (e *entry) delete() (hadValue bool) {
  for {
    p := atomic.LoadPointer(&e.p)
    if p == nil || p == expunged {
      return false
    }
    if atomic.CompareAndSwapPointer(&e.p, p, nil) {
      return true
    }
  }
}
```

执行软删除以后，由于二级缓存 `m.dirty[key]` 指向的内存地址和 **m.read.Load().(readOnly).m[key]** 一样，所以也被“偷偷”更新了。

说说问题2，一个 key 被删除以后，数据会删除，但是索引表里还记录着占位符，占位符的比例太高的话也会影响整体性能。上面讲 Store 函数时，我们提到：

```
如果两级缓存都cache miss，且两级缓存“已同步”，则需要将一级缓存刷到二级缓存上，并将一级缓存置为“不同步”(amended=true)，然后把新 entry 写入二级缓存
```

为什么需要将一级缓存刷到二级缓存呢？我们先看看刷的时候做了什么：

```golang
func (m *Map) dirtyLocked() {
  // 两级缓存已同步，所以 m.dirty 肯定是nil
  if m.dirty != nil {
    return
  }

  read, _ := m.read.Load().(readOnly)
  m.dirty = make(map[interface{}]*entry, len(read.m))
  for k, e := range read.m {
    // 读取entry的删除状态
    // 未删除的entry 才拷贝到 m.dirty
    if !e.tryExpungeLocked() {
      m.dirty[k] = e
    }
  }
}
```

检查删除状态的函数 **tryExpungeLocked** 写得非常风骚，这篇文章的细节已经很多了，这里不再列出来，有兴趣的同学可以看源码。

# 五、小结一下

sync.Map 的底层是典型的 lock-free hash table，一级缓存通过原子操作提供性能保障，二级缓存借助锁保障逻辑的正确性和完备性。两级缓存的同步是 sync.Map 实现里最复杂的部分，本文尽可能从全局上说清楚，细节上肯定会有不少疏漏，希望大家能阅读源码，对 sync.Map 有比较透彻的理解。


小明又买了一张复活卡，...

![还能再救一下](../assets/2020-02-08/try_saved_again.jpg)



# References

1. [sync.Map](https://golang.org/pkg/sync/#Map)
2. [The World's Simplest Lock-Free Hash Table](https://preshing.com/20130605/the-worlds-simplest-lock-free-hash-table/)
