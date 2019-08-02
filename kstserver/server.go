package kstserver

import (
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"../kstdb"
)

type server struct{}
var db = kstdb.KstCtx{}

//创建bucket
func (s *server)CreateBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("create bucket[%s]\n", req.BucketName)
	db.CreateBucket(req.BucketName)
	return &kstinter.Rsp{}, nil
}


//删除bucket
func (s *server)DelBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("delete bucket[%s]\n", req.BucketName)
	db.DelBucket(req.BucketName)
	return &kstinter.Rsp{}, nil
}


//插
func (s *server)InsertKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("bucket[%s]:insert key[%s:%s]\n",req.BucketName, req.Key, req.Value)
	db.InsertKey(req.BucketName, req.Key, req.Value)
	return &kstinter.Rsp{}, nil
}


//删
func (s *server)DelKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("bucket[%s]:delete key[%s]\n",req.BucketName, req.Key)
	db.DelKey(req.BucketName, req.Key)
	return &kstinter.Rsp{}, nil
}

//查
func (s *server)GetKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	v := db.GetKey(req.BucketName, req.Key)
	fmt.Printf("bucket[%s]:get key[%s:%s]\n",req.BucketName, req.Key, v)
	return &kstinter.Rsp{Value: v}, nil
}

//按前缀查
func (s *server)GetKeyWithPrefix(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	m := db.GetKeyWithPrefix(req.BucketName, req.Prefix)
	fmt.Printf("bucket[%s]:get key with prefix[%s]\n",req.BucketName, req.Prefix)
	for k := range m {
		fmt.Printf("[%s:%s]\n",k, m[k])
	}
	return &kstinter.Rsp{KVs:m}, nil
}

const(
	ADDRESS = "localhost:"
	PORT = "12345"
	DBNAME = "myboltdb"
)

//初始化服务端
func InitServer(port string){

	p, err := strconv.Atoi(port)
	if err != nil || p <= 1024 {
		port = PORT
	}

	lis, err := net.Listen("tcp", ADDRESS + port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening at %s%s\n", ADDRESS, port)

	db.OpenKstDB(DBNAME)

	//1 生成一个grpc服务对象，提供远程调用功能
	s := grpc.NewServer()

	//2 将实现了xx.proto文件中定义的接口的对象注册到protobuff服务端
	kstinter.RegisterKstinterServer(s, &server{})
	log.Printf("start server\n")
	//3 在指定的监听端口上启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


