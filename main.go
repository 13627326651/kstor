package main

/*
实现一个简单的键值存储项目 （kstor）

服务端  kstor server --port port

功能1：create、delete  bucket
kstor bucket create  --name bucketname --addr addr
kstor bucket delete --name bucketname --addr addr

功能2：set、get、delete key
kstor key set --key key1 --value value1 --bucket  bucketname  --addr addr
kstor key get --key key1 --bucket  bucketname  --addr addr
kstor key get --key key1 --bucket  bucketname --prefix   (输出以key1开头的所有key value)   --addr addr
kstor key delete --key key1 --bucket  bucketname   --addr addr

功能3：backup , restore
kstor backup --filename xxx.db --addr addr
kstor restore --filename xxx.cb --addr

性能测试：
--threads线程数 --count请求数量
key get test --threads 1 --count 1000
Test durate  0.883370735 s with 1 groutines for 1000 get reqs.
吞吐量: 1132.0275399433513 reqs/s
平均响应时间:0.000883s

key set test --threads 1 --count 1000
Test durate  52.158495922 s with 1 groutines for 1000 set reqs.
吞吐量: 19.17233199161728 reqs/s
平均响应时间:0.052158s


要求：
使用map作为内存缓存，数据落地存储在boltdb 中
命令行（client） 通过grpc 把请求发给 server , server 负责读取和写入数据。
键值存储落地使用 boltdb  https://github.com/boltdb/bolt
命令行框架使用 https://github.com/spf13/cobra
grpc   https://github.com/grpc/grpc-go
*/

import "./kstcmd"
func main(){
	kstcmd.RootCmd.Execute()
}
