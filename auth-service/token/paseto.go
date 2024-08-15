package token

import (
	"github.com/saeidalz13/lifestyle/auth-service/internal/autherr"
	"crypto/rand"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type TokenManger interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(pasetoToken string) (PasetoPayload, error)
}

type PasetoTokenManager struct {
	pasetoV2   *paseto.V2
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func BuildPasetoTokenManager() (*PasetoTokenManager, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &PasetoTokenManager{
		pasetoV2:   paseto.NewV2(),
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (p *PasetoTokenManager) CreateToken(email string, duration time.Duration) (string, error) {
	return p.pasetoV2.Sign(
		p.privateKey,
		PasetoPayload{
			Email:     email,
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(duration)},
		nil)
}

func (p *PasetoTokenManager) VerifyToken(pasetoToken string) (PasetoPayload, error) {
	var pp PasetoPayload
	if err := p.pasetoV2.Verify(pasetoToken, p.publicKey, &pp, nil); err != nil {
		return pp, err
	}

	if pp.isExpired() {
		return pp, autherr.ErrAuthExpiredToken
	}

	return pp, nil
}

var _ TokenManger = (*PasetoTokenManager)(nil)
