package tools

import (
	"testing"
)

func TestJWT(t *testing.T) {
	token, err := CreateToken("myname", "role")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	issuer, err := VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(issuer)

	issuer, err = VerifyToken(token + "11111")
	if err == nil {
		t.Error("invalid token pass")
	} else {
		t.Log(err)
	}
	t.Log(issuer)
}
