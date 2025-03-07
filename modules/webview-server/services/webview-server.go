package services

import (
	"context"
	"fmt"
	"notification-server/modules/webview-server/domain"
	dto "notification-server/modules/webview-server/dtos"
	"notification-server/modules/webview-server/models"
	"notification-server/modules/webview-server/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebViewService struct {
	repo *repositories.WebViewRepository
}

func NewWebviewService(repo *repositories.WebViewRepository) *WebViewService {
	return &WebViewService{repo: repo}
}

func (s *WebViewService) GetWebviewListService(ctx context.Context, keyword string, status string, limit int, nextPageToken string) (domain.WebViewResponse, error) {
	webviews, lastID, err := s.repo.GetWebviewList(ctx, keyword, status, limit, nextPageToken)
	if err != nil {
		return domain.WebViewResponse{}, err
	}

	nextPageToken = ""
	if lastID != "" {
		nextPageToken = lastID
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data: domain.GetWebViewList{
			List:          webviews,
			NextPageToken: nextPageToken,
		},
	}, nil
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

	objectID, updateErr := s.repo.ChangeWebviewStatus(ctx, req.ID, req.Status)
	if updateErr != nil {
		return domain.WebViewResponse{
			Message: "failed to update WebView status",
			Code:    500,
			Data:    nil,
		}, updateErr
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    objectID,
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

	deletedID, deleteErr := s.repo.DeleteWebview(ctx, req.ID)
	if deleteErr != nil {
		return domain.WebViewResponse{
			Message: "failed to delete WebView",
			Code:    500,
			Data:    nil,
		}, deleteErr
	}

	return domain.WebViewResponse{
		Message: "success",
		Code:    200,
		Data:    deletedID,
	}, nil
}
