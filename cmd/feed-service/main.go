package main


import (
    "log"
    "net/http"

    "github.com/dgraph-io/badger/v4"
	"github.com/redis/go-redis/v9"

    "mini-feed/internal/services/fanout"
    "mini-feed/internal/services/feed"
    badgerdb "mini-feed/internal/storage/badger"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./data/feed"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	feedStore := badgerdb.NewFeedStore(db)

	followStore := badgerdb.NewFollowStore(db)

	fanout.StartWorker(followStore, feedStore, rdb)

	handler := feed.NewHandler(feedStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/feed", handler.GetFeed)

	log.Println("Feed service running on :8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}