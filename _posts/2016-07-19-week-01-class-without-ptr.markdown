---
layout: post
title:  "C++: 不含指针的class(geekband)"
date:   2016-07-19 23:28:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*class 是C++语言的核心概念之一，也是面向对象的基石。侯捷老师把class分为两类：*
*不含指针的class（如complex）和含指针的class（如string）。本文是听了侯捷老*
*师关于“不包含指针的class”的课程以后，总结而成的。*

*第一周：设计和实现一个不含指针的class。*

## 设计class的注意事项
1. 数据要封装（尽量private)
2. 参数和返回值尽量用reference（返回local variable除外）
3. 能用const一定要用const
4. 构造函数尽量用initializer_list进行参数初始化

## 临时对象
临时对象是指在`函数作用域`或`块作用域`内定义的变量。
临时对象在需要返回一个新对象时，可以有效减少代码的冗余，比如:
{% highlight c++ %}
// complex.h, derived from ppt
// ...
// 省略代码
// ...

inline complex
operator + (double x, const complex&y) 
{
  return complex(x + real(y), imag(y));
}

{% endhighlight %}

## friend（友元）
C++ class属性可以被public、protected、private修饰，从而对外表现出不同的形态，实现对对象的封装。
为了使用上的便捷，C++引入了friend修饰符，允许其他对象访问当前对象的私有成员。但同时破坏了封装性。什
么时候使用friend，是一个涉及到权衡和取舍的问题。

常用 friend的一个地方是 辅助函数，比如complex类的 `__doapl(complex*, const complex&)`
{% highlight c++ %}
// complex.h, derived from ppt

inline complex&
__doapl (complex* this, const complex& r)
{
  // ...
  return *this;
}
{% endhighlight %}
非常tricky的一点是：`同一个class的所有实例互为友元。`

## inline（内联函数）

### 什么是内联函数？ (参考[Microsoft MSDN](https://msdn.microsoft.com/en-us/library/bw1hbe6y.aspx))

inline函数执行起来更高效，关于其原因，MSDN上有一段描述：

~~~
The inline and __inline specifiers instruct the compiler to insert a copy 
of the function body into each place the function is called.

The insertion (called inline expansion or inlining) occurs only if the 
compiler's cost/benefit analysis show it to be profitable. Inline expansion 
alleviates the function-call overhead at the potential cost of larger code size.
~~~

大致意思是：编译时，调用inline函数的位置，都会被inline函数的函数体替换，而不是存放该函数的地址。

### 声明内联函数
inline关键字告诉编译器在编译时，将函数设置为inline。但是编译器不一定会这么做，有两类函数编译器肯定
不会编译成内联函数:

* 递归函数 (Recursive functions)
* 该函数在其它位置以指针的形式传递 (Functions that are referred to through a pointer elsewhere in the translation unit)
