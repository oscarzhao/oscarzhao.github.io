



关于 sync.RWMutex，`Lock()` 优先级高于 `RLock()`，简单介绍下几个方法的功能。

`Lock()` 会：
1. 获取锁，阻塞其他将要调用 `Lock()` 的协程；
2. 阻塞**将要**调用 `RLock()` 的协程（将 readerCount 设置为负数）；
3. 等待已经调用 `RLock()` 的协程，直到所有 `RUnlock()` 被调用；

`Unlock()` 会：
1. 释放被阻塞的 `RLock()`（恢复 readerCount 为正数）;
2. 释放其他的 `Lock()`，如果有；

`RLock()` 会：
1. 增加 readerCount
2. 如果有writer，则等待 writer 写入完成

`RUnlock()` 会：
1. 降低 readerCount
2. 如果有 writer 在等待 (readerCount < 0)
  - 2.a 如果当前 reader 是最后一个 reader，则释放信号量
  - 2.b 如果当前 reader 不是最后一个reader，则只是降低 readerWait

