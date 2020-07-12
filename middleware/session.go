package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

// Store session存储
var Store memstore.Store

const WeekSeconds = 7 * 86400

func Session(secret string) gin.HandlerFunc {
	Store = memstore.NewStore([]byte(secret))
	Store.Options(sessions.Options{HttpOnly: true, MaxAge: WeekSeconds, Path: "/"})

	return sessions.Sessions("vinki-session", Store)
}
