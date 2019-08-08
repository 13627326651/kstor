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
http备份到本地myboltdb.back
curl http://localhost:8888/backup > myboltdb.back

tcp恢复本地备份文件
kstor restore --filename myboltdb.back


功能4：测试10个客户端同时读写的性能。

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
