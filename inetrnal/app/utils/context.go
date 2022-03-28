package utils

import (
	"context"
	"time"
)

func GetTimeoutContext(timeoutSec int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
}
