package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func TrackPageView(c *gin.Context, pageName string) {
	txn := newrelic.FromContext(c.Request.Context())
	if txn == nil {
		return
	}

	txn.AddAttribute("page.name", pageName)
	txn.AddAttribute("page.path", c.Request.URL.Path)

	if userRole, exists := c.Get("user_role"); exists {
		txn.AddAttribute("user.role", userRole)
	}
}

func TrackCustomEvent(c *gin.Context, eventName string, attributes map[string]interface{}) {
	txn := newrelic.FromContext(c.Request.Context())
	if txn == nil {
		return
	}

	app := txn.Application()
	if app == nil {
		return
	}

	app.RecordCustomEvent(eventName, attributes)
}

func AddTransactionAttribute(c *gin.Context, key string, value interface{}) {
	txn := newrelic.FromContext(c.Request.Context())
	if txn == nil {
		return
	}

	txn.AddAttribute(key, value)
}

func NoticeError(c *gin.Context, err error) {
	txn := newrelic.FromContext(c.Request.Context())
	if txn == nil {
		return
	}

	txn.NoticeError(err)
}
