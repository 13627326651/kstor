package kstclient

import (
	"../kstinter"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var address string

func output(e error, v interface{}){
	if e != nil{
		fmt.Println(e)
		fmt.Println(1)
		return
	}
	fmt.Println(v)
	fmt.Println(0)
}

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
	output(err, "")
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
	output(err, "")
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
	output(err, "")
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
	output(err, "")
}


func GetKey(bn string, k string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	r, err := cli.GetKey(ctx, &kstinter.Req{BucketName:bn, Key:k})
	var v string
	if r != nil{
		v = r.Value
	}
	output(err, v)
}


func GetKeyWithPrefix(bn string, prefix string){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := kstinter.NewKstinterClient(conn)
	ctx := context.Background()

	r, err := cli.GetKeyWithPrefix(ctx, &kstinter.Req{BucketName:bn, Prefix: prefix})
	var ks map[string]string
	if r != nil{
		ks = r.KVs
	}
	output(err, ks)
}


func UploadFile(filename string){
	if filename == ""{
		log.Printf("cannot upload file[%s]\n", filename)
		return
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()
	cli := kstinter.NewKstinterClient(conn)
	stream, err := cli.Restore(context.Background())
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF{
			break
		}
		stream.Send(&kstinter.Frame{Data:buf[:n], Size:int32(n)})
	}
	_, err = stream.CloseAndRecv()
	output(err, "")
}

func Backup(name string){
	f, err := os.OpenFile(name, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	cli := kstinter.NewKstinterClient(conn)
	stream, err := cli.BackUp(context.Background(), &kstinter.Req{})
	if err != nil{
		fmt.Println(err)
		return
	}
	for{
		frame, err := stream.Recv()
		if err == io.EOF{
			output(nil, "")
			return
		}
		f.Write(frame.Data)
	}
}

var wg = sync.WaitGroup{}

func TestGet(threads int, count int){
	wg.Add(threads)
	c := count / threads
	start := time.Now()
	for i := 0; i < threads; i++{
		go testGet(c, strconv.Itoa(i))
	}
	wg.Wait()
	elapse := time.Now().Sub(start).Seconds()
	fmt.Println("Test durate ", elapse, "s with", threads ,"groutines for", count, "get reqs.")
	fmt.Println("吞吐量:", float64(count)/elapse, "reqs/s")
	fmt.Printf("平均响应时间:%fs\n", elapse/float64(count))
}


func TestSet(threads int, count int){
	wg.Add(threads)
	c := count / threads
	start := time.Now()
	for i := 0; i < threads; i++{
		go testSet(c, strconv.Itoa(i))
	}
	wg.Wait()
	elapse := time.Now().Sub(start).Seconds()
	fmt.Println("Test durate ", elapse, "s with", threads ,"groutines for", count, "set reqs.")
	fmt.Println("吞吐量:", float64(count)/elapse, "reqs/s")
	fmt.Printf("平均响应时间:%fs\n", elapse/float64(count))
}


func testGet(count int, suffix string){
	defer wg.Done()
	k := "key_" + suffix
	for i := 0; i < count; i++{
		GetKey("mybucket", k + "_" + strconv.Itoa(i))
		time.Sleep(10)
	}
}

func testSet(count int, suffix string){
	defer wg.Done()
	k := "key_" + suffix
	for i := 0; i < count; i++{
		InsertKey("mybucket", k +"_" + strconv.Itoa(i), k)
		time.Sleep(10)
	}
}