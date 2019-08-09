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



要求：使用map作为内存缓存，数据落地存储在boltdb 中

其他功能：自由发挥


要求：
命令行（client） 通过grpc 把请求发给 server , server 负责读取和写入数据。
键值存储落地使用 boltdb  https://github.com/boltdb/bolt
命令行框架使用 https://github.com/spf13/cobra
grpc   https://github.com/grpc/grpc-go


输出boltdb grpc cobra protobuf使用文档  和kstor使用文档



# 如果已经安装了proto和protoc-gen-go的话就不用安装了
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
或者
*git clone https://github.com/golang/protobuf.git $GOPATH/src/github.com/golang/protobuf/


# 下载grpc-go
*git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc

# 下载golang/net
*git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net

# 下载golang/text
*git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text

# 下载go-genproto
git clone https://github.com/google/go-genproto.git $GOPATH/src/google.golang.org/genproto

# 安装
cd $GOPATH/src/./main
go install google.golang.org/grpc

#安装cobra
git clone https://github.com/spf13/cobra.git $GOPATH/src/github.com/spf13/cobra
依赖：
git clone https://github.com/inconshreveable/mousetrap.git $GOPATH/github.com/inconshreveable/mousetrap
git clone https://github.com/spf13/pflag.git $GOPATH/src/github.com/spf13/pflag


protobuff理解、教程
https://blog.csdn.net/samdy2008/article/details/52139047
https://www.ibm.com/developerworks/cn/linux/l-cn-gpb/

go中使用grpc、protobuf
https://blog.csdn.net/sureSand/article/details/82858047
https://blog.csdn.net/xuduorui/article/details/78278808
*https://blog.csdn.net/weixin_42654444/article/details/82945195

protoc -I ./ ./helloworld.proto --go_out=plugins=grpc:./
该命令需要protoc-go-gen.exe插件的支持，并加入PATH环境变量

boltdb的使用教程
https://studygolang.com/articles/12446
https://studygolang.com/articles/10433



./main bucket create --name bucket_1
./main bucket delete --name bucket_1


./main bucket create --name bucket_2
./main key set --bucket bucket_2 --key a --value 0
./main key set --bucket bucket_2 --key b --value 1
./main key set --bucket bucket_2 --key c --value 2d
./main key set --bucket bucket_2 --key aad --value 3
./main key set --bucket bucket_2 --key aaaf --value 4
./main key set --bucket bucket_2 --key abc --value 5
./main key set --bucket bucket_2 --key dd --value 6
./main key set --bucket bucket_2 --key gg --value 7

./main key get --bucket bucket_2 --key abc
./main key get --bucket bucket_2 --key a  --prefix

./main key delete --bucket bucket_2 --key abc
./main key get --bucket bucket_2 --key abc


测试指令
cd ~/gowork/kstor/; ./kstortest.sh a

linux挂载windows共享
sudo mount -t cifs //10.30.1.194/distributed-storage /home/test/storage/ -o username=admin,password=xgp123//,vers=2.0

curl http://localhost:8888/backup > myboltdb.back




