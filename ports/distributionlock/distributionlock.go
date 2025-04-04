package distributionlock

import "context"

type IDistributionLock interface {
	Lock(ctx context.Context, key string) (isAllow bool, err error)
	Unlock(ctx context.Context, key string) error
}
