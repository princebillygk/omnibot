package subscription

import (
	"context"
	"errors"
	"fmt"

	"github.com/princebillygk/omnibot/internal/config"
	"github.com/princebillygk/omnibot/internal/services/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Subscription struct {
	ID                         primitive.ObjectID `bson:"_id,omitempty"`
	SubscriberId               primitive.ObjectID `bson:"subscriber_id"`
	Subject                    string             `bson:"subject"`
	MessengerNotificationToken string             `bson:"messeger_notification_token,omitempty"`
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

type SubscriptionOptions struct {
	Subject                    string
	MessengerNotificationToken string
}

func (s *Service) Subscribe(ctx context.Context, user users.User, opts SubscriptionOptions) error {
	var err error
	var subs Subscription

	switch user.Platform {
	case users.MessengerPlatform:
		if opts.MessengerNotificationToken == "" {
			return errors.New("Messenger Notification is required")
		}

		subs.SubscriberId = user.ID
		subs.Subject = opts.Subject
		subs.MessengerNotificationToken = opts.MessengerNotificationToken
	default:
		return fmt.Errorf("Subscription service is not yet available for %s platform yet.", user.Platform)
	}

	_, err = s.collection.InsertOne(ctx, subs)
	if me := err.(mongo.WriteException); me.HasErrorCode(11000) {
		return config.ApplicationError{
			HttpStatus:   200,
			Message:      fmt.Sprintf("You are already subscribed to %s!", opts.Subject),
			DebugMessage: fmt.Sprintf("User %s is already subscribed to %s!", user.ID, opts.Subject),
		}
	}
	return err
}
