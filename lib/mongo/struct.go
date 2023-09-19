package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdventCalendarDay struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Day     uint8              `bson:"day" json:"day" binding:"required" validate:"required,min=1,max=31"`
	Year    uint16             `bson:"year" json:"year" binding:"required" validate:"required,min=1,max=9999"`
	Title   string             `bson:"title" json:"title" binding:"required" validate:"required,min=1,max=255"`
	Content string             `bson:"content" json:"content" binding:"required" validate:"required,min=1,max=65555"`
}

type AdventCalendarDayUpdate struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required" validate:"required"`
	Day     uint8              `bson:"day" json:"day" binding:"required" validate:"required,min=1,max=31"`
	Year    uint16             `bson:"year" json:"year" binding:"required" validate:"required,min=1,max=9999"`
	Title   string             `bson:"title" json:"title" binding:"required" validate:"required,min=1,max=255"`
	Content string             `bson:"content" json:"content" binding:"required" validate:"required,min=1,max=65555"`
}

type DayIDRequest struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required" validate:"required"`
}

type Repository struct {
	Db  *mongo.Client
	Ctx *context.Context
}
