package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

//限流控制

type LimiterIface interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBucket(rules ...LimiterBucketRule) MethodLimiter
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capcity      int64
	Quantum      int64
}
