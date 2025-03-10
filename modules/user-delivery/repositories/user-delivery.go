package repositories

import (
	"context"
	"notification-server/helpers"
	"notification-server/modules/user-delivery/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDeliveryRepository struct {
	list   *mongo.Collection
	client *mongo.Client
}

func NewUserDeliveryRepository(db *mongo.Database, client *mongo.Client) *UserDeliveryRepository {
	return &UserDeliveryRepository{
		list:   db.Collection("user-deliveries"),
		client: client,
	}
}

func (r *UserDeliveryRepository) StartSession(ctx context.Context) (mongo.Session, error) {
	return r.client.StartSession()
}

func (r *UserDeliveryRepository) GetUserDeliveryList(ctx context.Context, keyword string, status string, limit int, nextPageToken string) ([]models.UserDelivery, string, error) {
	var userDeliveries []models.UserDelivery
	filter := bson.M{}

	if keyword != "" {
		filter["name"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	if status != "" {
		filter["status"] = status
	}
	if nextPageToken != "" {
		tokenID, err := helpers.StringToObjectID(nextPageToken)
		if err == nil {
			filter["_id"] = bson.M{"$gt": tokenID}
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	limitInt64 := int64(limit)
	cursor, err := r.list.Find(ctx, filter, &options.FindOptions{
		Limit: &limitInt64,
		Sort:  bson.M{"_id": 1},
	})
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var lastID string
	for cursor.Next(ctx) {
		var UserDelivery models.UserDelivery
		var temp struct {
			ID        primitive.ObjectID `bson:"_id"`
			CreatedAt time.Time          `bson:"createdAt"`
			UpdatedAt time.Time          `bson:"updatedAt"`
			Name      string             `bson:"name"`
			Status    string             `bson:"status"`
		}

		if err := cursor.Decode(&temp); err != nil {
			return nil, "", err
		}

		UserDelivery.ID = helpers.ObjectIDToString(temp.ID)
		UserDelivery.CreatedAt = temp.CreatedAt
		UserDelivery.UpdatedAt = temp.UpdatedAt
		UserDelivery.Name = temp.Name
		UserDelivery.Status = temp.Status

		userDeliveries = append(userDeliveries, UserDelivery)
		lastID = UserDelivery.ID
	}

	return userDeliveries, lastID, nil
}

func (r *UserDeliveryRepository) CreateUserDelivery(ctx context.Context, userDelivery *models.UserDelivery) error {

	objectID, err := primitive.ObjectIDFromHex(userDelivery.ID)
	if err != nil {
		return err
	}

	userDeliveryDocument := bson.M{
		"_id":       objectID,
		"createdAt": userDelivery.CreatedAt,
		"updatedAt": userDelivery.UpdatedAt,
		"name":      userDelivery.Name,
		"status":    userDelivery.Status,
	}

	_, err = r.list.InsertOne(ctx, userDeliveryDocument)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserDeliveryRepository) IsUserDeliveryExistsByName(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}

	count, err := r.list.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserDeliveryRepository) GetUserDeliveryByID(ctx context.Context, id string) (*models.UserDelivery, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var userDelivery models.UserDelivery
	err = r.list.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userDelivery)
	if err != nil {
		return nil, err
	}

	return &userDelivery, nil
}

func (r *UserDeliveryRepository) UpdateUserDelivery(ctx context.Context, id string, name string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"name": name,
		},
	}

	_, err = r.list.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserDeliveryRepository) IsUserDeliveryExistsByID(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	count, err := r.list.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserDeliveryRepository) ChangeUserDeliveryStatus(ctx context.Context, id string, status string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err = r.list.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserDeliveryRepository) DeleteUserDelivery(ctx context.Context, id string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	_, err = r.list.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserDeliveryRepository) IsUserDeliveryActive(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	var userDelivery models.UserDelivery
	err = r.list.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userDelivery)
	if err != nil {
		return false, err
	}

	return userDelivery.Status == string(models.StatusActive), nil // Assuming StatusActive is defined in models
}
