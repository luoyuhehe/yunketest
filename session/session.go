package session

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// ClearAll deletes all values in the session.
	ClearAll()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save(ctx *gin.Context) error
}

type session struct {
	*sessions.Session
}

func GetSession(ctx *gin.Context, store Store, name string) (Session, error) {
	s, err := store.Get(ctx.Request, name)
	if err != nil {
		return nil, errors.WithMessage(err, "获取会话失败")
	}
	return &session{s}, nil
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session.Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session.Values[key] = val
}

func (s *session) Delete(key interface{}) {
	delete(s.Session.Values, key)
}

func (s *session) ClearAll() {
	for key := range s.Session.Values {
		s.Delete(key)
	}
}

func (s *session) Options(options Options) {
	s.Session.Options = options.ToGorillaOptions()
}

func (s *session) Save(ctx *gin.Context) error {
	err := s.Session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		return errors.WithMessage(err, "session save failed")
	}
	return nil
}
