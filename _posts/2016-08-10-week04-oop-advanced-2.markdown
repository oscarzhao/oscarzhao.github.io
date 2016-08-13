---
layout: post
title:  "C++ 对象模型(二)（geekband)"
date:   2016-08-10 11:59:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*本文是听了侯捷老师关于“C++对象模型”的课程以后，总结而成的。*
*课程中讲到了范型编程和面向对象编程两种模式，上篇博客中，我们讨论了类和模板的一些高级话题。*
*本文中，我们对继承体系下的类的构成进行深入讨论，相关话题包括this指针、虚指针vptr、虚表vtable、虚函数、多态等。*
*另外，new 和 delete 的重载是本文的另一个话题。*

## Part 1：虚表、虚指针、动态绑定

这一部分，我们介绍下 继承体系下，类和对象的存储形式。

### 1.1 vptr 虚指针 和 vtable 虚表

对于虚指针和虚表的定义，这里引用一段 [quora](https://www.quora.com/What-are-vTable-and-VPTR-in-C++ "c++") 上的一个回复(这里我已经翻译成中文)：
如果一个类存在一个或多个虚函数，编译器会为这个类的实例 (对象) 创建一个隐藏的成员变量，即虚指针(virtual-pointer)，简称 vptr。
vptr 指向一个包含一组函数指针的表，我们称之为 虚表 (virtual table)，简称 vtable。虚表由编译器创建，虚表中的每一项均是
一个指向对应虚函数的指针。

为了实现动态绑定 (dynamic binding)，编译器为每一个拥有虚函数的类 (和它的子类) 创建一张虚表。编译器将虚函数的地址存放到对应
类的虚表中。 当通过基类指针 (或父类指针，Base * pb) 调用虚函数时，编译器插入一段在虚表中查找虚函数地址和获取 vptr 的代码。
所以才能够调用到"正确"的函数，实现动态绑定。

关于 vptr 和 vtable 的调用，这里用侯捷老师 PPT 上的一张图表示：

![vptr-vtable](http://obi1zst3q.bkt.clouddn.com/blog/cpp/vptr-vtable.jpg "vtable-vptr")

关于 类 A、B、C 的结构声明参考下面的代码 (注意这里不包含具体实现)：

``` c++
// 上图中 类 A、B、C 的声明
class A {
public:
  virtual void vfunc1();
  virtual void vfunc2();
  void func1();
  void func2();
private:
  int m_data1, m_data2;
}

class B: public A {
public:
  virtual void vfunc1();
  void func2();
private:
  int m_data3;
}

class C: public B {
public:
  virtual void vfunc1();
  void func2();
private:
  int m_data1, m_data4;
}
```


### 1.2 this pointer (template method)

在继承体系中，子类对象调用一个方法时，如果该类本身这个函数，那么会调用这个函数；如果本身没有，那么编译器会沿着继承树向上查找父类中是否有该方法。

侯捷老师PPT中的一张图很好地体现了这种调用机制：

![虚函数-子类调用](http://obi1zst3q.bkt.clouddn.com/blog/cpp/20160811-this-pointer-2.jpg "virtual-inheritance")

### 1.3 dynamic binding 动态绑定

#### 1.3.1 什么是动态绑定？

动态绑定是编程语言的一种特性（或机制），它允许程序在运行时决定执行操作的细节，而不是在编译时就确定。在设计一个软件时，通常会出现下面两类情况：

1. 类的接口已经确定，但是还不知道具体怎么实现
2. 开发者知道需要什么算法，但是不知道具体的操作

这两种情况下，开发者都需要延迟决定，延迟到什么时候呢？延迟到已经有足够的信息去做一个正确的决策。此时如果能不修改原先的实现，我们的目标就达到了。

动态绑定正是为了满足这些需求而存在，结果就是更灵活和可伸缩的软件架构。比如在软件开发初期，不需要做出所有设计决策。这里我们讨论下灵活性和可伸缩性：

1. flexibility (灵活性): 很容易将现存组件和新的配置合并到一起
2. extensibility (扩展性)： 很容易添加新组件

C++ 通过 虚表和虚指针机制 实现对动态绑定的支持，具体的机制我们在上面已经谈到，这里不再赘述。

#### 1.3.2 动态绑定在 C++ 中的体现

在 C++ 中，动态绑定的标志是在声明类方法时，在方法名前面显式地添加 virtual 关键字。比如下面这样：

``` c++
class Base {
public:
  virtual void vfunc1() { std::cout << "Base::vfunc1()" << std::endl; }
  void func1() { std::cout << "Base::func1()" << std::endl; }
}
```

只有类的成员函数才能被声明为虚函数，下面三种是不可以的：

1. 普通的函数 (不属于任何一个类)
2. 类的成员变量
3. 静态方法 (static 修饰的成员函数)

virtual 修饰的成员函数的接口是固定的，但是子类中的同名成员函数可以修改默认实现，比如像下面这样：

``` c++
class Derived_1 {
public:
  virtual void vfunc1() { std::cout << "Derived_1::vfunc1() " << std::endl; }
}
```

注意：上面的代码中， virtual 是可选的，即便不写，它仍然是虚函数！

在程序运行时，虚函数调度机制会根据对象的"动态类型"选择对应的成员函数。
被选择的成员函数依赖于被指针指向的对象，而不是指针的类型。看下面代码：

``` c++
void foo (Base *bp) { bp->vf1 (); /* virtual */ }
Base b;
Base *bp = &b;
bp->vf1 (); // 打印 "Base::vfunc1()"
Derived_1 d;
bp = &d;
bp->vf1 (); // 打印 "Derived_1::vfunc1()"
foo (&b); // 打印 "Base::vfunc1()"
foo (&d); // 打印 "Derived_1::vfunc1()"，这里存在一个隐式的向上转型
```

关于动态绑定，更多细节参考 [C++ dynamic binding](http://www.cs.wustl.edu/~schmidt/PDF/C++-dynamic-binding4.pdf "dynamic-binding")。

## Part 2: const 补充

这个小结中，关于 const 的所有例子均来自于 [msdn](https://msdn.microsoft.com/en-us/library/07x6b05d.aspx "const")。为了便于理解，
对代码进行了稍微的调整。

### 2.1 const 修饰指针
下面这个例子中， const 修饰的是指针，因此不能修改指针 aptr 的值，即 aptr 不能指向另一个位置。

``` c++
// constant_values3.cpp
int main() {
   char *mybuf = 0, *yourbuf;
   char* const aptr = mybuf;
   *aptr = 'a';   // OK
   aptr = yourbuf;   // C3892
} 
```

### 2.2 const 修饰指针指向的数据
下面这个例子中， const 修饰的是指针指向的数据，因此可以修改指针的值，但是不能修改指针指向的数据。

``` c++
// constant_values4.cpp
#include <stdio.h>
int main() {
   const char *mybuf = "test";
   char* yourbuf = "test2";
   printf_s("%s\n", mybuf);

   const char* bptr = mybuf;   // Pointer to constant data
   printf_s("%s\n", bptr);

   // *bptr = 'a';   // Error
}
```

### 2.3 const 修饰成员函数
在声明成员函数时，如果在函数末尾使用 const 关键字，那么可以称这个函数是"只读"函数。 const成员函数不能修改任何 非static的成员变量，
也不能调用任何 非const 成员函数。

const成员函数在`声明`和`定义`时，都必须带有 const 关键字。看下面这个例子：

``` c++
// constant_member_function.cpp
class Date
{
public:
   Date( int mn, int dy, int yr );
   int getMonth() const;     // A read-only function
   void setMonth( int mn );   // A write function; can't be const
private:
   int month;
};

int Date::getMonth() const
{
   return month;        // Doesn't modify anything
}
void Date::setMonth( int mn )
{
   month = mn;          // Modifies data member
}
int main()
{
   Date MyDate( 7, 4, 1998 );
   const Date BirthDate( 1, 18, 1953 );
   MyDate.setMonth( 4 );    // Okay
   BirthDate.getMonth();    // Okay
   BirthDate.setMonth( 4 ); // C2662 Error
}
```

## Part 3：new 和 delete

### 3.1 分解 new 和 delete

new 和 delete 都是表达式，因此不能被重载。它们均有不同步骤组成：

new 的执行步骤：

1. 调用operator new 分配内存 (malloc)
2. 对指针进行类型转换
3. 调用构造函数

delete 的执行步骤：

1. 调用析构函数
2. 调用operator delete释放内存 (free)

虽然，new 和 delete 不能被重载，但是 operator new 和 operator delete 可以被重载。