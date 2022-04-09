package deprecator

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/proxy"
	krakendgin "github.com/luraproject/lura/router/gin"
)

type Rejector func(c *gin.Context) bool

func RejectorFactory(cfg *Config) Rejector {
	deprecationWindow := cfg.Complete.Sub(cfg.Start).Milliseconds()

	return func(c *gin.Context) bool {
		currentTime := time.Now()
		if currentTime.Before(cfg.Start) {
			return false
		}
		if cfg.Complete.Before(currentTime) {
			return true
		}
		if deprecationWindow == 0 && cfg.Start.Before(currentTime) {
			return true
		}
		sinceStart := currentTime.Sub(cfg.Start).Milliseconds()
		rejectRate := float64(sinceStart) / float64(deprecationWindow)
		if rand.Float64() > rejectRate { //nolint
			return false
		}
		return true
	}
}

func HandlerFactory(next krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	// runs when handler is registered
	return func(remote *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handlerFunc := next(remote, p)

		cfg, hasConfig := ConfigGetter(remote.ExtraConfig)
		if !hasConfig {
			return handlerFunc
		}

		shouldReject := RejectorFactory(cfg)
		startStr := cfg.Start.String()

		// runs when request is executed
		return func(c *gin.Context) {
			c.Writer.Header().Set("Deprecation", "true")
			c.Writer.Header().Set("Sunset", startStr)

			if shouldReject(c) {
				for key, val := range cfg.Headers {
					c.Writer.Header().Add(key, val)
				}
				c.AbortWithStatusJSON(cfg.Status, cfg.Body)
				return
			}

			handlerFunc(c)
		}
	}
}
