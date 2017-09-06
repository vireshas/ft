package ft

import (
	"log"
    "github.com/boltdb/bolt"
)

type DB struct {
	instance *bolt.DB
}

var BUCKET = []byte("ft")


type errorString struct {
    s string
}

func (e *errorString) Error() string {
	return e.s
}

func NewDB(path string) DB{
    db, err := bolt.Open("./bolt.db", 0644, nil)
    if err != nil {
        log.Fatal(err)
    }

	return DB{db}
}

func (db DB) Read(key []byte) ([]byte, error) {
	var res []byte
    err := db.instance.View(func(tx *bolt.Tx) error {

        bucket := tx.Bucket(BUCKET)
        if bucket == nil {
			return &errorString{"Bucket not found!"}
        }

        res = bucket.Get(key)
		
		return nil
	})

	return res, err
}

func (db DB) Write(key []byte, value []byte) error {
    err := db.instance.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists(BUCKET)
        if err != nil {
            return err
        }

        err = bucket.Put(key, value)
        if err != nil {
            return err
        }
        return nil
    })

	return err
}
