package token_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"testing"
	queue "wall/pkg/queue/mocks"
	"wall/pkg/token"
	"wall/pkg/token/mocks"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var tokenService token.Service
var mockPublicKey rsa.PublicKey

func init() {
	mockTokenRepo := &mocks.Repository{}
	// mockDeveloperRepo := &dmocks.Repository{}
	mockQueue := &queue.Service{}
	priv, _ := rsa.GenerateKey(rand.Reader, 4096)
	mockPublicKey = priv.PublicKey
	tokenService = token.NewTokenService(mockTokenRepo, priv, mockQueue)
}

func TestIssue(t *testing.T) {

	t.Run("verify_token", func(t *testing.T) {
		tokenString, err := tokenService.Issue("5616712", "esundev", []string{"imf"})
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		jwtToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return publicKeyToBytes(&mockPublicKey), nil
		})
		//assert.NoError(t, err)
		tokenMap := jwtToken.Claims.(jwt.MapClaims)
		assert.Equal(t, tokenMap["aud"].(string), "esundev")
		assert.Equal(t, tokenMap["sub"].(string), "5616712")

	})
}

func TestRevoke(t *testing.T) {
	t.Run("revoke", func(t *testing.T) {
		err := tokenService.Revoke("j1hfa2t1", "esundev")
		assert.EqualError(t, err, "permission deny")
	})
}

func publicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, _ := x509.MarshalPKIXPublicKey(pub)

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RS256",
		Bytes: pubASN1,
	})

	return pubBytes
}

func BenchmarkIssueToken(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := "iter" + strconv.Itoa(i)
		tokenService.Issue(userID, "esundev", []string{"imf"})
	}
}
