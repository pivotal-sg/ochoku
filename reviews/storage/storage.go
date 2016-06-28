package storage

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
)

type Storer interface {
	Get(string) (proto.ReviewDetails, error)
	List() ([]*proto.ReviewDetails, error)
	Insert(proto.ReviewDetails) error
}

type BoltStore struct {
	store *bolt.DB
}

func New(storageFile string) (*BoltStore, error) {
	boltStore := &BoltStore{}

	db, err := bolt.Open(storageFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("reviews"))
		return err
	})

	boltStore.store = db
	return boltStore, nil
}

func (bs *BoltStore) Get(name string) (rd proto.ReviewDetails, err error) {
	reviewDetails := &rd
	err = bs.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("reviews"))
		v := b.Get([]byte(name))
		err := json.Unmarshal(v, reviewDetails)
		return err
	})
	return
}

func (bs *BoltStore) List() (rl []*proto.ReviewDetails, err error) {
	rl = make([]*proto.ReviewDetails, 0, 0)
	reviewList := &rl
	err = bs.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("reviews"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			review := &proto.ReviewDetails{}
			if err := json.Unmarshal(v, review); err != nil {
				return err
			}
			*reviewList = append(*reviewList, review)
		}
		return err
	})
	return
}

func (bs *BoltStore) Insert(review proto.ReviewDetails) error {
	return bs.store.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("reviews"))
		v, err := json.Marshal(review)
		if err != nil {
			return err
		}
		return b.Put([]byte(review.Name), v)
	})
}
