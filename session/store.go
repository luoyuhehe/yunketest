package session

import (
	"github.com/gorilla/sessions"
)

const (
	RedisStore = "redis"
	FileStore  = "file"
	MysqlStore = "mysql"
)

// Store is an interface for custom session stores.
//
// See CookieStore and FilesystemStore for examples.
type Store interface {
	sessions.Store
	Close() error
}
