package kstclient

import (
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)


type Client struct{
	addr string
}

func (c *Client)InitClient(addr string){
	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	c.addr = addr
}

//创建bucket
func (c *Client)CreateBucket(bn string){

	//1 grpc生成连接
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
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
func (c *Client)DelBucket(bn string){
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
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
func (c *Client)InsertKey(bn string, k string, v string){
	log.Printf("1 client insert key [%s:%s]\n", k, v)

	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("2 client insert key [%s:%s]\n", k, v)

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	log.Printf("3 client insert key [%s:%s]\n", k, v)

	_, err = cli.InsertKey(ctx, &kstinter.Req{BucketName:bn, Key:k, Value:v})
	if err != nil {
		log.Print(err)
	}else{
		fmt.Printf("success\n")
	}
}

//删除key
func (c *Client)DelKey(bn string, k string){
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
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


func (c *Client)GetKey(bn string, k string) string{
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
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

func (c *Client)GetKeyWithPrefix(bn string, prefix string)(m map[string]string){
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
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


