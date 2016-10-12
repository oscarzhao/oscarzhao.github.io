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

在我的第一份工作中，从来没有见过、也没有听说过开发人员自己去写测试代码，测试也仅仅局限于测试人员的白盒测试。
开发人员写完代码后，只进行简单的功能测试，保证最经常走的逻辑能走通，然后就移交给测试人员。结果是 bug 表现在
软件的顶层（成型的产品），一方面难以定位，而且随着系统复杂性的增高越来越难；另一方面 跨部门的沟通导致 bug 的反馈链路被拉长。

来到云计算行业以后，开始用 kubernetes，努力去理解它的设计理念，逐步了解了 DevOps、Microservice等概念。
DevOps 理念中，测试是非常重要的一部分；在Kubernetes 项目中，单元测试、集成测试都比比皆是。比葫芦画瓢，
在公司的项目中也开始尝试自己写测试用例。后来的事实证明，绝大多数 bug 都可以通过单元测试检查出来并修复，代码行为的可控性提高不少。

总之，测试非常重要。

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

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")