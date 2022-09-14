package tools

import (
	"encoding/json"
	"sipsimclient/errors"
	"time"

	"github.com/cristalhq/jwt/v4"
)

const (
	secret = "PxSvG+Lvi&+jFUmb"
)

type TokenInfo struct {
	UserName string
	Role     string
}

func CreateToken(userName string, role string) (string, error) {
	// create a Signer (HMAC in this example)
	key := []byte(secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}

	// create claims (you can create your own, see: Example_BuildUserClaims)
	claims := &jwt.RegisteredClaims{
		Audience: []string{userName},
		ID:       RandString(),
		Subject:  role,
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	// here is token as a string
	return token.String(), nil
}

func VerifyToken(token string) (*TokenInfo, error) {
	// create a Verifier (HMAC in this example)
	key := []byte(secret)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return nil, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return nil, err
	}

	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return nil, err
	}

	// get Registered claims
	var newClaims jwt.RegisteredClaims
	err = json.Unmarshal(newToken.Claims(), &newClaims)
	if err != nil {
		return nil, err
	}

	// verify claims as you wish
	isValid := newClaims.IsValidAt(time.Now())
	if !isValid {
		return nil, errors.ErrInvalidToken
	}
	return &TokenInfo{
		UserName: newClaims.Audience[0],
		Role:     newClaims.Subject,
	}, nil
}
