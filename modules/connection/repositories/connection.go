package repositories

import (
	"context"
	"fmt"
	"notification-server/helpers"
	"notification-server/modules/connection/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnectionRepository struct {
	collection *mongo.Collection
}

func NewConnectionRepository(db *mongo.Database) *ConnectionRepository {
	return &ConnectionRepository{
		collection: db.Collection("connections"),
	}
}

func (repo *ConnectionRepository) IsHavingSameConnection(userDeliveryId string, webviewServerId string) (bool, error) {
	filter := bson.M{
		"userDeliveryServerId": userDeliveryId,
		"webviewServerId":      webviewServerId,
	}

	count, err := repo.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ConnectionRepository) CreateConnection(connect models.Connection) error {
	objectID, err := primitive.ObjectIDFromHex(connect.ID)
	if err != nil {
		return err
	}

	newConnection := bson.M{
		"_id":                          objectID,
		"createdAt":                    connect.CreatedAt,
		"updatedAt":                    connect.UpdatedAt,
		"status":                       connect.Status,
		"webviewServerApiKey":          connect.WebviewServerApiKey,
		"userDeliveryServerApiKey":     connect.UserDeliveryServerApiKey,
		"webviewServerId":              connect.WebviewServerId,
		"userDeliveryServerId":         connect.UserDeliveryServerId,
		"userDeliveryServerWebHookUrl": connect.UserDeliveryServerWebHookUrl,
	}

	_, err = repo.collection.InsertOne(context.TODO(), newConnection)
	return err
}

func (repo *ConnectionRepository) GetConnections(ctx context.Context, userDeliveryId string, webviewID string, status string, limit int, nextPageToken string) ([]models.Connection, string, error) {
	var connections []models.Connection
	filter := bson.M{}

	if userDeliveryId != "" {
		objectID, err := helpers.StringToObjectID(userDeliveryId)
		if err != nil {
			return nil, "", fmt.Errorf("invalid userDeliveryId: %s", userDeliveryId)
		}
		filter["userDeliveryServerId"] = objectID
	}

	if webviewID != "" {
		objectID, err := helpers.StringToObjectID(webviewID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid webviewID: %s", webviewID)
		}
		filter["webviewServerId"] = objectID
	}

	if status != "" {
		filter["status"] = status
	}

	if nextPageToken != "" {
		tokenID, err := helpers.StringToObjectID(nextPageToken)
		if err != nil {
			return nil, "", fmt.Errorf("invalid nextPageToken: %s", nextPageToken)
		}
		filter["_id"] = bson.M{"$gt": tokenID}
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	limitInt64 := int64(limit)
	cursor, err := repo.collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limitInt64,
		Sort:  bson.M{"_id": 1},
	})
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var lastID string
	for cursor.Next(ctx) {
		var connection models.Connection
		if err := cursor.Decode(&connection); err != nil {
			return nil, "", err
		}

		connections = append(connections, connection)
		lastID = connection.ID
	}

	return connections, lastID, nil
}

func (repo *ConnectionRepository) IsHavingConnectionById(ctx context.Context, connectionId string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(connectionId)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id": objectID,
	}

	count, err := repo.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ConnectionRepository) UpdateUserDeliveryHookUrl(ctx context.Context, id string, newUserDeliveryHookUrl string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"userDeliveryServerWebHookUrl": newUserDeliveryHookUrl}}

	_, err = repo.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *ConnectionRepository) ChangeConnectionStatus(ctx context.Context, id string, status string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ConnectionRepository) GetConnectionByID(ctx context.Context, id string) (models.Connection, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Connection{}, err
	}

	var connection models.Connection
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&connection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Connection{}, nil
		}
		return models.Connection{}, err
	}

	return connection, nil
}

func (repo *ConnectionRepository) GetConnectionByUserDeliveryId(ctx context.Context, userDeliveryId string) ([]models.Connection, error) {
	filter := bson.M{"userDeliveryServerId": userDeliveryId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var connections []models.Connection
	for cursor.Next(ctx) {
		var connection models.Connection
		if err := cursor.Decode(&connection); err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

func (repo *ConnectionRepository) GetConnectionByWebviewId(ctx context.Context, webviewId string) ([]models.Connection, error) {
	filter := bson.M{"webviewServerId": webviewId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var connections []models.Connection
	for cursor.Next(ctx) {
		var connection models.Connection
		if err := cursor.Decode(&connection); err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

func (repo *ConnectionRepository) DeleteConnection(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = repo.collection.DeleteOne(ctx, filter)
	return err
}
