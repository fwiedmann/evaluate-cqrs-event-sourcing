package infra

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/evaluate-cqrs-event-sourcing/services/user-query/domain"
	"go.etcd.io/bbolt"
	"time"
)

const boltUserBucket = "users"

func NewUserRepositoryBoltDB() *UserRepositoryBoltDB {
	db, err := bbolt.Open("/data/user-query/db.bolt", 0777, bbolt.DefaultOptions)
	if err != nil {
		panic(err)
	}
	return &UserRepositoryBoltDB{db: db}
}

type UserRepositoryBoltDB struct {
	db *bbolt.DB
}

func (u *UserRepositoryBoltDB) CreateHistoryEntry(user domain.User) error {
	return u.db.Update(func(tx *bbolt.Tx) error {
		userHistoryBucket, err := tx.CreateBucketIfNotExists([]byte(user.Id))
		if err != nil {
			return err
		}
		version, err := userHistoryBucket.NextSequence()
		if err != nil {
			return err
		}

		bucketUser := BoltDBUser{
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
			Version:   version,
			Id:        user.Id,
			Name:      user.Name,
			Surname:   user.Surname,
		}

		content, err := json.Marshal(&bucketUser)
		if err != nil {
			return err
		}

		err = userHistoryBucket.Put([]byte(fmt.Sprintf("%d", version)), content)
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *UserRepositoryBoltDB) FindById(id string) (domain.User, error) {
	var user BoltDBUser
	err := u.db.View(func(tx *bbolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(boltUserBucket))
		if err != nil {
			return err
		}

		content := bucket.Get([]byte(id))
		if len(content) == 0 {
			return fmt.Errorf("could not find user for id %s", id)
		}
		if err := json.Unmarshal(content, &user); err != nil {
			return err
		}
		return nil
	})
	return domain.User{
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
	}, err

}

func (u *UserRepositoryBoltDB) Find() ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := u.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(boltUserBucket))
		if err != nil {
			return err
		}

		err = bucket.ForEach(func(_, v []byte) error {
			var user BoltDBUser
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			users = append(users, domain.User{
				Id:      user.Id,
				Name:    user.Name,
				Surname: user.Surname,
			})
			return nil
		})

		return err
	})
	return users, err
}

func (u *UserRepositoryBoltDB) Create(user domain.User) error {
	return u.db.Update(func(tx *bbolt.Tx) error {
		usersBucket, err := tx.CreateBucketIfNotExists([]byte(boltUserBucket))
		if err != nil {
			return err
		}

		bucketUser := BoltDBUser{
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
			Id:        user.Id,
			Name:      user.Name,
			Surname:   user.Surname,
		}

		content, err := json.Marshal(&bucketUser)
		if err != nil {
			return err
		}
		err = usersBucket.Put([]byte(bucketUser.Id), content)
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *UserRepositoryBoltDB) Update(user domain.User) error {
	return u.Create(user)
}

func (u *UserRepositoryBoltDB) Delete(id string) error {
	return u.db.Update(func(tx *bbolt.Tx) error {
		usersBucket, err := tx.CreateBucketIfNotExists([]byte(boltUserBucket))
		if err != nil {
			return err
		}
		return usersBucket.Delete([]byte(id))
	})
}

type BoltDBUser struct {
	Command   string `json:"command"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version   uint64 `json:"version"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
}
