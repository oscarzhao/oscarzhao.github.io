---
layout: post
title:  "C++ OOP: 不含指针的class (Geekband)"
date:   2016-07-19 23:28:00 +0800
categories: C++
---
# 第一周：设计和实现一个不含指针的class

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
