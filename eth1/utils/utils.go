package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"strings"
	"time"
)

func ToSSV(value *big.Int, format string) string {
	f := toSSV(value)
	str := strings.TrimRight(fmt.Sprintf(format, f), "0")
	str = strings.TrimRight(str, ".")
	return str
}

func toSSV(value *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).Quo(big.NewFloat(0).SetInt(value), big.NewFloat(params.GWei)), big.NewFloat(params.GWei))
}

type RetryConfig struct {
	MaxRetries int
	RetryDelay time.Duration
}

var DefaultRetryConfig = RetryConfig{
	MaxRetries: 3,
	RetryDelay: time.Second,
}

func Retry[T any](operation func() (T, error), config RetryConfig) (T, error) {
	var result T
	var err error

	for attempt := 0; attempt < config.MaxRetries; attempt++ {
		result, err = operation()
		if err == nil {
			return result, nil
		}

		if attempt < config.MaxRetries-1 {
			time.Sleep(config.RetryDelay)
		}
	}

	return result, fmt.Errorf("operation failed after %d attempts: %w", config.MaxRetries, err)
}
