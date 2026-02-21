package net

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/labstack/echo/v4"
)

type rateLimitConfig struct {
	// Limit - max requests per Window
	Limit int64
	// Window - fixed time window (e.g. 60s)
	Window time.Duration
	// If true: include route pattern in key (e.g. /api/login has separate bucket)
	// If false: per-IP global bucket for all routes
	Prefix            string
	PerRoute          bool
	TrustProxyHeaders bool
	Skipper           func(c echo.Context) bool
	OnLimit           func(c echo.Context, retryAfter time.Duration) error
}

func (cfg *rateLimitConfig) setDefaults() {
	if cfg.Prefix == "" {
		cfg.Prefix = "ratelimit"
	}
	if cfg.Limit <= 0 {
		cfg.Limit = 20
	}
	if cfg.Window <= 0 {
		cfg.Window = 1 * time.Minute
	}
	if cfg.Skipper == nil {
		cfg.Skipper = func(c echo.Context) bool { return false }
	}

	if cfg.OnLimit == nil {
		cfg.OnLimit = func(c echo.Context, retryAfter time.Duration) error {
			c.Response().Header().Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())))
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"error":       "rate_limit_exceeded",
				"retry_after": int(retryAfter.Seconds()),
			})
		}
	}
}

func (e *Echo) RateLimitMW(cfg rateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}
			key := defaultKey(cfg, c)
			ctx := c.Request().Context()
			resp, err := e.rdsAPI.API.RateLimit(ctx, &brzrpc.RateLimitRequest{
				Key:                key,
				WindowMilliseconds: cfg.Window.Milliseconds(),
			})
			if err != nil {
				return next(c)
			}

			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Limit))
			remaining := cfg.Limit - resp.Count
			if remaining < 0 {
				remaining = 0
			}
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(resp.Ttl)*time.Millisecond).Unix()))

			if resp.Count > cfg.Limit {
				return cfg.OnLimit(c, time.Duration(resp.Ttl)*time.Millisecond)
			}

			return next(c)
		}
	}
}

func defaultKey(cfg rateLimitConfig, c echo.Context) string {
	ip := clientIP(cfg, c)

	route := ""
	if cfg.PerRoute {
		route = c.Path()
	}

	raw := ip + "|" + route
	sum := sha1.Sum([]byte(raw))
	h := hex.EncodeToString(sum[:])

	if cfg.PerRoute {
		return fmt.Sprintf("%s:%s:%s", cfg.Prefix, "ip_route", h)
	}
	return fmt.Sprintf("%s:%s:%s", cfg.Prefix, "ip", h)
}

func clientIP(cfg rateLimitConfig, c echo.Context) string {
	req := c.Request()

	if cfg.TrustProxyHeaders {
		// X-Forwarded-For: client, proxy1, proxy2...
		if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			if len(parts) > 0 {
				ip := strings.TrimSpace(parts[0])
				if parsed := net.ParseIP(ip); parsed != nil {
					return ip
				}
			}
		}
		if xri := req.Header.Get("X-Real-IP"); xri != "" {
			if parsed := net.ParseIP(strings.TrimSpace(xri)); parsed != nil {
				return strings.TrimSpace(xri)
			}
		}
	}

	// Echo has c.RealIP() (it considers headers depending on config),
	// but here we control it ourselves.
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil && net.ParseIP(host) != nil {
		return host
	}

	// fallback
	return req.RemoteAddr
}
