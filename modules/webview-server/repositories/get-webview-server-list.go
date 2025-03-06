package repositories

import (
	"context"
	"notification-server/modules/webview-server/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WebViewRepository định nghĩa các hàm làm việc với MongoDB
type WebViewRepository struct {
	collection *mongo.Collection
}

// NewWebViewRepository khởi tạo repository
func NewWebViewRepository(db *mongo.Database) *WebViewRepository {
	return &WebViewRepository{
		collection: db.Collection("webviews"),
	}
}

// GetWebViews lấy danh sách webviews từ MongoDB
func (r *WebViewRepository) GetWebViews(ctx context.Context, keyword string, status string, limit int) ([]models.WebViewServer, error) {
	var webviews []models.WebViewServer
	filter := bson.M{}

	if keyword != "" {
		filter["name"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	if status != "" {
		filter["status"] = status
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var webview models.WebViewServer
		if err := cursor.Decode(&webview); err != nil {
			return nil, err
		}
		webviews = append(webviews, webview)
	}

	return webviews, nil
}
