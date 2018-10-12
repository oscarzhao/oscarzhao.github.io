---
layout: post
title:  "用Go写算法：求最小可用自然数"
date:   2018-10-12 11:59:00 +0800
categories: Golang go algorithm
---

# 前言

前一段时间在 reddit 上看到有人推广一篇名为 [GopherCon 2018 - Demystifying Binary Search Tree Algorithms](https://about.sourcegraph.com/go/gophercon-2018-binary-search-tree-algorithms/) 的博客，
博客中列举了传统大学里学习算法的种种弊端，并强调了用 Go 实现算法是多么简单有趣，然后拿二叉树举了个例子。读完这篇博客以后，我不得不说，真心没看出来 Go 写算法的优势在哪里。但是，配图确实萌翻了，下面盗图一副。

![binary search tree](http://oat5ddzns.bkt.clouddn.com/gophers-binary-search-tree.png)

虽然不会画画，但是并不妨碍用 Go 做一些算法实现的尝试。这里我从 "Pears of Functional Algorithms Design" 里拿了一道题：给定一个无序自然数数组 A，求出不在 A 中的最小自然数，约束条件如下：

1. A 中的元素个数是有限的，每个元素都是自然数，并且互不相同（自然数包含 0 和 正整数）。
2. 假设 A 中元素的大小在 [0, inf] 之间。

```{text}
Note: 由于计算机本身的限制，我们这里加上 inf = 2^31-1
```

举几个例子：

1. A = [1,2,9,4], output = 0
2. A = [0,1,2], output = 3
3. A = [0,1,16,8,4,2], output = 3
4. A = [], output = 0

首先我们声明一个函数来表述这个问题：

```{go}
func f(A []int) int
```

然后，我们思考如何解决这个问题。

此处建议思考两分钟，如果身边有笔有纸，建议写写画画。

# 方案一

首先映入脑海的方法是：声明一个 bool 数组 B，将 A 中元素作为下标，将 B 对应的元素置为 true，遍历 B，返回第一个 false 的下标。这个方法显然不靠谱，在某些条件下，它会使用大量内存。为了更直观，这里把具体实现贴出来：

```{go}
func f(A []int) int {
  maxElem := max(A) // 函数 max 用于获取数组的最大值
  B := make([]bool, maxElem+1, maxElem+1)
  for _, elem := range A {
    B[elem] = true
  }
  for idx, val := range B {
    if !val {
      return idx
    }
  }
}
```

在极限情况下，例如 A = [2^31-1]，上面代码中 B 会占用 (2^31-1)/8 ≈ 2^28 byte = 256 Mbyte 内存，显然很不合理。

# 方案二：引入排序函数

如果引入一个库函数呢，比如 `sort`。先用 `sort` 给 A 排序，然后从 0 开始递增，找到第一个不在 A 中的元素。这个方法可以解决问题，假设 n=len(A)，排序时间 O(nlogn), 检索效率 O(n)。我们看一下代码实现：

```{go}
func f(A []int) int {
  sort.Ints(A)
  for idx, val := range A {
    if idx != val {
      return idx
    }
  }
  return len(A)
}
```

这是一个很常规的方法。这里的实现有两个地方需要注意：

1. `if idx != val` 避免了创建额外一个变量和对其进行管理的成本
2. `return len(A)` 优雅地处理了边界情况

但是引入 `sort` 的代价也很明显：O(nlogn) 的时间复杂度比较高。有没有时间复杂度为 O(n) 的方法呢？

我们再次思考这个问题本身。对于乱序状态下的 A，对于 [0, inf] 的每一个自然数，搜索的耗时都是 O(n)，最多搜索 n 个自然数。要达到总体复杂度为 O(n)，一个方法是将搜每个自然数搜索的耗时降低到 O(1)，另一个方法是 只搜索 O(1) 个自然数。后一个方法看起来有点难以实现，但是前一个方法只需要我们引入一个 hashmap。

而 Go 语言内置的 map 就是基于 hashmap。

# 方案三：引入 map

引入 map 以后，大致步骤是：

1. 遍历 A，填充 map。时间和空间复杂度都是 O(n)
2. 遍历 [0, inf]，找到第一个不在 map 中的元素

代码如下：

```{go}
func f(A []int) int {
  mapping := make(map[int]struct{}, len(A))
  for _, val := range A {
    mapping[val] = struct{}{}
  }
  for i := 0; ; i++ {
    if _, ok := mapping[i]; !ok {
      return i
    }
  }
}
```

该方法的时间和空间复杂度都达到了 O(n), 理论上达到了最优。但是从实践的角度考虑，hashmap 对空间的利用略超出了 O(n)，超出的范围取决于负载因子。严格来说，这个方法并不是最优 O(n) 解法。那么有没有更优的解法呢？答案是有，但是需要一些想象力，或者归纳能力。

# 方案四：拿出纸和笔，找规律

我们可以随意举出几个例子，对数组和结果进行分析。假设我们对数组已经排序，会有以下几种情况：

## 情况一：

1. A = [1,2,4,9], output = 0
2. A = [2,3,4,9], output = 0
3. A = 任意不包含 0 的数组，output = 0

## 情况二：

1. A = [], output = 0
2. A = [0], output = 1
3. A = [0,1], output = 2
4. A = [0,1,2], output = 3
5. A = [0,1,2,...,n-1], output = n

## 情况三：

1. A = [0,2,3,...,n-1], output = 1
2. A = [0,1,2,100,...,n-1], output = 3

对于所有情况，我们会发现一个共性：`output <= n`。为什么会出现这种情况？我们不妨逆向思考一下。假设有一个连续的自然数数组 `[0,1,2,3,...n-1]`：

1. 不改变这个数组，则返回值是 n；
2. 要改变一个元素，则必须从中取出一个自然数 i，然后替换成一个大于 n-1 的自然数。改变后的数组返回 i；
3. 改变多个元素，上一条仍然成立；
4. 交换任意两个元素，不影响返回值。

我们回到刚才提到的共性 `output <= n`，基于这个共性，我们可以认为数组 A 中所有大于 n 的数是无意义的。换句话说，我们只关心 A 中 `<= n` 的数字。

所以，解决方法来了：我们可以创建一个长度为 n+1 的 bool 数组 B，遍历 A 中所有元素 i 作为下标，设置 B[i]=true；然后找出 B 中第一个 false 的下标。代码如下：

```{go}
func f(A []int) int {
  n := len(A)
  B := make([]bool, n+1, n+1)
  for _, elem := range A {
    if elem <= n {
      B[elem] = true
    }
  }
  for idx, val := range B {
    if !val {
      return idx
    }
  }
}
```

理论上这个方法的时间和空间复杂度都是 O(n)，在实际运行时，都要略优于基于 map 的实现。

# 小结

有时候，一支笔，一张纸，可能比想破脑袋都好用（高智商人士当我没说）。

开头我们提到一篇博客对 Go 写算法的赞美之词。经过这波实践，依旧没有发现比 C++ 更简洁、更优美（手动狗头）。Go 简洁的语法让每一行都看起来那么短，但强制的大括号增加了代码行数，综合起来和C++写出来的效果旗鼓相当吧。

# 相关链接：

1. [GopherCon 2018 - Demystifying Binary Search Tree Algorithms](https://about.sourcegraph.com/go/gophercon-2018-binary-search-tree-algorithms/)

2. [Go maps in action](https://blog.golang.org/go-maps-in-action)

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")