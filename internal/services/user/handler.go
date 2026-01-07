package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	store interface {
		CreateUser(string) error
		Follow(string, string) error
	}
	rdb *redis.Client
}

func NewHandler(store interface {
		CreateUser(string) error
		Follow(string, string) error
	}, rdb *redis.Client) *Handler {
		return &Handler{store: store, rdb: rdb}
	}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	h.store.CreateUser(id)
	w.Write([]byte("user created"))
}

func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("u")
	v := r.URL.Query().Get("v")

	if u == "" || v == "" {
		http.Error(w, "missing u or v", http.StatusBadRequest)
		return
	}
	h.store.Follow(u, v)

	event := map[string]string{"follower": u, "followee": v}
	data, _ := json.Marshal(event)
	h.rdb.Publish(context.Background(), "user_followed", data)


	w.Write([]byte("followed"))
}