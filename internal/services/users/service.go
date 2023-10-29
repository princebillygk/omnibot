package users

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Platform string

const (
	MessengerPlatform Platform = "MESSENGER"
)

type User struct {
	Platform       Platform  `bson:"platform"`
	userID         string    `bson:"userId,omitempty"`
	Name           string    `bson:"name,omitempty"`
	Email          string    `bson:"email,omitempty"`
	Phone          string    `bson:"phone,omitempty"`
	CreatedAt      time.Time `bson:"createdAt"`
	LastActivityAt time.Time `bson:"lastActivityAt"`
}

type Service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	coll := db.Collection("users")

	indexes, err := coll.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userID", Value: "1"}, {}},
			// Options: options.SetInde,
		},
	})
	if err != nil {
	}

	return &Service{db, coll}
}

func (s *Service) SaveUser(c context.Context, u *User) error {
	res, err := s.collection.InsertOne(c, u)
	fmt.Println(res)
	return err
}

func (s *Service) UpdateLastActivityTime(c context.Context, userID string, time time.Time) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "lastActivityAt", Value: time}}}}

	res, err := s.collection.UpdateByID(c, userID, update)
	fmt.Println(res)
	return err
}

func (s *Service) GetUser(c context.Context, userID string) error {
	filter := bson.D{{Key: "userID", Value: userID}}
	var user User
	err := s.collection.FindOne(c, filter).Decode(user)
	fmt.Println(user)
	return err
}

func (s *Service) DeleteUser(c context.Context, userID string) error {
	res, err := s.collection.DeleteOne(c, bson.D{{
		Key:   "_id",
		Value: userID,
	}})
	fmt.Println(res)
	return err
}
