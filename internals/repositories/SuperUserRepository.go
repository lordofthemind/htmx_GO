package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/htmx_GO/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SuperuserRepository interface {
	CreateSuperuser(ctx context.Context, superuser *types.SuperUserType) error
	FindSuperuserByEmail(ctx context.Context, email string) (*types.SuperUserType, error)
	FindSuperuserByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error)
	UpdateSuperuser(ctx context.Context, superuser *types.SuperUserType) error
	FindSuperuserByUsername(ctx context.Context, username string) (*types.SuperUserType, error)
	FindSuperuserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error)
	DeleteSuperuserByID(ctx context.Context, id uuid.UUID) error
	ListSuperusers(ctx context.Context, limit, skip int64) ([]*types.SuperUserType, error)
	UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (string, error)
	Enable2FA(ctx context.Context, id uuid.UUID, isEnabled bool) error
	SearchSuperusers(ctx context.Context, searchQuery string) ([]*types.SuperUserType, error)
	SoftDeleteSuperuser(ctx context.Context, id uuid.UUID) error
	FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error)
	UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error
	BulkUpdateSuperusers(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}) error
}

type MongoSuperuserRepo struct {
	db *mongo.Collection
}

func NewMongoSuperuserRepository(db *mongo.Database) SuperuserRepository {
	return &MongoSuperuserRepo{
		db: db.Collection("superusers"),
	}
}

// CreateSuperuser creates a new superuser.
func (r *MongoSuperuserRepo) CreateSuperuser(ctx context.Context, superuser *types.SuperUserType) error {
	superuser.CreatedAt = time.Now().Unix()
	superuser.UpdatedAt = time.Now().Unix()
	if superuser.ID == uuid.Nil {
		superuser.ID = uuid.New()
	}
	_, err := r.db.InsertOne(ctx, superuser)
	return err
}

// FindSuperuserByEmail finds a superuser by email.
func (r *MongoSuperuserRepo) FindSuperuserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superuser types.SuperUserType
	err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

// FindSuperuserByID finds a superuser by ID.
func (r *MongoSuperuserRepo) FindSuperuserByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	var superuser types.SuperUserType
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

// UpdateSuperuser updates a superuser's details.
func (r *MongoSuperuserRepo) UpdateSuperuser(ctx context.Context, superuser *types.SuperUserType) error {
	filter := bson.M{"_id": superuser.ID}
	update := bson.M{
		"$set": bson.M{
			"username":       superuser.Username,
			"password":       superuser.Password,
			"updated_at":     time.Now().Unix(),
			"is_2fa_enabled": superuser.Is2FAEnabled,
		},
	}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

// FindSuperuserByUsername finds a superuser by username.
func (r *MongoSuperuserRepo) FindSuperuserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superuser types.SuperUserType
	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

// FindSuperuserByResetToken finds a superuser by reset token.
func (r *MongoSuperuserRepo) FindSuperuserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superuser types.SuperUserType
	err := r.db.FindOne(ctx, bson.M{"reset_token": token}).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superuser, err
}

func (r *MongoSuperuserRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	var superuser types.SuperUserType
	filter := bson.M{"_id": id}

	// Only fetching the role field, for performance optimization
	projection := bson.M{"role": 1}

	err := r.db.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&superuser)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("superuser not found")
	}

	return superuser.Role, err
}

// DeleteSuperuserByID deletes a superuser by their ID.
func (r *MongoSuperuserRepo) DeleteSuperuserByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListSuperusers lists superusers with pagination.
func (r *MongoSuperuserRepo) ListSuperusers(ctx context.Context, limit, skip int64) ([]*types.SuperUserType, error) {
	var superusers []*types.SuperUserType
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	cursor, err := r.db.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &superusers); err != nil {
		return nil, err
	}
	return superusers, nil
}

// UpdateResetToken updates the reset token for a superuser.
func (r *MongoSuperuserRepo) UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"reset_token": token, "updated_at": time.Now().Unix()}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

// Enable2FA enables or disables 2FA for a superuser.
func (r *MongoSuperuserRepo) Enable2FA(ctx context.Context, id uuid.UUID, isEnabled bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_2fa_enabled": isEnabled, "updated_at": time.Now().Unix()}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

// SoftDeleteSuperuser marks a superuser as archived instead of permanently deleting.
func (r *MongoSuperuserRepo) SoftDeleteSuperuser(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"archived": true, "updated_at": time.Now().Unix()}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

// SearchSuperusers allows partial search by full_name, username, or email.
func (r *MongoSuperuserRepo) SearchSuperusers(ctx context.Context, searchQuery string) ([]*types.SuperUserType, error) {
	var superusers []*types.SuperUserType
	filter := bson.M{
		"$or": []bson.M{
			{"full_name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"username": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &superusers); err != nil {
		return nil, err
	}
	return superusers, nil
}

// FindAll2FAEnabledSuperusers finds all superusers with 2FA enabled.
func (r *MongoSuperuserRepo) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	var superusers []*types.SuperUserType
	filter := bson.M{"is_2fa_enabled": true}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &superusers); err != nil {
		return nil, err
	}
	return superusers, nil
}

// UpdateSuperuserRole updates the role of a superuser.
func (r *MongoSuperuserRepo) UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"role": role, "updated_at": time.Now().Unix()}}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

// BulkUpdateSuperusers updates multiple superusers at once.
func (r *MongoSuperuserRepo) BulkUpdateSuperusers(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": updates}
	_, err := r.db.UpdateMany(ctx, filter, update)
	return err
}
