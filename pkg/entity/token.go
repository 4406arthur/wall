package entity

type TokenInfo struct {
	Tid      string `json:"tid" bson:"tid"`
	OwnerID  string `json:"ownerID" bson:"ownerID"`
	Enabled  bool   `json:"enabled" bson:"enabled"`
	CreateAt int64  `json:"createAt" bson:"createAt"`
}

type CustomInfo struct {
	Scopes []string `json:"scopes"`
}
