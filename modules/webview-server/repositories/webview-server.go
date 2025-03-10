package repositories

import (
	"context"
	"notification-server/helpers"
	"notification-server/modules/webview-server/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebViewRepository struct {
	list   *mongo.Collection
	client *mongo.Client
}

func NewWebviewRepository(db *mongo.Database, client *mongo.Client) *WebViewRepository {
	return &WebViewRepository{
		list:   db.Collection("webviews"),
		client: client,
	}
}

func (r *WebViewRepository) StartSession(ctx context.Context) (mongo.Session, error) {
	return r.client.StartSession()
}

func (r *WebViewRepository) GetWebviewList(ctx context.Context, keyword string, status string, limit int, nextPageToken string) ([]models.WebViewServer, string, error) {
	var webviews []models.WebViewServer
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
		var webview models.WebViewServer
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

		webview.ID = helpers.ObjectIDToString(temp.ID)
		webview.CreatedAt = temp.CreatedAt
		webview.UpdatedAt = temp.UpdatedAt
		webview.Name = temp.Name
		webview.Status = temp.Status

		webviews = append(webviews, webview)
		lastID = webview.ID
	}

	return webviews, lastID, nil
}

func (r *WebViewRepository) CreateWebview(ctx context.Context, webview *models.WebViewServer) error {

	objectID, err := primitive.ObjectIDFromHex(webview.ID)
	if err != nil {
		return err
	}

	webviewDocument := bson.M{
		"_id":       objectID,
		"createdAt": webview.CreatedAt,
		"updatedAt": webview.UpdatedAt,
		"name":      webview.Name,
		"status":    webview.Status,
	}

	_, err = r.list.InsertOne(ctx, webviewDocument)
	if err != nil {
		return err
	}

	return nil
}

func (r *WebViewRepository) IsWebviewExistsByName(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}

	count, err := r.list.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *WebViewRepository) GetWebviewByID(ctx context.Context, id string) (*models.WebViewServer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var webview models.WebViewServer
	err = r.list.FindOne(ctx, bson.M{"_id": objectID}).Decode(&webview)
	if err != nil {
		return nil, err
	}

	return &webview, nil
}

func (r *WebViewRepository) UpdateWebview(ctx context.Context, id string, name string) (string, error) {
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

func (r *WebViewRepository) IsWebviewExistsByID(ctx context.Context, id string) (bool, error) {
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

func (r *WebViewRepository) ChangeWebviewStatus(ctx context.Context, id string, status string) (string, error) {
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

func (r *WebViewRepository) DeleteWebview(ctx context.Context, id string) (string, error) {
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

func (r *WebViewRepository) IsWebviewActive(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	var webview models.WebViewServer
	err = r.list.FindOne(ctx, bson.M{"_id": objectID}).Decode(&webview)
	if err != nil {
		return false, err
	}

	return webview.Status == string(models.StatusActive), nil
}
