package lock

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"fmt"
)

const DefaultTimeout = 10 * time.Minute

//lua script to unlock the redis key
var unlockScript = redis.NewScript(1, `
if redis.call("get",KEYS[1]) == ARGV[1]
then
	return redis.call("del",KEYS[1])
else
	return 0
end
`)

type Lock struct {
	resource string
	token    string
	conn     redis.Conn //这是redis的一个链接
	timeout  time.Duration
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *Lock) tryLock() (ok bool, err error) {
	//尝试去上锁
	status, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int64(lock.timeout/time.Second), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return status == "OK", nil
}

//解锁
func (lock *Lock) UnLock() (err error) {
	//lock的Key
	_, err = unlockScript.Do(lock.conn, lock.key(), lock.token)
	return
}

//尝试对某个key进行加锁
func TryLock(conn redis.Conn, resource string) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resource, DefaultTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, timeout time.Duration) (lock *Lock, ok bool, err error) {
	//初始化Lock结构体
	lock = &Lock{resource, uuid.New(), conn, timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}
