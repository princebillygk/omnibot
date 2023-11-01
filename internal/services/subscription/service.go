package subscription

import (
	"context"
	"fmt"

	"github.com/princebillygk/omnibot/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Subscription struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  string             `bson:"userID"`
	Subject string             `bson:"subject"`
}

type Service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	ctx := context.TODO()

	coll := db.Collection("subscriptions")

	coll.Indexes().DropAll(ctx)

	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userID", Value: 1},
			{Key: "subject", Value: 1},
		},
			Options: options.Index().SetUnique(true),
		}})
	if err != nil {
		panic(err)
	}

	return &Service{db, coll}
}

func (s *Service) Subscribe(ctx context.Context, userID string, subject string) error {
	_, err := s.collection.InsertOne(ctx, Subscription{
		UserID:  userID,
		Subject: subject,
	})

	if me := err.(mongo.WriteException); me.HasErrorCode(11000) {
		return config.ApplicationError{
			HttpStatus:   200,
			Message:      fmt.Sprintf("You are already subscribed to %s!", subject),
			DebugMessage: fmt.Sprintf("User %s is already subscribed to %s!", userID, subject),
		}
	}
	return err
}
