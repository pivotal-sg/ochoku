package storage

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	proto "github.com/pivotal-sg/ochoku/userservice/proto"
)

type Storer interface {
	Get(string) (proto.UserDetails, error)
	Insert(proto.UserDetails) error
}

type BoltStore struct {
	store *bolt.DB
}

// Create a new storage, using boltdb, and the bucket 'users'
func New(storageFile string) (*BoltStore, error) {
	boltStore := &BoltStore{}

	db, err := bolt.Open(storageFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})

	boltStore.store = db
	return boltStore, nil
}

// Get a user by its username.
func (bs *BoltStore) Get(username string) (ud proto.UserDetails, err error) {
	userDetails := &ud
	err = bs.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get([]byte(username))
		err := json.Unmarshal(v, userDetails)
		return err
	})
	return
}

// Insert a user into the store; hash its password.
func (bs *BoltStore) Insert(user proto.UserDetails) error {
	return bs.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("reviews"))
		v, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.Username), v)
	})
}
