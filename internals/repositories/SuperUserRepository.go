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
	UpdateSuperuser(ctx context.Context, superuser *types.Superuser) error
	FindSuperuserByUsername(ctx context.Context, username string) (*types.Superuser, error)
	FindSuperuserByResetToken(ctx context.Context, token string) (*types.Superuser, error)
	GetRoles(ctx context.Context) ([]string, error)
	FindActivityLogsByUserID(ctx context.Context, userID string) ([]types.UserActivityLog, error)
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

func (r *superuserRepo) UpdateSuperuser(ctx context.Context, superuser *types.Superuser) error {
	filter := bson.M{"_id": superuser.ID}
	update := bson.M{
		"$set": bson.M{
			"username":       superuser.Username,
			"password":       superuser.Password,
			"updated_at":     superuser.UpdatedAt,
			"is_2fa_enabled": superuser.Is2FAEnabled,
		},
	}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

func (r *superuserRepo) FindSuperuserByUsername(ctx context.Context, username string) (*types.Superuser, error) {
	var superuser types.Superuser
	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

func (r *superuserRepo) FindSuperuserByResetToken(ctx context.Context, token string) (*types.Superuser, error) {
	var superuser types.Superuser
	err := r.db.FindOne(ctx, bson.M{"reset_token": token}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

func (r *superuserRepo) GetRoles(ctx context.Context) ([]string, error) {
	// Assume roles are stored in a separate collection or in a predefined list
	// Example: Fetch roles from a roles collection
	var roles []string
	cursor, err := r.db.Database().Collection("roles").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var role string
		if err := cursor.Decode(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *superuserRepo) FindActivityLogsByUserID(ctx context.Context, userID string) ([]types.UserActivityLog, error) {
	var logs []types.UserActivityLog
	cursor, err := r.db.Database().Collection("activity_logs").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var log types.UserActivityLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}
