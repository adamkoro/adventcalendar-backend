package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdventCalendarDay struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Day     int32              `bson:"day,omitempty"`
	Year    int32              `bson:"year,omitempty"`
	Title   string             `bson:"title,omitempty"`
	Content string             `bson:"content,omitempty"`
}

type Repository struct {
	Db  *mongo.Client
	Ctx *context.Context
}
