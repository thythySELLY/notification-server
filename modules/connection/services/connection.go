package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"notification-server/helpers"
	"notification-server/modules/connection/domain"
	dto "notification-server/modules/connection/dtos"
	"notification-server/modules/connection/models"
	connectionRepositories "notification-server/modules/connection/repositories"
	userDeliveryRepositories "notification-server/modules/user-delivery/repositories"
	webviewRepositories "notification-server/modules/webview-server/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionService struct {
	connectionRepo   *connectionRepositories.ConnectionRepository
	userDeliveryRepo *userDeliveryRepositories.UserDeliveryRepository
	webviewRepo      *webviewRepositories.WebViewRepository
}

func NewConnectionService(connectionRepo *connectionRepositories.ConnectionRepository, userDeliveryRepo *userDeliveryRepositories.UserDeliveryRepository, webviewRepo *webviewRepositories.WebViewRepository) *ConnectionService {
	return &ConnectionService{
		connectionRepo:   connectionRepo,
		userDeliveryRepo: userDeliveryRepo,
		webviewRepo:      webviewRepo,
	}
}

func (service *ConnectionService) CreateConnection(ctx context.Context, req dto.CreateConnection) (primitive.ObjectID, error) {
	if req.WebviewServerId != "" {
		webviewExists, err := service.webviewRepo.IsWebviewExistsByID(ctx, req.WebviewServerId)
		if err != nil {
			return primitive.NilObjectID, err
		}
		if !webviewExists {
			return primitive.NilObjectID, errors.New("webview server does not exist")
		}
	}

	if req.UserDeliveryServerId != "" {
		userDeliveryExists, err := service.userDeliveryRepo.IsUserDeliveryExistsByID(ctx, req.UserDeliveryServerId)
		if err != nil {
			return primitive.NilObjectID, err
		}
		if !userDeliveryExists {
			return primitive.NilObjectID, errors.New("user delivery server does not exist")
		}
	}

	exists, err := service.connectionRepo.IsHavingSameConnection(req.UserDeliveryServerId, req.WebviewServerId)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if exists {
		return primitive.NilObjectID, errors.New("connection already exists")
	}

	webviewServerApiKey, err := generateRandomAPIKey()
	if err != nil {
		return primitive.NilObjectID, err
	}
	userDeliveryServerApiKey, err := generateRandomAPIKey()
	if err != nil {
		return primitive.NilObjectID, err
	}

	objectID := primitive.NewObjectID()

	newConnection := models.Connection{
		ID:                           objectID.Hex(),
		CreatedAt:                    time.Now(),
		UpdatedAt:                    time.Now(),
		Status:                       "inactive",
		WebviewServerApiKey:          webviewServerApiKey,
		UserDeliveryServerApiKey:     userDeliveryServerApiKey,
		WebviewServerId:              req.WebviewServerId,
		UserDeliveryServerId:         req.UserDeliveryServerId,
		UserDeliveryServerWebHookUrl: req.UserDeliveryServerWebHookUrl,
	}
	err = service.connectionRepo.CreateConnection(newConnection)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objectID, nil
}

func generateRandomAPIKey() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (service *ConnectionService) GetConnections(ctx context.Context, req dto.GetConnections) (domain.ConnectionResponse, error) {
	cacheKey := fmt.Sprintf("connections:%s:%s:%s:%d:%s", req.UserDeliveryServerId, req.WebviewServerId, req.Status, req.Limit, req.PageToken)

	cachedData, err := helpers.GetCache(cacheKey)
	if err == nil {
		var cachedResponse domain.ConnectionResponse
		if jsonErr := json.Unmarshal([]byte(cachedData), &cachedResponse); jsonErr == nil {
			return cachedResponse, nil
		}
	}

	if req.WebviewServerId != "" {
		webviewExists, err := service.webviewRepo.IsWebviewExistsByID(ctx, req.WebviewServerId)
		if err != nil {
			return domain.ConnectionResponse{}, err
		}
		if !webviewExists {
			return domain.ConnectionResponse{}, errors.New("webview server does not exist")
		}
	}

	if req.UserDeliveryServerId != "" {
		userDeliveryExists, err := service.userDeliveryRepo.IsUserDeliveryExistsByID(ctx, req.UserDeliveryServerId)
		if err != nil {
			return domain.ConnectionResponse{}, err
		}
		if !userDeliveryExists {
			return domain.ConnectionResponse{}, errors.New("user delivery server does not exist")
		}
	}

	connections, nextPageToken, err := service.connectionRepo.GetConnections(ctx, req.UserDeliveryServerId, req.WebviewServerId, req.Status, req.Limit, req.PageToken)
	if err != nil {
		return domain.ConnectionResponse{}, err
	}

	response := domain.ConnectionResponse{
		Message: "success",
		Code:    200,
		Data: domain.GetUserDeliveryList{
			List:          connections,
			NextPageToken: nextPageToken,
		},
	}

	jsonData, _ := json.Marshal(response)
	_ = helpers.SetCache(cacheKey, string(jsonData))

	return response, nil
}

func (service *ConnectionService) UpdateWebHookUrl(ctx context.Context, dto dto.UpdateUserDelivery) error {
	exists, err := service.connectionRepo.IsHavingConnectionById(ctx, dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("connection with ID %s does not exist", dto.ID)
	}

	return service.connectionRepo.UpdateUserDeliveryHookUrl(ctx, dto.ID, dto.UserDeliveryServerWebHookUrl)
}

func (s *ConnectionService) ChangeConnectionStatus(ctx context.Context, req dto.ChangeConnectionStatus) (domain.ConnectionResponse, error) {
	connection, err := s.connectionRepo.GetConnectionByID(ctx, req.ID)
	if err != nil {
		return domain.ConnectionResponse{
			Message: "Connection not found, no status change applied",
			Code:    500,
			Data:    nil,
		}, err
	}

	if connection.Status == req.Status {
		return domain.ConnectionResponse{
			Message: "Connection status is already the same as the requested status",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("connection with id '%s' already has the requested status '%s'", req.ID, req.Status)
	}

	if req.Status == "active" {
		isWebviewActive, err := s.webviewRepo.IsWebviewActive(ctx, connection.WebviewServerId)
		if err != nil {
			return domain.ConnectionResponse{
				Message: "Error checking webview server status",
				Code:    500,
				Data:    nil,
			}, err
		}
		if !isWebviewActive {
			return domain.ConnectionResponse{
				Message: "Webview server is not active, cannot change status to active",
				Code:    400,
				Data:    nil,
			}, fmt.Errorf("webview server with id '%s' is not active", connection.WebviewServerId)
		}

		isUserDeliveryActive, err := s.userDeliveryRepo.IsUserDeliveryActive(ctx, connection.UserDeliveryServerId)
		if err != nil {
			return domain.ConnectionResponse{
				Message: "Error checking user delivery server status",
				Code:    500,
				Data:    nil,
			}, err
		}
		if !isUserDeliveryActive {
			return domain.ConnectionResponse{
				Message: "User delivery server is not active, cannot change status to active",
				Code:    400,
				Data:    nil,
			}, fmt.Errorf("user delivery server with id '%s' is not active", connection.UserDeliveryServerId)
		}
	}

	objectID, updateErr := s.connectionRepo.ChangeConnectionStatus(ctx, req.ID, req.Status)
	if updateErr != nil {
		return domain.ConnectionResponse{
			Message: "failed to update Connection status",
			Code:    500,
			Data:    nil,
		}, updateErr
	}

	return domain.ConnectionResponse{
		Message: "success",
		Code:    200,
		Data:    objectID,
	}, nil
}

func (service *ConnectionService) DeleteConnection(ctx context.Context, dto dto.DeleteConnection) error {
	exists, err := service.connectionRepo.IsHavingConnectionById(ctx, dto.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("connection with ID %s does not exist", dto.ID)
	}

	return service.connectionRepo.DeleteConnection(ctx, dto.ID)
}
