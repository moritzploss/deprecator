package deprecator

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/proxy"
	krakendgin "github.com/luraproject/lura/router/gin"
)

type Rejector func(c *gin.Context) bool

func HeadsUpFactory(cfg *Config) func(t time.Time) bool {
	isHeadsUp := make([]func(t time.Time) bool, len(cfg.HeadsUp.Dates))
	for i, date := range cfg.HeadsUp.Dates {
		startHeadsUp := date
		endHeadsUp := startHeadsUp.Add(cfg.HeadsUp.Duration.Duration)
		isHeadsUp[i] = func(t time.Time) bool {
			if startHeadsUp.Before(t) && endHeadsUp.After(t) {
				return true
			}
			return false
		}
	}
	return func(t time.Time) bool {
		for i := range isHeadsUp {
			if isHeadsUp[i](t) {
				return true
			}
		}
		return false
	}
}

func RejectorFactory(cfg *Config) Rejector {
	isHeadsUp := HeadsUpFactory(cfg)

	return func(c *gin.Context) bool {
		currentTime := time.Now()
		if isHeadsUp(currentTime) {
			return true
		}
		if currentTime.Before(cfg.Deprecate) {
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
		sunset := cfg.Sunset.String()

		// runs when request is executed
		return func(c *gin.Context) {
			c.Writer.Header().Set("Deprecation", "true")
			c.Writer.Header().Set("Sunset", sunset)

			if shouldReject(c) {
				for key, val := range cfg.Response.Headers {
					c.Writer.Header().Add(key, val)
				}
				c.AbortWithStatusJSON(cfg.Response.Status, cfg.Response.Body)
				return
			}

			handlerFunc(c)
		}
	}
}
