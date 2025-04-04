package cache

import "time"

type ICache interface {
	Set(key, value any, ttl time.Duration) error
	Get(key any) (value any, err error)
}
