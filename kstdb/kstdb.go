package kstdb

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"strconv"
)

var	db *bolt.DB

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
func OpenKstDB(dbname string) {
	if dbname == ""{
		log.Fatalf("open kstdb[%s] error", dbname)
	}
	log.Printf("Open kstdb[%s]\n", dbname)

	CloseKstDB()

	var err error
	db, err = bolt.Open(dbname,0600, nil)
	checkErrFatal(err)
}

//关闭数据库
func CloseKstDB() {
	if db != nil {
		db.Close()
	}
	db = nil
}

//创建bucket
func CreateBucket(bucket_name string){
	if db == nil{
		log.Fatal("database is nil")
	}
	db.Update(func(tx *bolt.Tx) error{
		_, err := tx.CreateBucketIfNotExists(([]byte(bucket_name)))
		checkErrFatal(err)
		return nil
	})
}

//删除bucket
func DelBucket(bucket_name string){
	if db == nil{
		log.Fatal("database is nil")
	}

	db.Update(func(tx *bolt.Tx) error{
		err := tx.DeleteBucket([]byte(bucket_name))
		checkErr(err)
		return nil
	})
}

//插入key/val
func InsertKey(bn string, k string, v string){
	if db == nil{
		log.Fatal("database is nil")
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bn))
		err := b.Put([]byte(k), []byte(v))
		return checkErr(err)
	})
}

//删除key
func DelKey(bn string, k string){
	if db == nil{
		log.Fatal("database is nil")
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bn))
		if b == nil {
			log.Print("no bucket")
			return nil
		}
		err := b.Delete([]byte(k))
		return checkErr(err)
	})
}

func GetKey(bn string, k string) (v string){
	if db == nil{
		log.Fatal("database is nil")
	}
	db.View(func(tx *bolt.Tx) error {

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

func GetKeyWithPrefix(bn string, prefix string)(m map[string]string){
	if db == nil{
		log.Fatal("database is nil")
	}
	db.View(func(tx *bolt.Tx) error {
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


func BackUp(w http.ResponseWriter, req *http.Request)error{
	if db == nil{
		log.Fatal("database is nil")
	}
	err := db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.FormatInt(tx.Size(), 10))
		fmt.Printf("backup size:%d\n", tx.Size())
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err
}
