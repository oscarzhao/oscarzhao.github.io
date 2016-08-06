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

## Part 1：面向对象的语法基础

### 1.1 构造函数和析构函数
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

### 1.2 static 关键字

#### 1.2.1 static 修饰成员变量
特点： 该成员变量属于 class，不属于特定的class 实例。

应用场景： 银行利率。

初始化： 在类的外部进行初始化。

#### 1.2.2 static 修饰成员函数
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

#### 1.2.3 一个用途： Singleton （单例模式）
单例模式下，构造函数被设置为 private， 通过 <class名>::getInstance() 函数获取 全局唯一的实例。

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

### 1.3 class template 和 function template
模板类(函数) 允许定义一类功能类似的 类(函数)，在`编译期展开`为具体的类（函数）。
两者的区别： class template在使用时，需要声明类型；function template 使用时，模板会自动推导类型，使用上和正常的函数一样。

标准库内置了很多模板类和模板函数。

常用的模板类有：

1. std::vector
2. std::deque
3. std::map
4. std::array
5. std::stack 
6. 更多参考[C++中的容器 (Containers)](http://www.cplusplus.com/reference/stl/ "containers")

常用的模板函数有：

1. std::less
2. std::less_equal
3. std::greater
4. std::greater_equal
5. 更多参考 [C++中的函数模板](http://www.cplusplus.com/reference/functional/ "functional")

### 1.4 namespace (命名空间)
在大学学习C++ 时，写的代码都比较简单，默认使用 std 作为 namespace，对自定义命名空间没有太大的需求。
但是在公司的大型项目中，多个小组实现不同的功能时，可能会存在函数或类名冲突的问题，此时可以引入命名空间。

题外话：命名空间的概念也存在于 linux 内核中，用于资源隔离，docker的诞生严重依赖于这种隔离机制。

#### 1.4.1 声明一个命名空间
可以使用下面的语法声明一个命名空间，在多个文件中重复声明是有效的：

``` c++
namespace std {
    // ...
}
```

#### 1.4.2 使用命名空间
从使用的角度，C++ 提供了三种方式：

1. using directive
该方法会将 整个命名空间下的所有函数引入到使用者的命名空间，懒人专用。

``` c++
#include<iostream>
using namespace std; 

int main() {
    cin >>  ...;
    cout << ...;
}
```

2. using declaration
该方法可以单独引入某个命名空间下的特定函数，便于定制化。

``` c++
#include<iostream>
using std::cout;

int main() {
    std::cin >> ...;
    cout << ...;
}
```

3. 使用函数（类）的全名
该方法可以发挥命名空间的优势，但是每次调用都要带上 命名空间名，想对麻烦。

``` c++
#include<iostream>

int main() {
    std::cin >> ...;
    std::cout << ...;
}
```

## Part 2：面向对象中的三种关系

面向对象的三大关系：继承(inheritance)、组合(composition)、委托(delegation)。

### 2.1 组合 (composition)

组合就是我们通常所说的 has-a 关系。举个栗子：

``` c++
template <typename T, class Sequence = deque<T> >
class queue {
//...
protected:
    Sequence c; // 底层容器
public:
    // 下面的函数均通过 c 的成员函数完成
    bool empty() const {return c.empty();}
    size_type size() const { return c.size();}
    reference front() { return c.front(); }
    reference back() { return c.back(); }
    // deque 的两端均可进出，queue 是 first in first out
    void push(const value_type& x) { c.push_back(x); }
    void pop() { c.pop_front(); }
};
```

这个例子是 标准库中 queue (队列) 的实现。队列的特征是 first in first out，
deque 是双端队列，一般由一个双向链表实现，头尾都可以进出。很容易发现，deque 的操作 是 queue 的超集，
因此 queue 可以借助于 deque 实现。

has-a 关系 用实心菱形的箭头表示：

![has-a](http://obi1zst3q.bkt.clouddn.com/has_a "has-a")

组合关系下，两个类的内存结构如下：

![has-a](http://obi1zst3q.bkt.clouddn.com/20160807_has_a "has-a")

以 queue 和 deque 为例。创建时，先调用 queue 的构造函数；销毁时，先调用 deque 的析构函数。两者的生命周期是一致的。

### 2.2 委托 (delegation)

委托是指 类本身只提供接口，具体的实现“委托”给另一个类。为了便于理解，我们参考 标准库 string 的实现：

``` c++
// file string.hpp
class StringRep;
class String {
public:
    String();
    String(const char* s);
    String(const String& s);
    String& operator=(const String& s);
    ~String();
// 省略一些函数
private:
    StringRep* rep; // 负责具体的实现
};

// file string.cpp
#include "string.hpp"
namespace {
class StringRep{
    friend class String;
    StringRep(const char* s);
    ~StringRep();
    int count;
    char* rep;
};

String::String() { ... }
}
```

这个例子中，类 String 定义了一些接口，但具体的实现全部交给了 StringRep 处理。StringRep 实现了引用计数 (reference counting)机制。
需要注意的是：这里 String  和 StringRep 的生命周期并不一致，这里也正是 引用计数 的巧妙所在。

很多地方谈到引用计数时，只有一句话：引用计数为0时，释放内存。这里我总结了引用计数生效的几个关键时间点：

1. 使用 拷贝构造函数、赋值构造函数 初始化一个新的 String时，底层 StringRep 对象的 count  +1。
2. 销毁 一个 String 时，其背后的 StringRep 对象的 count -1。
3. 当 StringRep 的 count 为 0 时， String 的析构函数 调用  delete rep，释放底层的 StringRep 对象。
4. 修改 一个 string 对象时，不会影响 和它具有同一个底层 StringRep 的 string 对象。如果 底层 StringRep 对象的 count == 1，则修改 StringRep 对象；否则为其创建一个新的StringRep对象，并将原 rep 的count -1。

C++ 没有垃圾回收机制，通过推理只能得到这样的结论。这里给出一段代码示例：

``` c++
#include <iostream>

int main() {
    string s("Hello");  // 创建 string s，假设背后的 StringRep 为 sr; sr.count == 1
    string s1(s);       // 创建 string s1, 背后仍然是 sr; sr.count == 2
    s1[0] = 'h';        // 修改 s1 不影响 s，为 s1 创建一个新的 StringRep 对象; src.count == 1
    string s2 = s;         // 创建 string s2, 背后仍然是 sr; sr.count == 2
    string* ps1 = new string(s); // 创建 string 指针 ps1, 背后对象仍然是 sr; sr.count == 4
    delete ps1;            // 删除 ps1，sr.count == 2
    
    string t("Hello");   // 创建 string 对象 t，t 背后的 StringRep 不是 sr，而是一个全新的 StringRep对象
    
    // main 函数退出时，这些变量都会被析构掉。
}
```

String 和 StringRep 对象的关系如下图所示：

![reference counting](http://obi1zst3q.bkt.clouddn.com/20160807_reference_counting "reference")

在 UML 图中，使用空心菱形的箭头表示委托关系：

![delegation](http://obi1zst3q.bkt.clouddn.com/20160807_delegation "delegation")

### 2.3 继承 (inheritance)
继承的语法是 定义一个类时，使用冒号。在UML图中，使用空心三角形的箭头表示：

![继承](http://obi1zst3q.bkt.clouddn.com/20160807_inheritance "inheritance")

继承关系中，我们称被继承的类为 父类 (base class)， 继承者类 称为 子类 (derived class)。

在讨论继承时，我们所说的是 public 继承。当然也有 private 和 protected 继承，不同的关键字意味着
子类对象对父类成员变量不同的访问权限。具体的权限可以查看
msdn上关于 [member access controll](https://msdn.microsoft.com/en-us/library/kktasw36.aspx "member access controll") 的描述。
这里我们只讨论 public 继承。 

在继承体系下，创建子类对象时，先调用父类构造函数，然后调用子类的构造函数；销毁子类对象时，先调用子类的析构函数，然后才调用父类的析构函数。

### 2.4 委托 + 继承 （设计模式：观察者模式）

#### 2.4.1 应用场景
先举个例子，在实现多窗口应用时，不同的窗口需要同步数据。比如 ppt的多窗口，Dota2的大地图和小地图。这里就用到了观察者模式。

![dota2](http://obi1zst3q.bkt.clouddn.com/20160807_Dota-2-review-7.jpg "dota2")

![ppt](http://obi1zst3q.bkt.clouddn.com/20160807_ppt "ppt")

首先，我们做一个抽象：

后台真正去处理和更新数据的类，我们称之为 Subject

根据数据渲染 UI 的类，我们称之为 Observer

我们将两个类同步数据的代码 抽象出来，如下所示：

``` c++
// Subject 类的设计
class Subject {
    int m_value;
    vector<Observer*> m_views;
public:
    void attach(Observer* obs) {
        m_views.push_back(obs);
    }
    void set_val(int value) {
        m_value = value;
        notify();
    }
    void notify() {
        for (int i = 0;i < m_views.size(); i++) {
            m_views[i]->update(this, m_value);
        }
    }
};

// Observer 类的设计
class Observer {
public:
    virtual void update(Subject* sub, int value) = 0;
}
```

两个类的关系用 UML 图表示为：

![Oberser](http://obi1zst3q.bkt.clouddn.com/20160807_observer "oberser")

以 ppt 为例，四个窗口用四个 Observer 对象表示，后台的数据用一个 Subject 对象表示。 
如果要增加一个窗口，则创建一个 Observer 对象，并使用Subject::attach方法与 Subject 
对象建立关系。这里未列出 取消两者关系的方法。