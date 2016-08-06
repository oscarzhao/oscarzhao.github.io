---
layout: post
title:  "C++: 面向对象简介(geekband)"
date:   2016-07-31 11:59:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*class 是C++语言的核心概念之一，也是面向对象的基石。侯捷老师把class分为两类：*
*不含指针的class（如complex）和含指针的class（如string）。本文是听了侯捷老*
*师关于“包含指针的class”的课程以后，总结而成的。由于上个博客对设计class进行了*
*很多讲解，本文重点讲解class面向对象的特征。*

*第二周：设计和实现一个包含指针的class。*

## 构造函数和析构函数
对于包含指针作为成员变量的class，需要设计构造函数和析构函数，以便内存管理。这三个函数的声明方式如下：

``` c++
class String
{
public:                                 
   String(const char* cstr=0);                     
   String(const String& str);                    
   String& operator=(const String& str);         
   ~String();                                    
   char* get_c_str() const { return m_data; }
private:
   char* m_data;
};
```

注意： 

1. 拷贝构造函数和赋值构造函数的传参为reference。
2. 析构函数要考虑释放 new 分配的内存
3. 构造函数、拷贝构造函数和赋值构造函数 需要注意 分配内存 （必要时）

## static 关键字

### static 修饰成员变量
特点： 该成员变量属于 class，不属于特定的class 实例。

应用场景： 银行利率。

初始化： 在类的外部进行初始化。

### static 修饰成员函数
特点： 该方法属于 class，不属于特定的 class 实例。

应用场景： 处理static成员变量。

调用： 可以使用 class 名称进行调用；或者通过 class 的任何实例调用。

看下面这个例子：

``` c++
class Account
{
public:
    static double m_rate;
    static void set_rate(const double& x) {m_rate = x;}
};

double Account::m_rate = 8.0;

int main() {
  Account::set_rate(5.0);  // 通过 class 名调用

  Account a;
  a.set_rate(7.0);  // 通过 class 实例调用
}
```

### 一个用途： Singleton （单例模式）
看代码：

``` c++
class A
{
public:
    static A& getInstance();
    setup(){/*...*/}
private:
    A();
    A(const A& rhs);
    // ...
};

A& A::getInstance() {
  static A a;  // 第一次调用时，才会创建 a
  return a;
}
```

## class template 和 function template
区别： class template在使用时，需要声明类型；function template 使用时，模板会自动推导类型，使用上和正常的函数一样。

## 未完待续 ...