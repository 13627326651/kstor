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

功能4：测试10个客户端同时读写的性能。

要求：使用map作为内存缓存，数据落地存储在boltdb 中

其他功能：自由发挥


要求：
命令行（client） 通过grpc 把请求发给 server , server 负责读取和写入数据。
键值存储落地使用 boltdb  https://github.com/boltdb/bolt
命令行框架使用 https://github.com/spf13/cobra
grpc   https://github.com/grpc/grpc-go


输出boltdb grpc cobra protobuf使用文档  和kstor使用文档

