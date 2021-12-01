package token

import (
	"wall/pkg/entity"
)

//Repository interface
type Repository interface {
	InsertOne(tid, ownerID string) (bool, error)
	FindOne(tid string) (*entity.TokenInfo, error)
	UpdateTokenStatus(tid, ownerID string) (bool, error)
}

//Service  interface
type Service interface {
	Issue(userID, developerID string, scopes []string) (string, error)
	Revoke(tokenID, developerID string) error
}
