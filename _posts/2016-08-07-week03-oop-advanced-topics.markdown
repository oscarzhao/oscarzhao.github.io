---
layout: post
title:  "C++ 对象模型"
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

使用 explicit 修饰构造函数以后，上面的默认转换就会失败。

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
下面这个是 C++11 中 std::less 的实现：

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

#### 1.5 命名空间 (namespace)
命名空间用于 模块分离和解耦。为了更好地说明一些细节，这里使用从 [msdn](https://msdn.microsoft.com/en-us/library/5cb46ksf.aspx "msdn") 摘取一段话：


``A namespace is a declarative region that provides a scope to the identifiers (the names of types, functions, variables, etc) inside it.`` 

``Namespaces are used to organize code into logical groups and to prevent name collisions that can occur especially when your code base includes multiple libraries.`` 

``All identifiers at namespace scope are visible to one another without qualification.`` 

``Identifiers outside the namespace can access the members by using the fully qualified name for each identifier, for example std::vector<std::string> vec;, or else by a using Declaration for a single identifier (using std::string), or a using Directive (C++) for all the identifiers in the namespace (using namespace std;).`` 

``Code in header files should always use the fully qualified namespace name.``

## Part 2 模板 (template)

### 2.1 class template (类模板)

前面几篇博客对类模板有所涉及，这里不再赘述。C++标准库的容器都是类模板的范例，比如：

1. [std::vector](http://www.cplusplus.com/reference/vector/ "vector")
2. [std::stack](http://www.cplusplus.com/reference/stack/ "stack")
3. [std::array](http://www.cplusplus.com/reference/stack/ "array")
4. [std::map](http://www.cplusplus.com/reference/map/ "map")
5. [and so on](http://www.cplusplus.com/reference/, "etc")

### 2.2 function template (函数模板)

对于 function template ，前面几篇博客也都有所涉及。C++标准库 algorithm 分类下有将近 90 个函数模板，
这里我列出几个：

1. [std::min](http://www.cplusplus.com/reference/algorithm/min/ "min")
2. [std::max](http://www.cplusplus.com/reference/algorithm/max/ "max")
3. [std::minmax](http://www.cplusplus.com/reference/algorithm/minmax/ "minmax")
4. [std::sort](http://www.cplusplus.com/reference/algorithm/sort/ "sort")
5. [std::copy](http://www.cplusplus.com/reference/algorithm/copy/ "copy")
6. [std::for_each](http://www.cplusplus.com/reference/algorithm/for_each/ "for_each")
7. [and so on](http://www.cplusplus.com/reference/algorithm/ "algorithm")

下面我们以 std::for_each 为例，看下如何使用函数模板：

``` c++
// for_each example (来源：http://www.cplusplus.com/reference/algorithm/for_each/)
#include <iostream>     // std::cout
#include <algorithm>    // std::for_each
#include <vector>       // std::vector

void myfunction (int i) {  // function:
  std::cout << ' ' << i;
}

struct myclass {           // function object type:
  void operator() (int i) {std::cout << ' ' << i;}
} myobject;

int main () {
  std::vector<int> myvector;
  myvector.push_back(10);
  myvector.push_back(20);
  myvector.push_back(30);

  std::cout << "myvector contains:";
  for_each (myvector.begin(), myvector.end(), myfunction);
  std::cout << '\n';

  // or:
  std::cout << "myvector contains:";
  for_each (myvector.begin(), myvector.end(), myobject);
  std::cout << '\n';

  return 0;
}
```

在这个例子中，注意函数 myfunction 和 仿函数 myobject 的用法，think twice about that。

另外，使用函数模板时，不需要指定特化类型，因为编译器会根据参数进行自动推导。

### 2.3 Member method (成员模板，默认为成员函数模板)

从使用者的角度来看，成员模板 比 类模板 具有更大的自由度。由于C++强大的继承机制，成员模板也有一些使用场景。
这里以 shared_ptr 为例：

``` c++
// 定义 类模板 shared_ptr
template <typename _Tp>
class shared_ptr : pubic __shared_ptr<_Tp> {
  //... 省略代码 ...

  template <typename _Tp1>
  explicit shared_ptr(_Tp1* __p) : __shared_ptr<_TP>(__p) {}

  // ... 省略代码 ...
};

// 使用 shared_ptr 的模板构造函数
// Derived1 类是 Base1 的子类
int main() {
  Base1 *ptr = new Derived1;  // 向上转型

  shared_ptr<Base1> sptr(new Derived1);  // 支持向上转型
}
```

这个例子中，成员模板允许 shared_ptr 支持接收子类对象的指针，构造一个父类shared_ptr。

### 2.4 specialization (模板特化)

模板本身是泛化的，允许用户在使用时进行特化。所谓“特化”，其实是指 在编译器的展开。
但是模板的设计有时候不能满足所有特化类型的要求，比如 std::vector 容纳 bool 时会有问题，
所有有了 `std::vector<bool>` 的特化版本。

#### 2.4.1 模板偏特化

模板偏特化 可以分为两类： 

1. 个数上的“偏”

  例如 `std::vector<int, typename Alloc=.....>` 相对于 `std::vector<typename T, typename Alloc=......>`

2. 类型上的“偏” (由对象扩展到 指针类型)

  这里直接看一个例子：

``` c++
// 泛化版本
template <typename T>
class C {
  //... 
};

// 扩展到指针的 特化版本
template <typename T>
class C<T*> {
  //...
};

// 使用 特化版本
int main() {
  C<string> obj1;    // 正常的特化版本
  C<string*> obj2;   // 特化的指针版本
}
```

#### 2.4.2 模板模板参数（模板嵌套)
模板模板参数是指 一个模板作为另一个模板的参数存在。这样的设计在使用上更为灵活。这里直接上一个例子：

``` c++
// 定义一个类模板，它使用一个模板作为模板参数
template <typename T,
  template <typename T>
  class Container
>
class XCls {
private:
  Container<T> c;
public:
  // ...
};

// 定义Lst
template<typename T>
using Lst = list<T, allocator<T> >;     // 注意： Lst 只有一个模板参数，而 list 有两个模板参数
// 使用该模板
int main() {
  XCls<string, list> mylist1;   // 合法的定义

//XCls<string, Lst> mylist2;    // 不合法，因为 XCls 的第二个模板参数只接受一个参数（有点绕，think about it）
}
```

这个模板的灵活性在于，第二个模板参数，你可以使用 std::list, std::stack, std::vector 
等迭代器的特化版本作为参数，也就是说底层可以接入不同的“内存管理方案”（这个词相对准确）。

## Part 3：C++语言层面的相关主题

### 3.1 C++标准库概论

这里用一张图表示

![stl](http://obi1zst3q.bkt.clouddn.com/20160808_C++%20%E6%A0%87%E5%87%86%E5%BA%93%20%E5%85%A8%E5%9B%BE "stl")

### 3.2 variadic templates：模板的可变参数列表 (C++11)

模板的可变参数列表与 正常的可变参数列表是一样的，只是语法上有些特殊。
下面是一个 print 的例子：

``` c++
// 定义 print 函数
void print() {}

template<typename T, typename... Types>
void print(const T& firstArg, const Types&... args) {
  cout << firstArg << endl;
  print(args...);
}

// 使用 print 函数
int main() {
  print(7.5, "hello", bitset<16>(377),42);
}
```

另外，对于模板参数，C++ 提供了辅助函数，用来获取可变参数列表的长度，函数签名为 `size_type sizeof...(args)`。

### 3.3 auto (C++11) 

auto 允许用户不声明变量的类型，而是留给编译器去推导。它是C++11加入的语法糖，可以减少有效代码量。

关于 auto，更多细节参考 [msdn](https://msdn.microsoft.com/en-us/library/dd293667.aspx "msdn")

### 3.4 range-based for (c++11)

这是C++11 新增加的语法，可以有效减少代码量，与 auto 配合使用更佳。考虑到是否需要修改数组的值，决定是否采用引用，看代码：

``` c++
vector<double> vec;
for (auto elem: vec) {  // 按值传递，不会修改数组的值
  cout << elem << endl;  
  elem *= 3;   // 即便这样写， 也只是修改了一个副本，不会修改 vec 的值。
}

for (auto& elem: vec){  // 按引用传递
  elem *= 3;  // 使用引用会修改数组的值
}
```

更多参考 [msdn上的描述](https://msdn.microsoft.com/en-us/library/jj203382.aspx)。

### 3.5 关于 reference (一些引起误解的地方)

#### 3.5.1 reference 的特征

reference的两个特征：

1. reference类型的变量一旦 代表某个变量，就永远不能代表另一个变量
2. reference类型的变量  大小和地址 与 原对象相同 (即 sizeof 和 operator& 的返回值)

下面用侯捷老师PPT上的一段代码来说明：

``` c++
int main() {
  int x = 0;
  int* p = &x;
  int& r = x;  // r 代表 x，两者的值都是 0
  int x2 = 5;

  r = x2;      // 这一行赋值的结果是：x 的值也变成了 5
  int& r2 = r; // r2、r 都代表 x，即值都是 5
}
```

上面这个例子中，需要注意：

1. sizeof(x) == sizeof(r)
2. &x == &r

#### 3.5.2 应用场景

reference 通常用在两个地方：

1. 参数传递 (比较快)
2. 返回类型 

### 3.6 构造和析构 (时间先后)

本小节主要讲解构造和析构 在继承和组合体系下的运作机制。

#### 3.6.1 继承体系中的构造和析构

构造：由内而外。内是指Base，外指Derived
析构：由外而内。先析构Derived Class，再析构Base Class 的部分

注意：Base Class 的析构函数必须是 virtual 的，否则会报出 undefined behaviors 的错误。
下面这段代码重现了这个错误：

``` c++
// 这段代码 来源于 stackoverflow ，但是经过了大量修改
// http://stackoverflow.com/questions/461203/when-to-use-virtual-destructors

#include <iostream>

class Node {
public:
    Node() { std::cout << "Node()" << std::endl; }
    ~Node() { std::cout << "~Node()" << std::endl; }
};

class Base 
{
public:
    Base() { std::cout << "Base()" << std::endl; }
    ~Base() { std::cout << "~Base()" << std::endl; }
};

class Derived : public Base
{
public:
    Derived()  { std::cout << "Derived()"  << std::endl;  m_pNode = new Node(); }
    ~Derived() { std::cout << "~Derived()" << std::endl;  delete m_pNode; }

private:
    Node* m_pNode;
};

int main() {
    // 注意：Base的析构函数设置为 virtual
    Base *b = new Derived();
    // 使用 b
    delete b; // 结果是：调用了 Base的构造函数

    std::cout << "execute complete" << std::endl;
}


``` 
上面这段代码打印结果是：

```
Base()
Derived()
Node()
~Base() // 为什么只打印了 这个？？？
execute complete
```

注意： 在实际的测试中，代码没有报出 undefined behavior 错误。但是出现了内存泄漏 m_pNode 的内存没有被释放。
关于这段代码的解释，我联想到侯捷老师讲到的静态绑定和动态绑定，网上一张相关的ppt，
点击[C++-dynamic-binding](http://www.cs.wustl.edu/~schmidt/PDF/C++-dynamic-binding4.pdf "vtable")查看。

然后，我给 ~Base() 和 ~Derived() 都加上了 virtual (这里就不再列出代码)，结果仍然令人疑惑，结果如下：

```
Base()
Derived()
Node()
~Derived()
~Node()
~Base()    // 为什么还会打印这个？？？
execute complete
```

又查了下文档，在 msdn的文档和 C++ dynamic binding.pdf 文档中，都提到 destructor 是不可继承的（看下图）：

![destructor](http://obi1zst3q.bkt.clouddn.com/blog/cpp/20160810-oop-pure-virtual-destructor.jpg "pic")

![destructor](http://obi1zst3q.bkt.clouddn.com/blog/cpp/20160810_msdn_destructors_override.jpg "msdn")

～Base() 虽然为 virtual 函数，但其不可继承（所以总是被override），因此析构的时候，会先调用 ~Derived(), 然后调用 ~Base()。

关于继承体系下，析构的顺序，可以参考 [msdn](https://msdn.microsoft.com/en-us/library/6t4fe76c.aspx "msdn")。

本文到此为止，谢谢耐心阅读。