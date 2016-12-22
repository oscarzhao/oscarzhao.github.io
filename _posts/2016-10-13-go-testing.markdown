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
最近几年新兴的大型项目，尤其是有多人参与的，大都有较高的代码覆盖率，如 [kubernetes](https://github.com/kubernetes/kubernetes "kubernetes"), [react](https://github.com/facebook/react "react")。
由于 Go 语言工程化程度比较高，对测试支持比较完善，Github 上 Go 语言项目中测试更是随处可见。 

在阅读本文之前，建议您对 Go 语言的 package 有一定的了解，并在实际项目中使用过，下面是一些基本的要求：

1. 了解如何在项目中 import 一个外部的package
2. 了解如何将自己的项目按照功能模块划分 package
3. 了解 struct、struct字段、函数、变量名字首字母大小写的含义（非必需）
4. 了解一些 Go语言的编译选项，比如 +build !windows（非必需）

如果你对 1、2都不太了解，建议阅读一下这篇文章[How to Write Go Code](https://golang.org/doc/code.html, "go code")，动手实践一下。

## 写一个简单的测试用例
为了便于理解，我们首先给出一个代码片段（如果你已经使用过go 的单元测试，可以跳过这个环节）：

```
// simple/equal.go
package simple

// a function to check if two numbers equals to each other.
func equal(a, b int) bool {
	return a == b
}

// simple/equal_test.go
package simple
import (
	"testing"
)

func TestEqual(t *testing.T) {
	a := 1
	b := 1
	shouldBe := true
	if real := equal(a, b); real == shouldBe {
		t.Errorf("equal(%d, %d) should be %v, but is:%v\n", a, b, shouldBe, real)
	}
}

```

上面这个例子中，如果你从来没有使用过单元测试，建议在本地开发环境中运行一次。这里有几点需要注意一下：

1. 这两个文件的父目录必须与包名一致（这里是 simple）
2. 测试用例的函数命名必须符合 TestXXX 格式，并且参数是 t *testing.T 
3. 了解一下 t.Errorf 与 t.Fatalf 的行为差异  


## Table Driven Test 

上面的测试用例中，我们一次只能测试一种情况，如果我们希望在一个 TestXXX 函数中进行很多项测试，Table Driven Test 就派上了用场。
举个例子，假设我们实现了自己的 [Sqrt](https://golang.org/pkg/math/#Sqrt "sqrt") 函数 mymath.Sqrt，我们需要对其进行测试：

首先，我们需要考虑一些特殊情况：

1. Sqrt(+Inf) = +Inf
2. Sqrt(±0) = ±0
3. Sqrt(x < 0) = NaN
4. Sqrt(NaN) = NaN

然后，我们需要考虑一般情况：

1. Sqrt(1.0) = 1.0
2. Sqrt(4.0) = 2.0
3. ...

注意：在一般情况中，我们对结果进行验证时，需要考虑小数点精确位数的问题。由于文章篇幅限制，这里不做额外的处理。

有了思路以后，我们可以基于 Table Driven Test 实现测试用例：

```
func TestSqrt(t *testing.T) {
	var shouldSuccess = []struct {
		input    float64 // input
		expected float64 // expected result
	}{
		{math.Inf(1), math.Inf(1)}, // positive infinity
		{math.Inf(-1), math.NaN()}, // negative infinity
		{-1.0, math.NaN()},
		{0.0, 0.0},
		{-0.0, -0.0},
		{1.0, 1.0},
		{4.0, 2.0},
	}
	for _, ts := range shouldSuccess {
		if actual := Sqrt(t.input); actual != ts.expected {
			t.Fatalf("Sqrt(%f) should be %v, but is:%v\n", ts.input, ts.expected, actual)
		}
	}
}
``` 

## 辅助函数 (test helper)

## 初始化

## 临时文件
如果待测试的功能模块涉及到文件操作，临时文件是一个不错的解决方案。go语言的 ioutil 包提供了 TempDir 和 
TempFile 方法，供我们使用。

我们以 etcd 创建 wal 文件为例，来看一下 TempDir 的用法：

```
// github.com/coreos/etcd/wal/wal_test.go

func TestNew(t *testing.T) {
	p, err := ioutil.TempDir(os.TempDir(), "waltest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(p)  // 千万不要忘记删除目录

	w, err := Create(p, []byte("somedata"))
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if g := path.Base(w.tail().Name()); g != walName(0, 0) {
		t.Errorf("name = %+v, want %+v", g, walName(0, 0))
	}
	defer w.Close()

	// 将文件 waltest 中的数据读取到变量 gb []byte 中 
	// ...

	// 根据 "somedata" 生成数据，存储在变量 wb byte.Buffer 中
	// ...

	// 临时文件中的数据（gb）与 生成的数据（wb）进行对比
	if !bytes.Equal(gd, wb.Bytes()) {
		t.Errorf("data = %v, want %v", gd, wb.Bytes())
	}
}
```

上面这段代码是从 etcd 中摘取出来的，源码查看 [coreos/etcd - Github](https://github.com/coreos/etcd/blob/2353cbca719f6661c8642d666dd8e16591f5ebb5/wal/wal_test.go "coreos/etcd")。
需要注意的是，使用 [TempDir](https://golang.org/pkg/io/ioutil/#TempDir "TempDir") 和 [TempFile](https://golang.org/pkg/io/ioutil/#TempFile "TempFile") 创建文件以后，需要自己去删除。

### 相关链接：

1. [golang.org/pkg/testing](https://golang.org/pkg/testing/ "testing")
2. [Testing Techniques](https://talks.golang.org/2014/testing.slide “testing techniques")
3. [Table Driven Test](https://github.com/golang/go/wiki/TableDrivenTests "table driven test")
4. [Learn Testing](https://github.com/golang/go/wiki/LearnTesting "learn testing")

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")