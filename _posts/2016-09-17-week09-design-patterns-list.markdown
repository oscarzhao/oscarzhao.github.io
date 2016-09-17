---
layout: post
title:  "设计模式 (下)（geekband)"
date:   2016-09-17 11:59:00 +0800
categories: C++
---

*这是C++语言的一系列文章，内容取自于网易微专业《C++开发工程师（升级版）》。*

*本文的主题是设计模式。设计模式在软件设计中是为解决常见问题逐渐累积出来的一套设计理念，旨在提高软件的可重用性、可扩展性和可维护性。*
*关于设计模式，本文只是第一部分，会讲解C++ 中常见的25种设计模式。*

上篇文章中，我们讲了五种通用的设计原则，本文要讲的是设计模式，粒度上更细一些。本文引用了 oodesign 上的主要内容，
也将设计模式分为三大类：Creational Patterns(对象创建)、Behavioral Patterns(对象行为)和 Structural Patterns(程序结构)。
下面我们对这三类设计模式分别进行解读。

引用：[Object Oriented Design](http://www.oodesign.com/ "oodeisgn")

## Part 1 Creational Patterns 对象创建

Creational Patterns 包含七个设计模式，七个模式都是关于如何优雅地、可扩展地创建对象。它们分别是：

1. [Object Pool Pattern 对象池模式](http://www.oodesign.com/object-pool-pattern.html "object pool")
2. [Prototype Pattern 原型模式](http://www.oodesign.com/object-pool-pattern.html "prototype")
3. [Factory Method Pattern 工厂方法模式](http://www.oodesign.com/factory-method-pattern.html "factory method")
4. [Builder Pattern 建造者模式](http://www.oodesign.com/builder-pattern.html "builder")
5. [Factory Pattern 工厂模式](http://www.oodesign.com/factory-pattern.html "factory")
6. [Abstract Factory Pattern 抽象工厂模式](http://www.oodesign.com/abstract-factory-pattern.html "abstract factory")
7. [Singleton Pattern  单例模式](http://www.oodesign.com/singleton-pattern.html "singleton")

## Part 2 Behavioral Pattern 对象行为

Behavioral Patterns 包含十一个设计模式，它们通过不同的方式定义了类的行为。它们分别是：

1. [Observer Pattern 观察者模式](http://www.oodesign.com/observer-pattern.html "observer")
2. [Command Pattern ](http://www.oodesign.com/command-pattern.html "command")
3. [Strategy Pattern 策略模式](http://www.oodesign.com/strategy-pattern.html "strategy")
4. [Visitor Pattern 访问者模式](http://www.oodesign.com/strategy-pattern.html "visitor")
5. [Chain of Responsibility Pattern 职责链模式](http://www.oodesign.com/chain-of-responsibility-pattern.html "chain of responsibility")
6. [Mediator Pattern ](http://www.oodesign.com/mediator-pattern.html "mediator pattern")
7. [Iterator Pattern 迭代器模式](http://www.oodesign.com/iterator-pattern.html "iterator")
8. [Template Method Pattern 模板方法模式](http://www.oodesign.com/template-method-pattern.html "template method")
9. [Memento Pattern](http://www.oodesign.com/memento-pattern.html "memento")
10. [Interpreter Pattern](http://www.oodesign.com/interpreter-pattern.html "interpreter")
11. [Null Object Pattern](http://www.oodesign.com/null-object-pattern.html "null object")

## Part 3 Structural Patterns 

Structural Patterns 包含了六个设计模式，它们定义了类与类（接口）之间的关系。它们分别是：

1. [Adapter 适配器模式](http://www.oodesign.com/adapter-pattern.html "adapter")
2. [Bridge 桥接模式](http://www.oodesign.com/bridge-pattern.html "bridge")
3. [Composite 组合模式](http://www.oodesign.com/composite-pattern.html "composite")
4. [Decorator 装饰器模式](http://www.oodesign.com/decorator-pattern.html "decorator")
5. [Flyweight](http://www.oodesign.com/flyweight-pattern.html "fly weight")
6. [Proxy](http://www.oodesign.com/proxy-pattern.html "proxy")