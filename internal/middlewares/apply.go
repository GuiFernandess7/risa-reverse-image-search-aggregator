package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ApplyMiddlewares(e *echo.Echo) *echo.Echo {
	e.Use(middleware.Recover())
	//e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
			return next(c)
		}
	})

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	return e
}
