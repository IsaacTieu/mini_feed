package badgerdb

import "github.com/dgraph-io/badger/v4"

type UserStore struct {
	db *badger.DB
}

func NewUserStore(db *badger.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("user:"+id), []byte{})
	})
}

func (s *UserStore) Follow(user, target string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := []byte("follow:" + user + ":" + target)
		return txn.Set(key, []byte{})
	})
}