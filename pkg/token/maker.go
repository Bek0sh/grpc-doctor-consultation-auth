package token

import "time"

type Maker interface {
	CreateToken(id int, name string, role string, dur time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
