package middleware

import (
	"fmt"
	"net/http"
	"openapphub/internal/model"
	"openapphub/pkg/cache"
	"openapphub/pkg/serializer"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

// RateLimiterConfig holds the configuration for rate limiting
type RateLimiterConfig struct {
	RateString  string // Rate limit string (e.g., "100-H" for 100 requests per hour)
	LimitByUser bool   // If true, limit by user ID; if false, limit by IP
}

// RateLimiter returns a Gin middleware for rate limiting.
// It can limit by IP or user ID, depending on the configuration.
func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	// Initialize rate limiter components
	rate, store, instance := setupRateLimiter(config.RateString)

	return func(c *gin.Context) {
		// Determine the identifier for rate limiting
		key := getIdentifier(c, config.LimitByUser)

		// Apply rate limiting
		var _ limiter.Store = store
		var _ limiter.Rate = rate
		err := applyRateLimit(c, instance, key)
		if err != nil {
			// Use the global logger to log the error
			GetZapLogger().Error("Rate limit error", zap.Error(err))
		}

		// If the context was aborted (due to rate limiting), don't continue
		if c.IsAborted() {
			return
		}

		// Continue processing the request
		c.Next()
	}
}

// setupRateLimiter initializes the rate limiter components
func setupRateLimiter(rateString string) (limiter.Rate, limiter.Store, *limiter.Limiter) {
	// Parse the rate limit string
	rate, err := limiter.NewRateFromFormatted(rateString)
	if err != nil {
		panic(err)
	}

	// Create a Redis store for rate limiting
	store, err := sredis.NewStoreWithOptions(cache.RedisClient, limiter.StoreOptions{
		Prefix: "limiter_", // Prefix for Redis keys
	})
	if err != nil {
		panic(err)
	}

	// Create a new rate limiter instance
	instance := limiter.New(store, rate)
	return rate, store, instance
}

// getIdentifier returns the identifier for rate limiting.
// If limitByUser is true and a user is logged in, it uses the user ID.
// Otherwise, it uses the client's IP address.
func getIdentifier(c *gin.Context, limitByUser bool) string {
	if limitByUser {
		// Get the user from the context
		user, exists := c.Get("user")
		if exists {
			// Check if the user is of type *model.User
			if u, ok := user.(*model.User); ok {
				return fmt.Sprintf("user:%d", u.ID)
			}
		}
	}
	// Use the client's IP address as the identifier
	return "ip:" + c.ClientIP()
}

// applyRateLimit applies the rate limit and sets appropriate headers
func applyRateLimit(c *gin.Context, instance *limiter.Limiter, key string) error {
	// Get the rate limit context
	context, err := instance.Get(c, key)
	if err != nil {
		response := serializer.Err(serializer.CodeInternalServerError, "Failed to get rate limit info", err)
		c.JSON(http.StatusInternalServerError, response)
		GetZapLogger().Error("Failed to get rate limit info", zap.Any("response", response))
		c.Abort()
		return err
	}

	// Check if the rate limit has been exceeded
	if context.Reached {
		response := serializer.Response{
			Code: serializer.CodeRateLimitExceeded,
			Msg:  "Rate limit exceeded",
		}
		c.JSON(http.StatusTooManyRequests, response)
		GetZapLogger().Warn("Rate limit exceeded", zap.String("key", key))
		c.Abort()
		return nil
	}

	// Set the rate limit headers
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
	c.Header("X-RateLimit-Reset", time.Unix(context.Reset, 0).Format(time.RFC3339))

	return nil
}
