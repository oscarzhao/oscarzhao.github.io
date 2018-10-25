---
layout: post
title:  "使用testify和mockery库简化单元测试"
date:   2018-10-25 11:59:00 +0800
categories: Golang go testing
---

# 前言

2016年我写过一篇关于[Go语言单元测试](2016-10-13-go-testing.markdown)的文章，简单介绍了 testing 库的使用方法。后来发现 [testify/require 和 testify/assert](https://github.com/stretchr/testify) 可以大大简化单元测试的写法，完全可以替代 `t.Fatalf` 和 `t.Errorf`，而且代码实现更为简短、优雅。

再后来，发现了 [mockery](https://github.com/vektra/mockery)，它可以为 Go interface 生成一个 mocks struct。通过 mocks struct，在单元测试中我们可以模拟所有 normal cases 和 corner cases，让 bug 无处藏身。在测试无状态函数 (对应 FP 中的 pure function) 时，mocks 的意义不大。mocks 的应用场景主要在于不可控的第三方服务、数据库、磁盘读写等，将细节封装到 interface 内部，只把方法暴露给调用方。这样一来，调用方的任何逻辑更改都可以被充分测试到。

关于 interface 的诸多用法，我会单独拎出来一篇文章来讲。本文中，我会通过两个例子展示 `testify/require` 和 `mockery` 的用法，分别是：

1. 使用 `testify/require` 简化 table driven test
2. 使用 `mockery` 和 `testify/mock` 为 lazy cache 写单元测试

# `testify/require`

@todo


# 小结

@todo

# 相关链接：

@todo

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")