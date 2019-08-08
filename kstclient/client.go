package kstclient

import (
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var address string

func InitClient(addr string){
	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	address = addr
}

//创建bucket
func CreateBucket(bn string){

	//1 grpc生成连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//2 使用指定连接生成protobuff客户端对象
	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	//3 客户端调用接口函数
	//para1: context.Context类型
	//para2: x.proto文件定义的接口参数类型指针
	//return: x.proto文件定义的接口返回类型指针
	log.Printf("remote call(CreateBucket) start\n")
	_, err = cli.CreateBucket(ctx, &kstinter.Req{BucketName: bn})
	if err != nil {
		log.Print(err)
	}else{
		fmt.Printf("success\n")
	}
}

//删除bucket
func DelBucket(bn string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	_, err = cli.DelBucket(ctx, &kstinter.Req{BucketName:bn})
	if err != nil {
		log.Print(err)
	}else{
		fmt.Printf("success\n")
	}
}

//插入key/val
func InsertKey(bn string, k string, v string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	_, err = cli.InsertKey(ctx, &kstinter.Req{BucketName:bn, Key:k, Value:v})
	if err != nil {
		log.Print(err)
	}else{
		fmt.Printf("success\n")
	}
}

//删除key
func DelKey(bn string, k string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	_, err = cli.DelKey(ctx, &kstinter.Req{BucketName:bn, Key:k})
	if err != nil {
		log.Print(err)
	}else{
		fmt.Printf("success\n")
	}
}


func GetKey(bn string, k string) string{
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	r, err := cli.GetKey(ctx, &kstinter.Req{BucketName:bn, Key:k})
	if err != nil {
		log.Print(err)
		return ""
	}else{
		fmt.Printf("[%s:%s]\n", k, r.Value)
	}
	return r.Value
}

func GetKeyWithPrefix(bn string, prefix string)(m map[string]string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	r, err := cli.GetKeyWithPrefix(ctx, &kstinter.Req{BucketName:bn, Prefix: prefix})
	if err != nil {
		log.Print(err)
		return nil
	}else{
		for k := range r.KVs {
			fmt.Printf("[%s:%s]\n", k, r.KVs[k])
		}
	}
	return r.KVs
}


func UploadFile(filename string){
	log.Println("start to upload db", filename)
	if filename == ""{
		log.Printf("upload filename is nil")
		return
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()


	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil || n < 0{
			break
		}

		buf = buf[:n]
		conn.Write(buf)
	}
	log.Printf("upload success\n")
}