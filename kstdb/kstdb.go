package kstdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strings"
	"sync"
)


type record struct{
	bucket_name string
	key string
	value string
	myfunc func(bn, k, v string) error
}

type dbcontext struct{
	db *bolt.DB
	c chan record
	mutex sync.Mutex
}


var dbctx = dbcontext{}

func (dbctx *dbcontext)openKstDB(dbname string) (err error){
	dbctx.db, err = bolt.Open(dbname,0600, nil)
	if err == nil {
		dbctx.c = make(chan record)
		go updateBoltdb(dbctx.c)
	}
	return
}

//写数据库子线程
func updateBoltdb(c chan record){
	for r := range c{
		//start := time.Now()
		r.myfunc(r.bucket_name, r.key, r.value)
		//duration := time.Now().Sub(start).Seconds()
		//fmt.Printf("updateBoltdb[%f]s\n", duration)
	}
}

var myCache map[string] map[string]string

//同步数据到缓存
func loadFromDB(){
	myCache = make(map[string] map[string]string)
	dbctx.db.View(func(tx *bolt.Tx) error {
		//遍历bucket
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			ks := make(map[string]string)
			//遍历key
			b.ForEach(func(k, v []byte) error {
				ks[string(k)] = string(v)
				return nil
			})
			myCache[string(name)] = ks
			return nil
		})
		return nil
	})
}


func Insert(bn, k, v string)error{
	if bn == "" || k == ""{
		return fmt.Errorf("cannot insert key. bn[%s],k[%s],v[%s]\n", bn, k, v)
	}
	ks := myCache[bn]
	if ks == nil{
		return fmt.Errorf("no bucket[%s]\n", bn)
	}
	ks[k] = v
	dbctx.c <- record{
		bucket_name:bn,
		key:k,
		value:v,
		myfunc: func(bn, k, v string) error {
			return dbctx.db.Update(func(tx *bolt.Tx) error{
				b := tx.Bucket([]byte(bn))
				if b == nil {
					return fmt.Errorf("no bucket[%s]\n", bn)
				}
				return b.Put([]byte(k), []byte(v))
			})
		},
	}
	return nil
}

func Get(bn, k string)(string, error){
	if bn == "" || k == ""{
		return "", fmt.Errorf("cannot get key.bn[%s], k[%s]\n", bn, k)
	}
	if myCache[bn] == nil{
		return "", fmt.Errorf("no bucket[%s]\n", bn)
	}
	v := myCache[bn][k]
	return v, nil
}

func GetWithPrefix(bn, prefix string)(map[string]string, error){
	if bn == "" || prefix == ""{
		return nil, fmt.Errorf("cannot get key with prefix.bn[%s], k[%s]\n", bn, prefix)
	}
	if myCache[bn] == nil{
		return nil, fmt.Errorf("no bucket[%s]\n", bn)
	}
	rtn := make(map[string]string)
	for k, v := range myCache[bn] {
		if strings.HasPrefix(k, prefix){
			rtn[k] = v
		}
	}
	return rtn, nil
}

func Delete(bn, k string)error{
	if bn == "" || k == ""{
		return fmt.Errorf("cannot delete key.bn[%s], k[%s]\n", bn, k)
	}
	if myCache[bn] == nil{
		return fmt.Errorf("no bucket[%s]\n", bn)
	}
	delete(myCache[bn], k)
	dbctx.c <- record{
		bucket_name:bn,
		key:k,
		value:"",
		myfunc: func(bn, k, v string) error {
			return dbctx.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bn))
				if b == nil {
					return fmt.Errorf("no bucket[%s]\n", bn)
				}
				return b.Delete([]byte(k))
			})
		},
	}
	return nil
}

func CreateBucket(bn string)error{
	if bn == ""{
		return fmt.Errorf("cannot create bucket. bucketname[%s]\n", bn)
	}
	myCache[bn] = make(map[string]string)
	dbctx.c <- record{
		bucket_name:bn,
		key:"",
		value:"",
		myfunc: func(bn, k, v string) error {
			return dbctx.db.Update(func(tx *bolt.Tx) error{
				_, err := tx.CreateBucketIfNotExists(([]byte(bn)))
				return err
			})
		},
	}
	return nil
}

func DeleteBucket(bn string)error{
	if bn == ""{
		return fmt.Errorf("cannot delete bucket. bucketname[%s]\n", bn)
	}
	delete(myCache, bn)
	dbctx.c <- record{
		bucket_name:bn,
		key:"",
		value:"",
		myfunc: func(bn, k, v string) error {
			return dbctx.db.Update(func(tx *bolt.Tx) error{
				return tx.DeleteBucket([]byte(bn))
			})
		},
	}
	return nil
}

func Init(name string){
	if err := dbctx.openKstDB(name); err != nil{
		log.Fatal("open boltdb fail")
	}
	loadFromDB()
}


func Backup() (name string, err error){
	name = "myboltdb_backup"
	err = dbctx.db.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(name, os.ModePerm)
	})
	return
}

func Restore(name string) error {
	db, err := bolt.Open(name,0600, nil)
	if err != nil{
		os.Remove(name)
		return err
	}
	db.Close()
	if dbctx.db != nil{
		dbctx.db.Close()
	}
	Init(name)
	return nil
}

//
//
//func BackUp(w http.ResponseWriter, req *http.Request)error{
//	if db == nil{
//		log.Fatal("database is nil")
//	}
//	err := db.View(func(tx *bolt.Tx) error {
//		w.Header().Set("Content-Type", "application/octet-stream")
//		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
//		w.Header().Set("Content-Length", strconv.FormatInt(tx.Size(), 10))
//		fmt.Printf("backup size:%d\n", tx.Size())
//		_, err := tx.WriteTo(w)
//		return err
//	})
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//	return err
//}
