---
layout: post
title:  "使用testify和mockery库简化单元测试"
date:   2018-10-25 11:59:00 +0800
categories: Golang go testing
---

# 前言

2016年我写过一篇关于[Go语言单元测试](2016-10-13-go-testing.markdown)的文章，简单介绍了 testing 库的使用方法。后来发现 [testify/require 和 testify/assert](https://github.com/stretchr/testify) 可以大大简化单元测试的写法，完全可以替代 `t.Fatalf` 和 `t.Errorf`，而且代码实现更为简短、优雅。

再后来，发现了 [mockery](https://github.com/vektra/mockery) 库，它可以为 Go interface 生成一个 mocks struct。通过 mocks struct，在单元测试中我们可以模拟所有 normal cases 和 corner cases，让 bug 无处藏身。在测试无状态函数 (对应 FP 中的 pure function) 时，mocks 的意义不大。mocks 的应用场景主要在于不可控的第三方服务、数据库、磁盘读写等，将细节封装到 interface 内部，只把方法暴露给调用方。这样一来，调用方的任何逻辑更改都可以被充分测试到。

关于 interface 的诸多用法，我会单独拎出来一篇文章来讲。本文中，我会通过两个例子展示 `testify/require` 和 `mockery` 的用法，分别是：

1. 使用 `testify/require` 简化 table driven test
2. 使用 `mockery` 和 `testify/mock` 为 lazy cache 写单元测试

# testify/require

首先，我们通过一个简单的例子看下 require 的用法。我们针对函数 `Sqrt` 进行测试，其实现为：

```{go}
// Sqrt calculate the square root of a non-negative float64 
// number with max error of 10^-9. For simplicity, we don't 
// discard the part with is smaller than 10^-9.
func Sqrt(x float64) float64 {
  if x < 0 {
    panic("cannot be negative")
  }

  if x == 0 {
    return 0
  }

  a := x / 2
  b := (a + 2) / 2
  erro := a - b
  for erro >= 0.000000001 || erro <= -0.000000001 {
    a = b
    b = (b + x/b) / 2
    erro = a - b
  }

  return b
}
```

这里我们使用了一个常规的方法实现 `Sqrt`，该实现的最大精确度是到小数点后9位（为了方便演示，这里没有对超出9位的部分进行删除）。我们首先测试 `x < 0` 导致 panic 的情况，看 `require` 如何使用，下面是测试代码：

```{go}
func TestSqrt_Panic(t *testing.T) {
  defer func() {
    r := recover()
    require.Equal(t, "cannot be negative", r)
  }()
  _ = Sqrt(-1)
}
```

在上面的函数中，我们只使用 `require.Equal` 一行代码就实现了运行结果校验。如果使用 `testing` 来实现的话，变成了三行，而且需要手写一串描述：

```{go}
func TestSqrt_Panic(t *testing.T) {
  defer func() {
    r := recover()
        if r.(string) != "cannot be negative" {
            t.Fatalf("expect to panic with message \"cannot be negative\", but got \"%s\"\n", r)
        }
  }()
  _ = Sqrt(-1)
}
```

使用 `require` 之后，不仅使测试代码更易于编写，而且能够在测试运行失败时，格式化运行结果，方便定位和修改bug。这里你不妨把 `-1` 改成一个正数，运行 `go test`，查看运行结果。

上面我们能够看到 `require` 库带来的编码和调试效率的上升。在 table driven test 中，我们会有更深刻的体会。

## Table Driven Test

我们仍然以 `Sqrt` 为例，来看下如何在 table driven test 中使用 `require`。这里我们测试的传入常规参数的情况，代码实现如下：

```{go}
func TestSqrt(t *testing.T) {
  testcases := []struct {
    desc   string
    input  float64
    expect float64
  }{
    {
      desc:   "zero",
      input:  0,
      expect: 0,
    },
    {
      desc:   "one",
      input:  1,
      expect: 1,
    },
    {
      desc: "a very small rational number",
      input: 0.00000000000000000000000001,
      expect: 0.0,
    },
    {
      desc:   "rational number result: 2.56",
      input:  2.56,
      expect: 1.6,
    },
    {
      desc:   "irrational number result: 2",
      input:  2,
      expect: 1.414213562,
    },
  }

  for _, ts := range testcases {
    got := Sqrt(ts.input)
    erro := got - ts.expect
    require.True(t, erro < 0.000000001 && erro > -0.000000001, ts.desc)
  }
}
```

在上面这个例子，有三点值得注意：

1. `匿名struct` 允许我们填充任意类型的字段，非常方便于构建测试数据集；
2. 每个`匿名struct`都包含一个 `desc string` 字段，用于描述该测试要处理的状况。在测试运行失败时，非常有助于定位失败位置；
3. 使用 `require` 而不是 `assert`，因为使用 `require` 时，测试失败以后，所有测试都会停止执行。

关于 `require`，除了本文中提到的 `require.True`, `require.Equal`，还有一个比较实用的方法是 `require.EqualValues`，它的应用场景在于处理 Go 的强类型问题，我们不妨看一段代码：

```{go}
func Test_Require_EqualValues(t *testing.T) {
	// tests will pass
	require.EqualValues(t, 12, 12.0, "compare int32 and float64")
	require.EqualValues(t, 12, int64(12), "compare int32 and int64")

	// tests will fail
	require.Equal(t, 12, 12.0, "compare int32 and float64")
	require.Equal(t, 12, int64(12), "compare int32 and int64")
}
```

更多 `require` 的方法参考 [require's godoc](https://godoc.org/github.com/stretchr/testify/require)。

# mockery



# 小结

@todo

# 相关链接：

@todo

扫码关注微信公众号“深入Go语言”

![在这里]( http://oat5ddzns.bkt.clouddn.com/qrcode_for_gh_9280bd217b46_430.jpg "qrcode")