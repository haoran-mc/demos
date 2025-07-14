package session_store

import (
	"errors"
	"net/http"

	"github.com/haoran-mc/demos/go-session/sessions"
)

var filesystemStore *sessions.FilesystemStore

var filesystemSessionName = "wx_login"

// 传入一个符串是用于 session 的认证加密
func InitFilesystemStore() {
	filesystemStore = sessions.NewFilesystemStore(
		".",
		[]byte("secret123456"),
	)
	filesystemStore.Options = &sessions.Options{
		MaxAge: 60 * 60 * 24, // 24 小时，一天
	}
}

func GetFileSystemSession(r *http.Request) (*sessions.Session, error) {
	session, err := filesystemStore.Get(r, filesystemSessionName)
	if err != nil {
		return nil, errors.New("get error")
	}
	return session, nil
}
