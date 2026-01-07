package post

import (
	"encoding/json"
	"net/http"
	"mini-feed/internal/models"
	"github.com/google/uuid"

	"context"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	store interface {
		CreatePost(*models.Post) error
	}
	rdb *redis.Client
}

func GenerateID() string {
	id := uuid.New().String()
	return id
}

func NewHandler(store interface {
		CreatePost(*models.Post) error
	}, rdb *redis.Client) *Handler {
		return &Handler{store: store, rdb: rdb}
	}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	text := r.URL.Query().Get("text")

	if userID == "" || text == "" {
		http.Error(w, "missing user or text", http.StatusBadRequest)
		return
	}

	post := &models.Post{
		ID: GenerateID(),
		UserID: userID,
		Text: text,
	}

	if err := h.store.CreatePost(post); err != nil {
		http.Error(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(post)
	h.rdb.Publish(context.Background(), "post_created", data)

	json.NewEncoder(w).Encode(post)

}