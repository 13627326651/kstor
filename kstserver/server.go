package kstserver

import (
	"../kstdb"
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type server struct{}

type myError struct{
	reson string
}

func (e myError) Error() string {
	return e.reson
}

//打开数据库
func (s *server)OpenKstDB(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("OpenKstDB [%s]\n", req.Dbname)

	if req.Dbname == ""{
		return &kstinter.Rsp{}, myError{reson:"boltdb name is nil\n"}
	}
	kstdb.OpenKstDB(req.Dbname)
	return &kstinter.Rsp{}, nil
}


//创建bucket
func (s *server)CreateBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("create bucket[%s]\n", req.BucketName)
	kstdb.CreateBucket(req.BucketName)
	return &kstinter.Rsp{}, nil
}


//删除bucket
func (s *server)DelBucket(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("delete bucket[%s]\n", req.BucketName)
	kstdb.DelBucket(req.BucketName)
	return &kstinter.Rsp{}, nil
}


//插
func (s *server)InsertKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("bucket[%s]:insert key[%s:%s]\n",req.BucketName, req.Key, req.Value)
	kstdb.InsertKey(req.BucketName, req.Key, req.Value)
	return &kstinter.Rsp{}, nil
}


//删
func (s *server)DelKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	fmt.Printf("bucket[%s]:delete key[%s]\n",req.BucketName, req.Key)
	kstdb.DelKey(req.BucketName, req.Key)
	return &kstinter.Rsp{}, nil
}

//查
func (s *server)GetKey(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	v := kstdb.GetKey(req.BucketName, req.Key)
	fmt.Printf("bucket[%s]:get key[%s:%s]\n",req.BucketName, req.Key, v)
	return &kstinter.Rsp{Value: v}, nil
}

//按前缀查
func (s *server)GetKeyWithPrefix(ctx context.Context, req *kstinter.Req)( *kstinter.Rsp, error){
	m := kstdb.GetKeyWithPrefix(req.BucketName, req.Prefix)
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


type HttpHandler struct{
}

var http_handler = HttpHandler{}

func (HttpHandler)BackUp(w http.ResponseWriter, r *http.Request){
	log.Print("backup request")
	err := kstdb.BackUp(w, r)
	if err != nil{
		log.Print(err)
	}
}

func initHttpServer(){
	log.Print("http server for backup is listening at localhost:8888")

	http.HandleFunc("/backup", http_handler.BackUp)

	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Print(err)
	}
}


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

	kstdb.OpenKstDB(DBNAME)

	//处理http备份请求
	go initHttpServer()
	//处理tcp恢复请求
	go initFileServer()

	//1 生成一个grpc服务对象，提供远程调用功能
	s := grpc.NewServer()

	//2 将实现了xx.proto文件中定义的接口的对象注册到protobuff服务端
	kstinter.RegisterKstinterServer(s, &server{})
	log.Printf("start server\n")
	//3 在指定的监听端口上启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	template.ParseFiles()
}


func initFileServer() {
	log.Println("initialize file server for restoring is listening at localhost:8080")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for{

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("accept new connection %s\n", conn.RemoteAddr().String())
		rand.Seed(time.Now().Unix())
		newFile := "myboltdb_restore" + strconv.Itoa(rand.Intn(10000))
		f, err := os.OpenFile(newFile, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		buff := make([]byte, 1024)
		done := false
		for{
			n, err := conn.Read(buff)
			if err == nil && n > 0 {
				f.Write(buff)
			}

			if err == io.EOF{
				fmt.Println("recv done")
				done = true
				break
			}
		}

		if done{
			kstdb.OpenKstDB(newFile)
		}

		f.Close()
		conn.Close()
	}
}

