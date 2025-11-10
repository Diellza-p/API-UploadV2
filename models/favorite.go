package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Favorite struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userID" bson:"userID"`
	ContentID string             `json:"contentID" bson:"contentID"`
	AlbumID   string             `json:"albumID" bson:"albumID"`
	DateAdded time.Time          `json:"dateAdded" bson:"dateAdded"`
}

type Album struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID       string             `json:"userID" bson:"userID"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	DateCreated  time.Time          `json:"dateCreated" bson:"dateCreated"`
	DateModified time.Time          `json:"dateModified" bson:"dateModified"`
}
