package infra

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/command/domain"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user/core"
	"go.etcd.io/bbolt"
	"time"
)

func NewUserRepositoryBoltDB() *UserRepositoryBoltDB {
	db, err := bbolt.Open("/data/command/db.bolt", 0777, bbolt.DefaultOptions)
	if err != nil {
		panic(err)
	}
	return &UserRepositoryBoltDB{db: db}
}

type UserRepositoryBoltDB struct {
	db *bbolt.DB
}

type BoltDBUser struct {
	Command   string `json:"command"`
	CreatedAt string `json:"created_at"`
	Version   uint64 `json:"version"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
}

func (u *UserRepositoryBoltDB) Append(command string, user core.User) error {
	return u.db.Update(func(tx *bbolt.Tx) error {

		var bucket *bbolt.Bucket
		var err error
		if command == domain.CreateCommand {
			bucket, err = tx.CreateBucketIfNotExists([]byte(user.Id))
		} else {
			bucket = tx.Bucket([]byte(user.Id))
		}

		if bucket == nil {
			return fmt.Errorf("could not get bucket")
		}

		b, err := tx.CreateBucketIfNotExists([]byte(user.Id))
		if err != nil {
			return err
		}

		version, err := b.NextSequence()
		if err != nil {
			return err
		}

		bucketUser := BoltDBUser{
			Command:   command,
			CreatedAt: time.Now().String(),
			Version:   version,
			Id:        user.Id,
			Name:      user.Name,
			Surname:   user.Surname,
		}

		fmt.Printf("bolt append: %+v", bucketUser)

		content, err := json.Marshal(&bucketUser)
		if err != nil {
			return err
		}

		return b.Put([]byte(fmt.Sprintf("%d", version)), content)
	})
}
