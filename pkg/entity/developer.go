package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

//Developer data struct
type Developer struct {
	ID       primitive.ObjectID `json:"-" bson:"_id" yaml:"-"`
	Name     string             `json:"name" bson:"name" yaml:"name"`
	Password string             `json:"password"  bson:"password" yaml:"password"`
	Grants   []string           `json:"grants" bson:"grants" yaml:"grants"`
}
