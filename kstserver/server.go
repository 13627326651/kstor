package kstserver

import (
	"../kstdb"
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type server struct{
	c chan []byte
	reqs []float64
}


//创建bucket
func (s *server)CreateBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("create bucket[%s]\n", req.BucketName)
	return &kstinter.Rsp{}, kstdb.CreateBucket(req.BucketName)
}


//删除bucket
func (s *server)DelBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("delete bucket[%s]\n", req.BucketName)
	return &kstinter.Rsp{}, kstdb.DeleteBucket(req.BucketName)
}

//插
func (s *server)InsertKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	start := time.Now()
	err := kstdb.Insert(req.BucketName, req.Key, req.Value)
	s.reqs = append(s.reqs, time.Now().Sub(start).Seconds())
	fmt.Printf("bucket[%s]:insert key[%s:%s] duration[%fs]\n",req.BucketName, req.Key, req.Value, s.reqs[len(s.reqs)-1])
	if len(s.reqs) % 100 == 0{
		var all float64
		for _,v := range s.reqs{
			all += v
		}
		fmt.Println("平均响应时间: ", all/float64(len(s.reqs)), "s/req")
	}

	return &kstinter.Rsp{}, err
}

//删
func (s *server)DelKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("bucket[%s]:delete key[%s]\n",req.BucketName, req.Key)
	return &kstinter.Rsp{}, kstdb.Delete(req.BucketName, req.Key)
}

//查
func (s *server)GetKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	start := time.Now()
	v, err := kstdb.Get(req.BucketName, req.Key)
	s.reqs = append(s.reqs, time.Now().Sub(start).Seconds())
	fmt.Printf("bucket[%s]:get key[%s:%s] duration[%fs]\n",req.BucketName, req.Key, v, s.reqs[len(s.reqs)-1])
	if len(s.reqs) % 100 == 0{
		var all float64
		for _, v := range s.reqs{
			all += v
		}
		fmt.Println("平均响应时间: ", all/float64(len(s.reqs)), "s/req")
		s.reqs = make([]float64, 0)
	}
	return &kstinter.Rsp{Value: v}, err
}

//按前缀查
func (s *server)GetKeyWithPrefix(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	m, err := kstdb.GetWithPrefix(req.BucketName, req.Prefix)
	return &kstinter.Rsp{KVs:m}, err
}

//数据备份
func (s *server)BackUp(req *kstinter.Req, stream kstinter.Kstinter_BackUpServer) error{
	fmt.Println("server backup")
	name, err := kstdb.Backup()
	if err != nil{
		return err
	}
	buff := make([]byte, 1024)
	f, err := os.Open(name)
	if err != nil{
		return err
	}
	for{
		n, err := f.Read(buff)
		if err == io.EOF{
			break
		}
		stream.Send(&kstinter.Frame{Data:buff[:n], Size:int32(n)})
	}
	f.Close()
	e := os.Remove(name)
	if e != nil{
		fmt.Println(e)
	}
	return nil
}

//数据恢复
func (s *server)Restore(stream kstinter.Kstinter_RestoreServer)(error){
	rand.Seed(time.Now().Unix())
	dbname := "myboltdb_" + strconv.Itoa(rand.Intn(10000))
	f, err := os.OpenFile(dbname, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777)
	if err != nil{
		return err
	}
	for{
		frame, err := stream.Recv()
		if err == io.EOF{
			stream.SendAndClose(&kstinter.Rsp{})
			break
		}
		f.Write(frame.Data)
	}
	f.Close()
	return kstdb.Restore(dbname)
}

const(
	ADDRESS = "localhost:"
	PORT = "12345"
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
	kstdb.Init("myboltdb")

	//1 生成一个grpc服务对象，提供远程调用功能
	s := grpc.NewServer()

	//2 将实现了xx.proto文件中定义的接口的对象注册到protobuff服务端
	kstinter.RegisterKstinterServer(s, &server{reqs:make([]float64, 0)})
	log.Printf("start server\n")

	//3 在指定的监听端口上启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
