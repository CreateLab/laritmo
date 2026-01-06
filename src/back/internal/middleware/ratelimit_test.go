package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows requests within limit", func(t *testing.T) {
		resetRateLimiter()

		router := gin.New()
		router.POST("/test", RateLimitMiddleware(60, 10), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("blocks requests exceeding limit", func(t *testing.T) {
		resetRateLimiter()

		router := gin.New()
		router.POST("/test", RateLimitMiddleware(5, 5), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12346"

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "Too many login attempts")
	})

	t.Run("different IPs have separate limits", func(t *testing.T) {
		resetRateLimiter()

		router := gin.New()
		router.POST("/test", RateLimitMiddleware(5, 5), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req1 := httptest.NewRequest("POST", "/test", nil)
		req1.RemoteAddr = "127.0.0.1:12347"

		req2 := httptest.NewRequest("POST", "/test", nil)
		req2.RemoteAddr = "127.0.0.2:12347"

		for i := 0; i < 5; i++ {
			w1 := httptest.NewRecorder()
			router.ServeHTTP(w1, req1)
			assert.Equal(t, http.StatusOK, w1.Code)

			w2 := httptest.NewRecorder()
			router.ServeHTTP(w2, req2)
			assert.Equal(t, http.StatusOK, w2.Code)
		}
	})

	t.Run("respects X-Forwarded-For header", func(t *testing.T) {
		resetRateLimiter()

		router := gin.New()
		router.POST("/test", RateLimitMiddleware(5, 5), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.100")
		req.RemoteAddr = "127.0.0.1:12348"

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("respects X-Real-IP header", func(t *testing.T) {
		resetRateLimiter()

		router := gin.New()
		router.POST("/test", RateLimitMiddleware(5, 5), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("X-Real-IP", "10.0.0.1")
		req.RemoteAddr = "127.0.0.1:12349"

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})
}

func TestGetClientIP(t *testing.T) {
	t.Run("returns X-Forwarded-For when present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("X-Forwarded-For", "192.168.1.1")

		ip := getClientIP(c)
		assert.Equal(t, "192.168.1.1", ip)
	})

	t.Run("returns X-Real-IP when X-Forwarded-For absent", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("X-Real-IP", "10.0.0.1")

		ip := getClientIP(c)
		assert.Equal(t, "10.0.0.1", ip)
	})

	t.Run("returns ClientIP when headers absent", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "127.0.0.1:12345"

		ip := getClientIP(c)
		require.NotEmpty(t, ip)
	})
}

func resetRateLimiter() {
	loginLimiter = nil
	once = sync.Once{}
}
