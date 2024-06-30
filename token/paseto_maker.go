package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	v4SymmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	symmetricKeyBytes := []byte(symmetricKey)
	key, err := paseto.V4SymmetricKeyFromBytes(symmetricKeyBytes)
	if err != nil {
		return nil, err
	}
	maker := &PasetoMaker{
		v4SymmetricKey: key,
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))
	token.SetString("username", username)

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	token.Set("id", tokenID)

	encryptedToken := token.V4Encrypt(maker.v4SymmetricKey, nil)
	return encryptedToken, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParserForValidNow()

	parsedToken, err := parser.ParseV4Local(maker.v4SymmetricKey, token, nil)
	if err != nil {
		if err.Error() == "the ValidAt time is after this token expires" {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	claimsData := parsedToken.ClaimsJSON()
	payload := &Payload{}
	err = json.Unmarshal(claimsData, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
