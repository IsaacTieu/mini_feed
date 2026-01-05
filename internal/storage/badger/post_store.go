package badgerdb

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"mini-feed/internal/models"
)

type PostStore struct {
	db *badger.DB
}

func NewPostStore(db *badger.DB) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) CreatePost(post *models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("post"+post.ID), data)
	})
}