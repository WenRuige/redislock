package main

import (
	lock "github.com/redis/lock"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
	"fmt"
)

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	lock, ok, err := lock.TryLock(conn, "user:123")
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		// User is in use - return to avoid duplicate work, race conditions, etc.
		fmt.Println("this is fucking locking")
		return
	}
	defer lock.UnLock()

	time.Sleep(time.Second * 5)


	// Do something with the user.
}
