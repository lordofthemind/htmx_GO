package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/lordofthemind/htmx_GO/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SuperuserRepository interface {
	CreateSuperuser(ctx context.Context, superuser *types.Superuser) error
	FindSuperuserByEmail(ctx context.Context, email string) (*types.Superuser, error)
	FindSuperuserByID(ctx context.Context, id primitive.ObjectID) (*types.Superuser, error)
}

type superuserRepo struct {
	db *mongo.Collection
}

func NewSuperuserRepository(db *mongo.Database) SuperuserRepository {
	return &superuserRepo{
		db: db.Collection("superusers"),
	}
}

func (r *superuserRepo) CreateSuperuser(ctx context.Context, superuser *types.Superuser) error {
	superuser.CreatedAt = time.Now().Unix()
	superuser.UpdatedAt = time.Now().Unix()
	_, err := r.db.InsertOne(ctx, superuser)
	return err
}

func (r *superuserRepo) FindSuperuserByEmail(ctx context.Context, email string) (*types.Superuser, error) {
	var superuser types.Superuser
	err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

func (r *superuserRepo) FindSuperuserByID(ctx context.Context, id primitive.ObjectID) (*types.Superuser, error) {
	var superuser types.Superuser
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}
