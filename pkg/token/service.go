package token

import (
	"crypto/rsa"
	"errors"

	// "wall/pkg/developer"
	"time"
	"wall/pkg/entity"
	"wall/pkg/queue"
	"wall/utils/rand"

	"github.com/dgrijalva/jwt-go"
)

type tokenService struct {
	tokenRepo  Repository
	privateKey *rsa.PrivateKey
	queue      queue.Service
}

type customClaims struct {
	*jwt.StandardClaims
	entity.CustomInfo
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewTokenService(tokenRepo Repository, pk *rsa.PrivateKey, q queue.Service) Service {
	return &tokenService{
		tokenRepo:  tokenRepo,
		privateKey: pk,
		queue:      q,
	}
}

func (s tokenService) Issue(userID, developerID string, scopes []string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	tid, _ := rand.GenerateRandomString(16)
	// set our claims
	t.Claims = &customClaims{
		&jwt.StandardClaims{
			Id:       tid,
			IssuedAt: time.Now().Unix(),
			Issuer:   "wall",
			Audience: developerID,
			Subject:  userID,
			// set the expire time default: 2 day
			ExpiresAt: time.Now().Add(time.Minute * 2880).Unix(),
		},
		entity.CustomInfo{
			Scopes: scopes,
		},
	}

	// Creat token string
	token, err := t.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	_, err = s.tokenRepo.InsertOne(tid, developerID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s tokenService) Revoke(tokenID, developerID string) error {

	// must check this token is generate by same developerID
	tokenInfo, err := s.tokenRepo.FindOne(tokenID)
	if err != nil {
		return err
	}

	if tokenInfo.OwnerID != developerID {
		err := errors.New("permission deny")
		return err
	}

	_, err = s.tokenRepo.UpdateTokenStatus(tokenID, developerID)
	if err != nil {
		return err
	}

	s.queue.Send(tokenID)
	return nil
}
