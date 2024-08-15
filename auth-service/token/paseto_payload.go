package token

import "time"

type PasetoPayload struct {
	Email     string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (pp *PasetoPayload) isExpired() bool {
	return pp.ExpiredAt.After(time.Now())
}
