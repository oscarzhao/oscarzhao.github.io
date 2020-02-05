---
layout: post
title: MacOS 安装 JupyterNotebook 和 Spark 插件
date: 2020-02-04 17:11:00 +0800
categories: Apache Spark Jupyter
---

# 一、安装多个版本 Python

## 1.1 安装 pyenv

pyenv 是一个 python 版本管理工具，github链接是： https://github.com/pyenv/pyenv。


```bash
brew update
brew install pyenv
```

[Advanced Configuration](https://github.com/pyenv/pyenv#advanced-configuration) 最好也配置一下：

```bash
cat >> ~/.bash_profile >>EOF
export PATH=$(pyenv root)/shims:$PATH
EOF

source $(pyenv root)/completions/pyenv.bash
```

## 1.2 安装 python 3

Mac 自带 Python 2.7，无法安装 Jupyter。Python 3.6.6 兼容 Spark 和 Jupyter 比较好，所以选择这个版本。

```bash
pyenv install 3.6.6
pyenv global 3.6.6
```

如果需要在某些地方使用 2.7，可以在特定目录执行 `pyenv local 2.7.16`

# 二、Apache Spark

## 1.1 安装

为了与生产环境兼容，Apache Spark 选择的版本是 2.3.2

```bash
brew install apache-spark@2.3.2
```

安装位置：/usr/local/Cellar/apache-spark/2.3.2

把 pyspark 和依赖的库 加入到 `PYTHONPATH`：

```bash
cat >> ~/.bash_profile <<EOF
export SPARK_HOME=/usr/local/Cellar/apache-spark/2.3.2/libexec
export PYTHONPATH=$SPARK_HOME/python/:$PYTHONPATH
export PYTHONPATH=$SPARK_HOME/python/lib/py4j-0.10.7-src.zip:$PYTHONPATH
EOF
```

## 1.2 启动 Spark Standalone 集群

这里我们只列出需要的命令，具体参考[官方文档](https://spark.apache.org/docs/2.3.2/spark-standalone.html)。

集群操作的命令均存放在 `$SPARK_HOME/sbin/` 下，我只列出用到的一些

```bash
export SPARK_MASTER_HOST=127.0.0.1
$SPARK_HOME/sbin/start-master.sh
ps aux|grep master
```

启动以后，通过 `ps aux|grep master` 查看 master-spark-URL 和 web url，大概是这个样子：

```
zhaoshuaihu    52994   0.0  0.0  4267932    672 s003  S+    5:43PM   0:00.00 grep spark
zhaoshuaihu    52796   0.0  1.1  6810632 188768 s003  S     5:42PM   0:03.71 /Library/Java/JavaVirtualMachines/jdk1.8.0_191.jdk/Contents/Home/bin/java -cp /usr/local/Cellar/apache-spark/2.3.2/libexec/conf/:/usr/local/Cellar/apache-spark/2.3.2/libexec/jars/* -Xmx1g org.apache.spark.deploy.master.Master --host 127.0.0.1 --port 7077 --webui-port 8080
```

可以看到，master-spark-URL 是 spark://127.0.0.1:7077， web UI URL 是 http://localhost:8080。

```
cd $SPARK_HOME/conf
cp slaves.template slaves
$SPARK_HOME/sbin/start-slave.sh 127.0.0.1:7077
```

启动成功之后，在浏览器打开 http://localhost:8080，就可以看到下图：

![Spark Web UI](../assets/2020-02-04/spark_web_ui.jpg)

# 三、Jupyter Notebook

## 3.1 安装 Jupyter notebook

由于 notebook 6.x 和 spark 插件不兼容，所以选择了 5.7.4。

```
pip install jupyter
pip install notebook==5.7.4
pip install ipykernel==4.9.0
```

Jupyter 支持很多插件，比如 

1. [sparkmagic](https://github.com/jupyter-incubator/sparkmagic)
2. [Apache Toree](https://github.com/apache/incubator-toree)
3. [gophernotes](https://github.com/gopherdata/gophernotes)
4. 更多查看[插件清单](https://github.com/jupyter/jupyter/wiki/Jupyter-kernels)


## 3.2 安装 jupyter-spark

[jupyter-spark](https://github.com/mozilla/jupyter-spark)

```
pip install jupyter-spark
jupyter serverextension enable --py jupyter_spark
jupyter nbextension install --py jupyter_spark
jupyter nbextension enable --py jupyter_spark
jupyter nbextension enable --py widgetsnbextension
```

## 3.3 安装 sparkmagic

[sparkmagic 使用指南完整版](https://github.com/jupyter-incubator/sparkmagic)

```bash
pip install sparkmagic
jupyter nbextension enable --py --sys-prefix widgetsnbextension 
```

## 3.3.1 Livy server

参考文档：[Livy & Jupyter Notebook & Sparkmagic = Powerful & Easy Notebook for Data Scientist](https://blog.chezo.uno/livy-jupyter-notebook-sparkmagic-powerful-easy-notebook-for-data-scientist-a8b72345ea2d)

[Livy官网](http://livy.incubator.apache.org./) 下载安装包，配置&启动 Livy Server

```bash
cd ~/python3/jupyter/
wget https://www-us.apache.org/dist/incubator/livy/0.7.0-incubating/apache-livy-0.7.0-incubating-bin.zip
unzip apache-livy-0.7.0-incubating-bin.zip
cd apache-livy-0.7.0-incubating-bin

# 配置 livy-server
ls conf
# 如何修改配置，让 livy 连接上 Spark

# 启动 livy-server
./bin/livy-server 
```

启动后通过下面命令检查是否启动正常：

```bash
curl -X http://localhost:8998/sessions
```

### 3.3.2 配置 sparkmagic

设置 sparkmagic 连接 livy-server

[配置文件example](https://github.com/jupyter-incubator/sparkmagic/blob/master/sparkmagic/example_config.json)


```bash
curl -X GET https://raw.githubusercontent.com/jupyter-incubator/sparkmagic/master/sparkmagic/example_config.json > ~/.sparkmagic/config.json
```

### 3.3.3 检查 sparkmagic 是否能正常运行

1. 启动 jupyter notebook

```bash
jupyter notebook
```

2. 创建一个 pyspark notebook，加入一句话

```
spark?
```

执行，发现跪了。。。  问题排查中


# 四、启动 Jupyter Notebook

# References

1. [cat to file >>EOF](https://stackoverflow.com/questions/2500436/how-does-cat-eof-work-in-bash)

