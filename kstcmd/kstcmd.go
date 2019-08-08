package kstcmd

import (
	"../kstclient"
	"../kstserver"
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "kstcmd",
	Short:"simple [K:V] storage",
	Long:"simple [K:V] storage",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		//fmt.Printf("RootCmd exec\n");
	},
}

//server命令上下文
type servercontext struct{
	server *cobra.Command

	flag_port string
}

//bucket命令上下文
type bucketcontext struct{
	bukcket *cobra.Command
	create *cobra.Command
	delete *cobra.Command

	flag_name string
	flag_addr string
}

//key命令上下文
type keycontext struct{
	key *cobra.Command
	set *cobra.Command
	delete *cobra.Command
	get *cobra.Command

	flag_key string
	flag_value string
	flag_bucket string
	flag_addr string
	flag_prefix bool
}

type restorecontext struct{
	restore *cobra.Command

	flag_file_name string
}
//命令上下文全局变量定义
var serverctx = servercontext{}
var buckctx = bucketcontext{}
var keyctx = keycontext{}
var restorectx = restorecontext{}

const(
	FLAG_ADDR = "addr"
	FLAG_ADDR_DEFAULT = "localhost:12345"
	FLAG_ADDR_DETAIL = "set a server port.(default)--port 123456"

	FLAG_NAME = "name"
	FLAG_NAME_DEFAULT = ""
	FLAG_NAME_DETAIL = ""

	FLAG_BUCKET = "bucket"
	FLAG_BUCKET_DEFALUT = ""
	FLAG_BUCKET_DETAIL = ""

	FLAG_PORT = "port"
	FLAG_PORT_DEFAULT = "12345"
	FLAG_PORT_DETAIL = ""

	FLAG_KEY = "key"
	FLAG_KEY_DEFAULT = ""
	FLAG_KEY_DETAIL = ""

	FLAG_VALUE = "value"
	FLAG_VALUE_DEFAULT = ""
	FLAG_VALUE_DETAIL = ""

	FLAG_PREFIX = "prefix"
	FLAG_PREFIX_DEFAULT = false
	FLAG_PREFIX_DETAIL = ""

	FLAG_FILE_NAME = "filename"
	FLAG_FILE_NAME_DEFAULT =""
	FLAG_FILE_NAME_DETAIL = "restore file name"

)

func init(){
	//serverctx初始化
	serverctx.server = &cobra.Command{
		Use:"server",
		Short: "start a server",
		Long: "start a server",
		Run: func(cmd *cobra.Command, args []string) {
			kstserver.InitServer(serverctx.flag_port)
		},
	}

	//bucketctx初始化
	buckctx.bukcket = &cobra.Command{
		Use:"bucket",
		Short: "create/delete a bucket",
		Long: "create/delete a bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	buckctx.create = &cobra.Command{
		Use:"create",
		Short: "create a bucket",
		Long: "create a bucket",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("create exec")
			kstclient.InitClient(buckctx.flag_addr)
			kstclient.CreateBucket(buckctx.flag_name)
		},
	}
	buckctx.delete = &cobra.Command{
		Use:"delete",
		Short: "delete a bucket",
		Long: "delete a bucket",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("delete exec")

			kstclient.InitClient(buckctx.flag_addr)
			kstclient.DelBucket(buckctx.flag_name)
		},
	}

	//keyctx初始化
	keyctx.key = &cobra.Command{
		Use:   "key",
		Short: "operation on key",
		Long:  "operation on key",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("key exec")
			cmd.Usage()
		},
	}
	keyctx.set = &cobra.Command{
		Use:   "set",
		Short: "set a key",
		Long:  "set a key",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("set exec")
			kstclient.InitClient(keyctx.flag_addr)
			kstclient.InsertKey(keyctx.flag_bucket, keyctx.flag_key, keyctx.flag_value)
		},
	}
	keyctx.delete = &cobra.Command{
		Use:   "delete",
		Short: "delete a key",
		Long:  "delete a key",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("delete exec")
			kstclient.InitClient(keyctx.flag_addr)
			kstclient.DelKey(keyctx.flag_bucket, keyctx.flag_key)
		},
	}
	keyctx.get = &cobra.Command{
		Use:   "get",
		Short: "get a key value",
		Long:  "get a key value",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("get exec")
			kstclient.InitClient(keyctx.flag_addr)
			if keyctx.flag_prefix {
				kstclient.GetKeyWithPrefix(keyctx.flag_bucket, keyctx.flag_key)
			}else{
				kstclient.GetKey(keyctx.flag_bucket, keyctx.flag_key)
			}

		},
	}

	//restorectx初始化
	restorectx.restore = &cobra.Command{
		Use:   "restore",
		Short: "restore a specific boltdb",
		Long:  "restore a specific boltdb",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("restore exec")
			kstclient.UploadFile(restorectx.flag_file_name)
		},
	}

	serverctx.server.Flags().StringVar(&serverctx.flag_port, FLAG_PORT, FLAG_PORT_DEFAULT, FLAG_PORT_DETAIL)
	RootCmd.AddCommand(serverctx.server)

	//添加bucket命令
	RootCmd.AddCommand(buckctx.bukcket)

	//添加create/delete
	buckctx.create.Flags().StringVar(&buckctx.flag_name, FLAG_NAME, FLAG_NAME_DEFAULT, FLAG_NAME_DETAIL)
	buckctx.create.Flags().StringVar(&buckctx.flag_addr, FLAG_ADDR, FLAG_ADDR_DEFAULT, FLAG_ADDR_DETAIL)
	buckctx.create.MarkFlagRequired(FLAG_NAME)
	buckctx.bukcket.AddCommand(buckctx.create)

	buckctx.delete.Flags().StringVar(&buckctx.flag_name, FLAG_NAME, FLAG_NAME_DEFAULT, FLAG_NAME_DETAIL)
	buckctx.delete.Flags().StringVar(&buckctx.flag_addr, FLAG_ADDR, FLAG_ADDR_DEFAULT, FLAG_ADDR_DETAIL)
	buckctx.delete.MarkFlagRequired(FLAG_NAME)
	buckctx.bukcket.AddCommand(buckctx.delete)

	//添加key
	RootCmd.AddCommand(keyctx.key)

	//添加set
	keyctx.set.Flags().StringVar(&keyctx.flag_key, FLAG_KEY, FLAG_KEY_DEFAULT, FLAG_KEY_DETAIL)
	keyctx.set.Flags().StringVar(&keyctx.flag_value, FLAG_VALUE, FLAG_VALUE_DEFAULT, FLAG_VALUE_DETAIL)
	keyctx.set.Flags().StringVar(&keyctx.flag_bucket, FLAG_BUCKET, FLAG_BUCKET_DEFALUT, FLAG_BUCKET_DETAIL)
	keyctx.set.Flags().StringVar(&keyctx.flag_addr, FLAG_ADDR, FLAG_ADDR_DEFAULT, FLAG_ADDR_DETAIL)
	keyctx.set.MarkFlagRequired(FLAG_BUCKET)
	keyctx.set.MarkFlagRequired(FLAG_KEY)
	keyctx.key.AddCommand(keyctx.set)

	//添加get
	keyctx.get.Flags().StringVar(&keyctx.flag_key, FLAG_KEY, FLAG_KEY_DEFAULT, FLAG_KEY_DETAIL)
	keyctx.get.Flags().StringVar(&keyctx.flag_bucket, FLAG_BUCKET, FLAG_BUCKET_DEFALUT, FLAG_BUCKET_DETAIL)
	keyctx.get.Flags().StringVar(&keyctx.flag_addr, FLAG_ADDR, FLAG_ADDR_DEFAULT, FLAG_ADDR_DETAIL)
	keyctx.get.Flags().BoolVar(&keyctx.flag_prefix, FLAG_PREFIX, FLAG_PREFIX_DEFAULT, FLAG_PREFIX_DETAIL)
	keyctx.get.MarkFlagRequired(FLAG_BUCKET)
	keyctx.get.MarkFlagRequired(FLAG_KEY)
	keyctx.key.AddCommand(keyctx.get)

	//添加delete
	keyctx.delete.Flags().StringVar(&keyctx.flag_key, FLAG_KEY, FLAG_KEY_DEFAULT, FLAG_KEY_DETAIL)
	keyctx.delete.Flags().StringVar(&keyctx.flag_bucket, FLAG_BUCKET, FLAG_BUCKET_DEFALUT, FLAG_BUCKET_DETAIL)
	keyctx.delete.Flags().StringVar(&keyctx.flag_addr, FLAG_ADDR, FLAG_ADDR_DEFAULT, FLAG_ADDR_DETAIL)
	keyctx.delete.MarkFlagRequired(FLAG_BUCKET)
	keyctx.delete.MarkFlagRequired(FLAG_KEY)
	keyctx.key.AddCommand(keyctx.delete)

	//添加restore
	restorectx.restore.Flags().StringVar(&restorectx.flag_file_name, FLAG_FILE_NAME, FLAG_FILE_NAME_DEFAULT, FLAG_FILE_NAME_DETAIL)
	RootCmd.AddCommand(restorectx.restore)
}
