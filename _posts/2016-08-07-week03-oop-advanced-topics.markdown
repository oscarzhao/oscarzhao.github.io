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

#### 1.3.1 shared_ptr 智能指针
实现一个智能指针类，则必须重载 (override) 两个成员函数： operator*() 和 operator->()。 

shared_ptr 代码的抽离出来一部分如下：

``` c++
// shared_ptr 模板类的定义
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
};

// Foo 类
struct Foo {
  // 省略其它部分，关注 method 
  void method() { //... }
};

// 使用 shared_ptr
int main() {
  shared_ptr<Foo> sp(new Foo);

  Foo f(*sp);  // 使用 解引用 操作符 *
  sp->method();  // 使用 操作符 ->  ( 等价于 px->method() )
}
```

*关于 operator -> ，从语法角度上来说，-> 是可重生的，所以下面 main函数中才可以这样使用。*

很容易发现，除了构造，shared_ptr 在使用上与裸指针几乎没有差别，但是不需要手动释放内存。
当然，仿指针类的能力远不止于自动释放内存，还有很多，这里我们看看标准库中 std::shared_ptr的附加功能。

std::shared_ptr 不仅提供了有限的垃圾回收特性，还提供了内存拥有权的管理 (ownership)，点击[这里查看详情](http://www.cplusplus.com/reference/memory/shared_ptr/?kw=shared_ptr "shared_ptr") 

#### 1.3.2 iterator 迭代器

pointer-like classes 在迭代器中也有广泛的应用。
标准库中所有的容器(std::vector等) 都有迭代器。换句话说，标准库的迭代器也实现了 operator* 和 operator-> 方法。

每个迭代器对象 指向 一个容器变量，但同时实现了下面几个方法：

1. operator==
2. operator!=
3. operator++
4. operator--

关于 迭代器中 operator* 和 operator-> 的实现，也相当值得考究：

``` c++
// 忽略上下文

reference operator*() const {
  return (*node).data;
}

pointer operator->() const { // 借助于 operator* 实现
  return &(operator*());
}
```

你可以像下面这样使用这两个方法：

``` c++
list<Foo>::iterator ite;

//... 省略一部分代码...

*ite;   // 获取 Foo 对象的引用

ite->method();  
// 意思是 调用 Foo::method()
// 相当于 (*ite).method();
// 相当于 (&(*ite))->method();

```

### 1.4 function-like classes (仿函数)
#### 1.4.1 什么是仿函数？
仿函数其实不是函数，是一个类，但是它的行为和函数类似。在实现的层面上，一个类一旦定义了 operator() 方法，就可以称之为仿函数。

C++标准库内置了很多[仿函数模板](http://www.cplusplus.com/reference/functional/, "function")。
我们先用 std::less 和 std::less_equal 为例，对仿函数的用法有一个直观的认识：

``` c++
// less example (http://www.cplusplus.com/reference/functional/less/)
// compile: g++ -o main main.cpp -lm
#include <iostream>     // std::cout
#include <functional>   // std::less
#include <algorithm>    // std::sort, std::includes

int main () {
  // 自己写的简单例子, 表达式 "std::less<int>()" 创建了一个临时对象  
  int a = 5, b = 4;
  std::cout << "std::less<int>()(" << a << ", " << b << "): " << std::less<int>()(a, b) << std::endl;
  std::cout << "std::less<int>()(" << b << ", " << a << "): " << std::less<int>()(b, a) << std::endl;

  std::cout << "std::less_equal<int>()(" << a << ", " << b << "): " << std::less_equal<int>()(a, b) << std::endl;
  std::cout << "std::less_equal<int>()(" << b << ", " << a << "): " << std::less_equal<int>()(b, a) << std::endl;
  std::cout << "std::less_equal<int>()(" << a << ", " << a << "): " << std::less_equal<int>()(a, a) << std::endl;

  // 网站上带的高级例子
  int foo[]={10,20,5,15,25};
  int bar[]={15,10,20};
  std::sort (foo, foo+5, std::less<int>());  // 5 10 15 20 25
  std::sort (bar, bar+3, std::less<int>());  //   10 15 20
  if (std::includes (foo, foo+5, bar, bar+3, std::less<int>()))
    std::cout << "foo includes bar.\n";
  return 0;
}
```
#### 1.4.2 仿函数的实现 

仿函数实际上是一个 类 (class)，这个类实现了 operator() 方法。
下面这个是 C++98 中 std::less 的实现：

``` c++
// C++11 中的实现（侯捷老师讲的是 C++98中的实现）
template <class T> struct less {
  bool operator() (const T& x, const T& y) const {return x<y;}
  typedef T first_argument_type;
  typedef T second_argument_type;
  typedef bool result_type;
};
```

注意：std::less 是类模板。在课程中，侯捷老师提到了 unary_function 和 binary_function，这两个类定义了参数类型，C++11中已经不再
使用，而是内置到 std::less 中，具体参考[这里](http://www.cplusplus.com/reference/functional/less/ "less")

### 1.5 未完待续
