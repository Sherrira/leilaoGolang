package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sherrira/leilaoGolang/configuration/logger"
	"github.com/Sherrira/leilaoGolang/internal/entity/user_entity"
	"github.com/Sherrira/leilaoGolang/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(ctx context.Context, id string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var userEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			msg := fmt.Sprintf("user with id %s not found", id)
			logger.Error(msg, err)
			return nil, internal_error.NewNotFoundError(msg)
		}

		logger.Error("error finding user by id", err)
		return nil, internal_error.NewInternalServerError("error finding user by id")
	}

	userEntity := &user_entity.User{
		Id:   userEntityMongo.Id,
		Name: userEntityMongo.Name,
	}

	return userEntity, nil
}
