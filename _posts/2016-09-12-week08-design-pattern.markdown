---
layout: post
title:  "设计模式 (上)（geekband)"
date:   2016-09-12 11:59:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*本文的主题是设计模式。设计模式在软件设计中是为解决常见问题逐渐累积出来的一套设计理念，旨在提高软件的可重用性、可扩展性和可维护性。*
*关于设计模式，本文只是第一部分，会讲解C++ 中常见的25种设计模式。*

相对于设计模式而言，设计模式背后的设计原则更为通用。越通用的东西，可复用程度就越高，也更不容易过时。从长远来看，更有利于自身的成长。
看《计算机程序设计》、《Unix网络编程》这类书，比看《MFC Windows程序设计》更有价值。 如果把设计原则比做理念，那么设计模式更多是
技法（很多是面向对象语言专用，函数式编程用不了），因此本文更多的关注点放在设计原则上。

关于设计模式，仅仅凭借网易这门课程，只能略懂皮毛。随着时间流逝，里面涉及到的内容会逐渐被遗忘。唯一能够期待的是，将来的某个时候，
当个人实现一个复杂的功能时，能够想起这些设计模式的使用场景，再去查阅资料，将某个模式应用到自己的代码中。

为了将来能够便于复习，这里推荐两个网站，作为文字版，更方便回头复习和加以理解。

1. [Design Principles](http://www.oodesign.com/design-principles.html "principles")
2. [tutorialspoint](https://www.tutorialspoint.com/design_pattern/index.htm "design pattern")

## Part 1 设计原则

软件开发过程中，我们经常讨论代码的可读性、可扩展性和重构。Google 内部每一年半左右都会对已有的项目进行重构，以满足新（和未来的）需求。
当我们对项目重构时，我们的出发点是满足现在（和未来）的需求，参考的是软件设计原则，目的是提高产品的扩展性、稳定性，让开发者更幸福。

在我看来，设计原则的终极目标是让开发者不用考虑设计原则也能开发出优秀的软件。今天，我们离这个目标已经很近，比如Web开发、Android开发、
ios开发、前端开发，我们都已经有非常成熟的框架，所以培训班出身没有计算机基础的人也能上岗；即便是后台开发，我们也能在各个领域拉出一堆框架，
比如消息队列框架、分布式存储框架、数据库套装。开发人员更多的关注点放在业务逻辑的实现上。然而这类开发人员是可替代的，你的价值在于你的
不可替代。所以，你还是需要研究算法，研究程序优化，研究设计原则。

本文的重点是讨论设计原则，首先我们看一下主流的五大设计原则：

### 1、Open Close Principle  开放关闭原则

*参考链接： http://www.oodesign.com/open-close-principle.html*

通常情况下，一个新的功能被添加到应用中时，很多代码都会发生变化。聪明的设计能够保证将变化的范围最小化，比如一两个类，而不是很多个模块。
开放关闭原则可以用来保证设计和实现代码时，保证将来增加新功能时，将对现有代码的改变最小化。

一句话描述开放关闭原则：对**扩展**开放，对**修改**关闭。

假如我们要实现一个画图程序，能够画出圆形(Circle)和矩形(Rectangle)，我们可以简单地以这样一种方式实现：

``` c++
// Open-Close Principle - Bad example
 class GraphicEditor {
 
 	public void drawShape(Shape s) {
 		if (s.m_type==1)
 			drawRectangle(s);
 		else if (s.m_type==2)
 			drawCircle(s);
 	}
 	public void drawCircle(Circle r) {....}
 	public void drawRectangle(Rectangle r) {....}
 }
 
 class Shape {
 	int m_type;
 }
 
 class Rectangle extends Shape {
 	Rectangle() {
 		super.m_type=1;
 	}
 }
 
 class Circle extends Shape {
 	Circle() {
 		super.m_type=2;
 	}
 } 
```

注意：这里我们使用 drawCircle 和 drawRectangle 两个函数分别画 圆形和矩形。如果将来要画三角形，我们要对代码进行**修改**，违背了开放关闭原则。
我们再看一看符合开放关闭原则的实现：

``` c++
// Open-Close Principle - Good example
 class GraphicEditor {
 	public void drawShape(Shape s) {
 		s.draw();
 	}
 }
 
 class Shape {
 	abstract void draw();
 }
 
 class Rectangle extends Shape  {
 	public void draw() {
 		// draw the rectangle
 	}
 } 
```

这段代码复合开放关闭原则，如果要增加一个新的功能，它会有下面几个优势：

1. 对新代码不需要重新写单元测试，使用以前的单元测试即可
2. 不需要理解 GraphEditor 的实现细节
3. 由于新功能在 新的 Shape 子类中实现，因此不需要担心牵连到旧代码

### 2、Dependency Inversion Principle 依赖倒置原则

*参考链接： http://www.oodesign.com/dependency-inversion-principle.html*

在一个软件中，通常会包含一些底层功能（操作）和一些高层功能（操作）。首先我们举几个例子：

1. mysql、postgres、sqlite 接口（低层）vs 业务逻辑层的数据读取和写入（高层） （go语言中通过 database/sql 实现抽象）
2. yaml、json、ini、csv 等文件类型接口（低层） vs 业务逻辑中数据读取和写入（高层）
3. usb 接口支持的多种设备驱动 （低层）    vs 操作系统
4. dota中一百多个英雄的技能指标（低层）    vs 一局对战中技能释放及伤害计算（高层）

在这类关系中，按照正常思路，高层依赖低层，但低层通常会发生改变。比如增加对oracle数据库的支持、dota中增加了孙悟空或某个英雄的技能发生了变化。
应对这种变化，我们需要依赖倒置原则。我们通过一个例子看一个常规的实现：

这个例子的场景是，一个 Manager 管理多个 Worker，后来我们需要增加一个 SuperWorker 让 Manager 去管理，然后代码写成了这样：

``` c++
// Dependency Inversion Principle - Bad example
class Worker {
	public void work() {
		// ....working
	}
};

class Manager {
	Worker worker;
	public void setWorker(Worker w) {
		worker = w;
	}
	public void manage() {
		worker.work();
	}
};

class SuperWorker {
	public void work() {
		//.... working much more
	}
};
```

为了让 Manager 去管理 SuperWorker，我们需要修改 Manager 的代码，因此产生下面三个问题：

1. 我们必须改变 Manager 类（设想 Manager 类已经包含了很多功能，代码可能已经超过500行）
2. Manager 类中现有的功能可能收到影响（很难知道哪些有影响，滋生bug）
3. 之前写的单元测试还要重写

使用依赖倒置原则后，我们让 Manager 依赖一个抽象接口 IWorker，因此可以这样写：

``` c++
// Dependency Inversion Principle - Good example
interface IWorker {
	public void work();
}

class Worker implements IWorker{
	public void work() {
		// ....working
	}
}

class SuperWorker  implements IWorker{
	public void work() {
		//.... working much more
	}
}

class Manager {
	IWorker worker;

	public void setWorker(IWorker w) {
		worker = w;
	}

	public void manage() {
		worker.work();
	}
}
```

Manager 类依赖抽象接口，而不是低层的具体实现以后，有三个优点：

1. 添加 SuperWorker 时，我们不需要对 Manager 做任何改变
2. 新增功能时，最小化 Manager 类中已有功能受到的影响
3. 不需要为 Manager 类重写单元测试


### 3、Interface Segregation Principle 接口分离原则
*参考链接： http://www.oodesign.com/interface-segregation-principle.html*

未完待续（请耐心等待）

### 4、Single Responsibility Principle 单一职责原则
*参考链接： http://www.oodesign.com/single-responsibility-principle.html*

未完待续（请耐心等待）

### 5、Liskov's Substitution Principle(LSP) 里氏替换原则
*参考链接：http://www.oodesign.com/liskov-s-substitution-principle.html*

未完待续（请耐心等待）