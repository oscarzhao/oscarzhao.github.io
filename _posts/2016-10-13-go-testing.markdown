---
layout: post
title:  "Go语言单元测试"
date:   2016-10-13 11:59:00 +0800
categories: Golang testing go
---

## 简介

Go 语言在设计之初就考虑到了代码的可测试性。一方面 Go 本身提供了 [testing](https://golang.org/pkg/testing/ "testing") 库，使用方法很简单;
另一方面 go 的 package 提供了很多编译选项，代码和业务逻辑代码很容易解耦，可读性比较强（不妨对比一下C++测试框架）。 本文中，我们讨论的重点是 Go 语言中
的单元测试，而且只讨论一些基本的测试方法，包括下面几个方面：

1. 写一个简单的测试用例
2. Table driven test
3. 使用辅助测试函数（test helper）
4. 初始化（比如数据库连接）
5. 临时文件

这里我们只涉及到一些通用的测试方法。关于 HTTP server/client 测试，涉及到 Go 内置的 http server 和 web 框架，我会单独写一篇文章去讨论。

## 阅读建议

**Testing shows the presence, not the absence of bugs** -- [Edsger W. Dijkstra](https://en.wikiquote.org/wiki/Edsger_W._Dijkstra "Dijkstra")

为了保证业务逻辑代码的正确性，测试代码应当被赋予同等的重要性。在Github 上的开源代码中，我们经常可以看到一个指标“coverage”，即测试覆盖率。
最近几年新兴的大型项目，尤其是有多人参与的，大都有较高的代码覆盖率，如 [kubernetes](https://github.com/kubernetes/kubernetes "kubernetes") 84%, 
[react](https://github.com/facebook/react "react") 88%。很多项目没有给出 test coverage，但稍加注意就能发现不少测试代码，比如 [socket.io](https://github.com/socketio/socket.io "socket.io")、
[prometheus](https://github.com/prometheus/prometheus "prometheus")、[tensorflow](https://github.com/tensorflow/tensorflow "tensorflow")、[kafka](https://github.com/apache/kafka "kafka")。
由于 Go 语言工程化程度比较高，对测试支持比较完善，Github 上 Go 语言项目中测试更是随处可见。 

在阅读本文之前，建议您对 Go 语言的 package 有一定的了解，并在实际项目中使用过，下面是一些基本的要求：

1. 了解如何在项目中 import 一个外部的package
2. 了解如何将自己的项目按照功能模块划分 package
3. 了解 struct、struct字段、函数、变量名字首字母大小写的含义（非必需）
4. 了解一些 Go语言的编译选项，比如 +build !windows（非必需）

如果你对 1、2都不太了解，建议阅读一下这篇文章[How to Write Go Code](https://golang.org/doc/code.html, "go code")，动手实践一下。

## 写一个简单的测试用例

为了便于理解，我们首先给出一个代码片段：



## Table Driven Test 

## 辅助函数 (test helper)

## 初始化

## 临时文件
ioutil.TempDir，需要自己去清理文件

## 


### 相关链接：

1. [golang.org/pkg/testing](https://golang.org/pkg/testing/ "testing")
2. [Testing Techniques](https://talks.golang.org/2014/testing.slide “testing techniques")
3. [Table Driven Test](https://github.com/golang/go/wiki/TableDrivenTests "table driven test")
4. [Learn Testing](https://github.com/golang/go/wiki/LearnTesting "learn testing")

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")