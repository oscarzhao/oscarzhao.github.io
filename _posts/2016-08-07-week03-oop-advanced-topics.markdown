---
layout: post
title:  "C++ 对象模型（geekband)"
date:   2016-08-07 11:59:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*本文是听了侯捷老师关于“C++对象模型”的课程以后，总结而成的。*
*课程中讲到了范型编程和面向对象编程两种模式，因此本文的主题是 template。*
*包括类模板、函数模板、成员模板、模板特化和偏特化、可变模板参数、模板嵌套等。*

*说到面向对象，本文的另一个主题是C++底层的对象模型，包括this指针、虚指针vptr、虚表vtable、虚函数、多态等。*

## 第三周：C++对象模型

## Part 1：class 的高级用法

这一部分，我们介绍一些C++ 类的高级用法：conversion function、non-explicit-one-argument constructor、pointer-like
classes、function-like classes。

### 1.1 conversion function (转换函数)

转换函数 可以将对象 默认转换为另一种类型，方便程序的调用。

其作用与只接收一个参数的构造函数（不使用 explicit修饰符）作用正好相反。关于构造函数的坑参考[这篇文章](http://tipsandtricks.runicsoft.com/Cpp/Explicit.html "explicit")。
下一个小结讲解的就是 non-explicit-one-argument constructor。

废话少说，先上个例子（注意观察它的使用方法）：

```c++
// Fraction 类
class Fraction {
public:
  Fraction(int num, int den=1): m_numerator(num), m_denominator(den){};
  operator double() const {return (double)m_numerator/m_denominator; };
  private:
    int m_numerator;
    int m_denominator;
};

// main 函数（使用 Fraction 的转换函数 operator double() const)
int main(int argc, char * argv[]) {
  Fraction f(3, 5);
  double d = 4.0 + f;      // f 调用 operator double() 转换为一个 double 值 
  std::cout << "d = " << d << std::endl;
}
```

编译器认为必要时，将对象转换为某种特定的类型。转换函数的三个特征：

1. 转换函数 的声明中不需要写返回类型 (因为返回类型和函数名称相同)
2. 函数名必须是 operator xxx
3. 必须使用 const 修饰符，因为转换函数不会修改 this 指针

### 1.2 non-explicit-one-argument constructor 

non-explicit-one-argument constructor 有下面两个语法特征：

1. 这类构造函数只接收一个参数
2. 没有使用 explicit 修饰符 （默认 implicit）

该函数在编译器认为需要的时候，将参数类型的变量转换成该类型的对象。具体看下面代码：

``` c++
// Fraction 类
class Fraction {
  public:
    Fraction(int num, int den=1): m_numerator(num), m_denominator(den){};
    Fraction operator+(const Fraction& f) {
      return Fraction(m_denominator * f.m_numerator + m_numerator * f.m_denominator, m_denominator * f.m_denominator); 
    };

    friend ostream& operator << (ostream& os, const Fraction& f);
  private:
    int m_numerator;
    int m_denominator;
};

ostream& operator<<(ostream& os, const Fraction& f) {
  return os << f.m_numerator << "/" << f.m_denominator;
}

// main 函数（注意变量 d 的类型）
int main(int argc, char * argv[]) {
  Fraction f(3, 5);
  Fraction d =  f + 4;  // here 4 is converted to a Fraction(4, 1), 4+f is not allowed
  std::cout << "d = " << d << std::endl;
}
```

使用 explicit 修饰构造函数以后，上面的默认转换就要失败。

### 1.3  pointer-like class (智能指针)
像指针的类，是指该 class 的对象表现出来像一个指针。这么做的目的是实现一个比指针更强大的结构。
标准库内值了一些指针类，如 std::shared_ptr, std::weak_ptr, std::unique_ptr，具体参考[cplusplus](http://www.cplusplus.com/reference/memory/ "memory")

实现一个智能指针类，则必须重载 (override) 两个成员函数： operator*() 和 operator->()。 

shared_ptr 代码的抽离出来一部分如下：

``` c++
template <typename T>
class shared_ptr{
public:
  T& operator*() const { // 解引用
    return *px;
  }
  T* operator->() const { // 取指针，这个方法有点诡异
    return px;
  }
  
  shared_ptr(T *p): px(p) {}
private:
  T* px;
  long* pn;
}
```

### 1.4 未完待续～～