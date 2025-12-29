package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var skipPaths = []string{
	"/v1/payments/webhook/",
}

func AddWebhookPath(path string) {
	skipPaths = append(skipPaths, path)
}

func shouldSkipPath(path string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

func RateLimiterSkipper(c echo.Context) bool {
	return shouldSkipPath(c.Request().URL.Path)
}

func GzipSkipper(c echo.Context) bool {
	return shouldSkipPath(c.Request().URL.Path)
}

func ApplySecurityMiddlewares(e *echo.Echo) *echo.Echo {
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: GzipSkipper,
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !shouldSkipPath(c.Request().URL.Path) {
				w := c.Response().Writer
				w.Header().Set("X-DNS-Prefetch-Control", "off")
				w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
				w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
				w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
				w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
				w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")
				w.Header().Set("X-Powered-By", "Django")
				w.Header().Set("Server", "")
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			}
			return next(c)
		}
	})

	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store:   middleware.NewRateLimiterMemoryStore(5),
		Skipper: RateLimiterSkipper,
	}))
	return e
}
