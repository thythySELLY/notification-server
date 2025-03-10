package services

import (
	"context"
	"encoding/json"
	"fmt"
	connectionRepositories "notification-server/modules/connection/repositories"
	"notification-server/modules/webview-server/domain"
	dto "notification-server/modules/webview-server/dtos"
	"notification-server/modules/webview-server/models"
	"notification-server/modules/webview-server/repositories"
	"notification-server/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebViewService struct {
	repo           *repositories.WebViewRepository
	connectionRepo *connectionRepositories.ConnectionRepository
}

func NewWebviewService(repo *repositories.WebViewRepository, connectionRepo *connectionRepositories.ConnectionRepository) *WebViewService {
	return &WebViewService{repo: repo, connectionRepo: connectionRepo}
}

func (s *WebViewService) GetWebviewListService(ctx context.Context, keyword string, status string, limit int, nextPageToken string) (domain.WebViewResponse, error) {
	cacheKey := fmt.Sprintf("webview_list:%s:%s:%d:%s", keyword, status, limit, nextPageToken)

	cachedData, err := helpers.GetCache(cacheKey)
	if err == nil {
		var cachedResponse domain.WebViewResponse
		if jsonErr := json.Unmarshal([]byte(cachedData), &cachedResponse); jsonErr == nil {
			return cachedResponse, nil
		}
	}

	webviews, lastID, err := s.repo.GetWebviewList(ctx, keyword, status, limit, nextPageToken)
	if err != nil {
		return domain.WebViewResponse{}, err
	}

	nextPageToken = ""
	if lastID != "" {
		nextPageToken = lastID
	}

	response := domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data: domain.GetWebViewList{
			List:          webviews,
			NextPageToken: nextPageToken,
		},
	}

	jsonData, _ := json.Marshal(response)
	_ = helpers.SetCache(cacheKey, string(jsonData))

	return response, nil
}

func (s *WebViewService) CreateWebviewService(ctx context.Context, req dto.CreateWebviewServer) (domain.WebViewResponse, error) {
	exists, err := s.repo.IsWebviewExistsByName(ctx, req.Name)
	if err != nil {
		return domain.WebViewResponse{
			Message: "failed to check existing WebView",
			Code:    500,
			Data:    nil,
		}, err
	}
	if exists {
		return domain.WebViewResponse{
			Message: "WebView name already exists",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("webview with name '%s' already exists", req.Name)
	}

	objectID := primitive.NewObjectID()
	now := time.Now()

	webview := models.WebViewServer{
		ID:        objectID.Hex(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      req.Name,
		Status:    string(models.StatusInactive),
	}

	err = s.repo.CreateWebview(ctx, &webview)
	if err != nil {
		return domain.WebViewResponse{
			Message: "failed to create WebView",
			Code:    500,
			Data:    nil,
		}, err
	}

	responseData := domain.CreateWebViewServer{ID: webview.ID}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    responseData,
	}, nil
}

func (s *WebViewService) UpdateWebviewService(ctx context.Context, req dto.UpdateWebviewServer) (domain.WebViewResponse, error) {
	exists, err := s.repo.IsWebviewExistsByID(ctx, req.ID)
	if err != nil {
		return domain.WebViewResponse{
			Message: "failed to check existing WebView",
			Code:    500,
			Data:    nil,
		}, err
	}
	if !exists {
		return domain.WebViewResponse{
			Message: "WebView not found",
			Code:    404,
			Data:    nil,
		}, fmt.Errorf("webview with id '%s' does not exist", req.ID)
	}

	if existsByName, err := s.repo.IsWebviewExistsByName(ctx, req.Name); err != nil {
		return domain.WebViewResponse{
			Message: "failed to check existing WebView by name",
			Code:    500,
			Data:    nil,
		}, err
	} else if existsByName {
		return domain.WebViewResponse{
			Message: "WebView name already exists",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("webview with name '%s' already exists", req.Name)
	}

	updateID, updateErr := s.repo.UpdateWebview(ctx, req.ID, req.Name)
	if updateErr != nil {
		return domain.WebViewResponse{
			Message: "failed to update WebView",
			Code:    500,
			Data:    nil,
		}, updateErr
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    updateID,
	}, nil
}

func (s *WebViewService) ChangeWebviewStatus(ctx context.Context, req dto.ChangeWebviewServerStatus) (domain.WebViewResponse, error) {
	webview, err := s.repo.GetWebviewByID(ctx, req.ID)
	if err != nil {
		return domain.WebViewResponse{
			Message: "WebView not found, no status change applied",
			Code:    500,
			Data:    nil,
		}, err
	}

	if webview.Status == req.Status {
		return domain.WebViewResponse{
			Message: "WebView status is already the same as the requested status",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("webview with id '%s' already has the requested status '%s'", req.ID, req.Status)
	}

	session, err := s.repo.StartSession(ctx)
	if err != nil {
		return domain.WebViewResponse{
			Message: "Failed to start transaction",
			Code:    500,
			Data:    nil,
		}, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		objectID, updateErr := s.repo.ChangeWebviewStatus(sessCtx, req.ID, req.Status)
		if updateErr != nil {
			return nil, updateErr
		}

		connections, connErr := s.connectionRepo.GetConnectionByWebviewId(sessCtx, req.ID)
		if connErr != nil {
			return nil, connErr
		}

		for _, conn := range connections {
			if conn.Status == "active" && req.Status == "inactive" {
				_, err := s.connectionRepo.ChangeConnectionStatus(sessCtx, conn.ID, "inactive")
				if err != nil {
					return nil, err
				}
			}
		}

		return objectID, nil
	})

	if err != nil {
		return domain.WebViewResponse{
			Message: "Failed to update WebView and associated connections",
			Code:    500,
			Data:    nil,
		}, err
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    result,
	}, nil
}

func (s *WebViewService) DeleteWebviewService(ctx context.Context, req dto.DeleteWebviewServer) (domain.WebViewResponse, error) {

	exists, err := s.repo.IsWebviewExistsByID(ctx, req.ID)
	if err != nil {
		return domain.WebViewResponse{
			Message: "failed to check existing WebView",
			Code:    500,
			Data:    nil,
		}, err
	}
	if !exists {
		return domain.WebViewResponse{
			Message: "WebView not found",
			Code:    404,
			Data:    nil,
		}, fmt.Errorf("webview with id '%s' does not exist", req.ID)
	}

	session, err := s.repo.StartSession(ctx)
	if err != nil {
		return domain.WebViewResponse{
			Message: "Failed to start transaction",
			Code:    500,
			Data:    nil,
		}, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		connections, connErr := s.connectionRepo.GetConnectionByWebviewId(sessCtx, req.ID)
		if connErr != nil {
			return nil, connErr
		}

		for _, conn := range connections {
			if err := s.connectionRepo.DeleteConnection(sessCtx, conn.ID); err != nil {
				return nil, err
			}
		}

		deletedID, deleteErr := s.repo.DeleteWebview(sessCtx, req.ID)
		if deleteErr != nil {
			return nil, deleteErr
		}

		return deletedID, nil
	})

	if err != nil {
		return domain.WebViewResponse{
			Message: "Failed to delete WebView and associated connections",
			Code:    500,
			Data:    nil,
		}, err
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    result,
	}, nil
}
