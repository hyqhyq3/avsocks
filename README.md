#README

##起源
前几天受到V2EX的[clowwindy](http://www.v2ex.com/t/32777)启发，想用GO实现一个这样的加密代理，就当是练手。

##优点
- 。。。。（你懂的）
- 效率（？？还没有测试，不过据说go的效率还行）
- 保密性：使用了AES加密，默认为128位

##缺点
- 需要在客户端和服务器端都部署
- 需要有国外服务器才能。。。。（你懂的）

##安装

windows客户端可以和linux服务器端通信

###windows用户
可以直接在Download页面下载已编译版本，包含server和client，配置文件config.ini。
服务器运行server.exe,客户端运行client.exe


###linux/mac用户
首先需要有**Google go**环境

	git clone http://github.com/hyqhyq3/newsocks
	cd newsocks/client  ;;;;;如果是服务器，那么进入newsocks/server目录
	go get github.com/kless/goconfig/config
	go build
	nohup ./client &	;;;;;如果是服务器，运行nohup ./server &

##重要提示
配置文件最后一行必须要是空行，否则最后一行不会被解析

##PS
配置文件里面的那个ip是我的龟爬服务器，可以直接用