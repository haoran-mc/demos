package session_store

import (
	"errors"
	"net/http"

	"github.com/haoran-mc/demos/go-session/sessions"
	// "github.com/gorilla/sessions"
)

var cookieStore *sessions.CookieStore

var cookieSessionName = "wx_login"

// 传入一个符串是用于 session 的认证加密
func InitCookieStore() {
	cookieStore = sessions.NewCookieStore(
		[]byte("secret123456"),
	)
	cookieStore.Options = &sessions.Options{
		MaxAge: 60 * 60 * 24,
	}
}

func GetCookieSession(r *http.Request) (*sessions.Session, error) {
	session, err := cookieStore.Get(r, cookieSessionName) // 这里并不唯一
	if err != nil {
		return nil, errors.New("get error")
	}
	return session, nil
}
