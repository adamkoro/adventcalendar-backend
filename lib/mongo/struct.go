package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdventCalendarDay struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id" `
	Day     uint8              `bson:"day,omitempty" json:"day,omitempty" binding:"required" validate:"required,min=1,max=31"`
	Year    uint16             `bson:"year,omitempty" json:"year,omitempty" binding:"required" validate:"required,min=1,max=9999"`
	Title   string             `bson:"title,omitempty" json:"title,omitempty" binding:"required" validate:"required,min=1,max=255"`
	Content string             `bson:"content,omitempty" json:"content,omitempty" binding:"required" validate:"required,min=1,max=65555"`
}

type Repository struct {
	Db  *mongo.Client
	Ctx *context.Context
}
