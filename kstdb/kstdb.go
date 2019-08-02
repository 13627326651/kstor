package kstdb

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
)

/*对外接口，调用者实例化*/
type KstCtx struct{
	db *bolt.DB
}

func checkErrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkErr(err error) error{
	if err != nil {
		log.Print(err)
	}

	return err
}

//打开数据库
func (ctx *KstCtx)OpenKstDB(dbname string) {
	if dbname == ""{
		log.Fatalf("open kstdb[%s] error", dbname)
	}
	log.Printf("Open kstdb[%s]\n", dbname)

	var err error
	ctx.db, err = bolt.Open(dbname,0600, nil)
	checkErrFatal(err)
}

//关闭数据库
func (ctx *KstCtx)CloseKstDB() {
	if ctx.db != nil {
		ctx.db.Close()
	}
	ctx.db = nil
}

//创建bucket
func (ctx *KstCtx)CreateBucket(bucket_name string){
	ctx.db.Update(func(tx *bolt.Tx) error{
		_, err := tx.CreateBucketIfNotExists(([]byte(bucket_name)))
		checkErrFatal(err)
		return nil
	})
}

//删除bucket
func (ctx *KstCtx)DelBucket(bucket_name string){
	ctx.db.Update(func(tx *bolt.Tx) error{
		err := tx.DeleteBucket([]byte(bucket_name))
		checkErr(err)
		return nil
	})
}

//插入key/val
func (ctx *KstCtx)InsertKey(bn string, k string, v string){
	ctx.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bn))
		err := b.Put([]byte(k), []byte(v))
		return checkErr(err)
	})
}

//删除key
func (ctx *KstCtx)DelKey(bn string, k string){
	ctx.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bn))
		if b == nil {
			log.Print("no bucket")
			return nil
		}
		err := b.Delete([]byte(k))
		return checkErr(err)
	})
}

func (ctx *KstCtx)GetKey(bn string, k string) (v string){
	ctx.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bn))
		if b == nil {
			log.Printf("no bucket")
			return nil
		}

		v = string(b.Get([]byte(k)))
		return nil
	})
	return
}

func (ctx *KstCtx)GetKeyWithPrefix(bn string, prefix string)(m map[string]string){

	ctx.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bn))
		if b == nil {
			log.Printf("no bucket")
			return nil
		}

		c := b.Cursor()
		m = make(map[string]string)
		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next(){
			m[string(k)] = string(v)
		}
		return nil
	})
	return
}

