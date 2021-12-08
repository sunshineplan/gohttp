package login

import (
	"fmt"
	"time"
)

func (s *Session) Kaptcha(kaptcha func(*Session) (string, error)) error {
	if s == nil {
		s = New()
	}
	if s.login == nil {
		return ErrNilLogin
	}

	data, err := kaptcha(s)
	if err != nil {
		time.Sleep(s.interval)
		if s.retry > 0 {
			s.retry--
			return s.Kaptcha(kaptcha)
		}
		return fmt.Errorf("max wrong retry: failed to get kaptcha: %s", err)
	}

	if err = s.login(s, data); err != nil {
		time.Sleep(s.interval)
		if s.retry > 0 {
			s.retry--
			return s.Kaptcha(kaptcha)
		}
		return fmt.Errorf("max wrong retry: failed to login: %s", err)
	}

	return nil
}
