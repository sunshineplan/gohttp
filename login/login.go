package login

import (
	"errors"
	"time"

	"github.com/sunshineplan/gohttp"
)

const defaultWrongRetry = 5
const defaultRetryInterval = 5 * time.Second

var ErrNilLogin = errors.New("nil login function")

var SetAgent = gohttp.SetAgent

type Session struct {
	*gohttp.Session
	login    func(*Session, interface{}) error
	retry    int
	interval time.Duration
}

func New() *Session {
	return &Session{
		Session:  gohttp.NewSession(),
		retry:    defaultWrongRetry,
		interval: defaultRetryInterval,
	}
}

func (s *Session) SetSession(session *gohttp.Session) *Session {
	s.Session = session
	return s
}

func (s *Session) SetWrongRetry(n int) *Session {
	s.retry = n
	return s
}

func (s *Session) SetRetryInterval(d time.Duration) *Session {
	s.interval = d
	return s
}

func (s *Session) SetLogin(fn func(*Session, interface{}) error) *Session {
	s.login = fn
	return s
}

func (s *Session) Login() error {
	if s.login == nil {
		return ErrNilLogin
	}
	if s.Session == nil {
		s.Session = gohttp.NewSession()
	}
	return s.login(s, nil)
}
