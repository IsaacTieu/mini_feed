package fanout

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"mini-feed/internal/models"
)

type FollowStore interface {
	GetFollowers(usedID string) []string
	AddFollow(follower, user string) error
}

type FeedStore interface {
	AddToFeed(userID, postID string) error
}

func StartWorker(followStore FollowStore, feedStore FeedStore, rdb *redis.Client) {
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "post_created", "user_followed")

	go func() {
		log.Println("Redis fanout worker started")
		ch := sub.Channel()

		for msg := range ch {
			switch msg.Channel {
			case "user_followed":
				var event map[string]string
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					log.Println("error unmarshaling event:", err)
					continue
				}
				if err := followStore.AddFollow(event["follower"], event["followee"]); err != nil {
					log.Println("Error syncing follow:", err)
				} else {
					log.Printf("Synced follow: %s -> %s", event["follower"], event["followee"])
				}

			case "post_created":
				var post models.Post
				if err := json.Unmarshal([]byte(msg.Payload), &post); err != nil {
					log.Println("error unmarshaling post event:", err)
					continue
				}

				followers := followStore.GetFollowers(post.UserID)
				for _, follower := range followers {
					err := feedStore.AddToFeed(follower, post.ID)
					if err != nil {
						log.Println("fanout error:", err)
					}
				}
				log.Printf("Fanned out post %s to %d followers", post.ID, len(followers))
			}
		}
	} ()
}