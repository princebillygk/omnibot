package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Platform string

type ApplicationError struct {
	ErrCode    int
	httpStatus int
	message    string
}

func (ae ApplicationError) Error() string {
	return ae.message
}

var UserNotFoundErr = ApplicationError{
	ErrCode:    4041,
	httpStatus: http.StatusNotFound,
	message:    "User not found",
}

const (
	MessengerPlatform Platform = "MESSENGER"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Platform       Platform           `bson:"platform"`
	UserID         string             `bson:"userID"`
	Name           string             `bson:"name,omitempty"`
	Email          string             `bson:"email,omitempty"`
	Phone          string             `bson:"phone,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt"`
	LastActivityAt time.Time          `bson:"lastActivityAt"`
}

type Service struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	ctx := context.TODO()

	coll := db.Collection("users")

	_, err := coll.Indexes().DropAll(ctx)
	if err != nil {
		panic(err)
	}

	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userID", Value: 1},
			{Key: "platform", Value: 1},
		},
			Options: options.Index().SetUnique(true),
		}})
	if err != nil {
		panic(err)
	}

	return &Service{db, coll}
}

func (s *Service) SaveUser(ctx context.Context, u *User) error {
	res, err := s.collection.InsertOne(ctx, u)
	fmt.Println(res)
	return err
}

func (s *Service) UpdateLastActivityTime(c context.Context, id string, time time.Time) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "lastActivityAt", Value: time}}}}

	res, err := s.collection.UpdateByID(c, id, update)
	fmt.Println(res)
	return err
}

type GetUserInput struct {
	Platform Platform
	UserID   string
}

func (s *Service) GetUser(ctx context.Context, input *GetUserInput) (*User, error) {
	filter := bson.D{{
		Key:   "platform",
		Value: input.Platform,
	}, {
		Key:   "userID",
		Value: input.UserID,
	}}

	var u User
	err := s.collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (*User, error) {
	var u User

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, UserNotFoundErr
	}

	err = s.collection.FindOne(ctx, bson.D{{
		Key:   "_id",
		Value: oid,
	}}).Decode(&u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserNotFoundErr
		} else {
		}
	}
	return &u, nil
}

func (s *Service) DeleteUserByID(c context.Context, id string) error {
	res, err := s.collection.DeleteOne(c, bson.D{{
		Key:   "_id",
		Value: id,
	}})
	fmt.Println(res)
	return err
}
