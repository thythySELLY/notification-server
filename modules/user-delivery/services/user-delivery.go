package services

import (
	"context"
	"fmt"
	connectionRepositories "notification-server/modules/connection/repositories"
	"notification-server/modules/user-delivery/domain"
	dto "notification-server/modules/user-delivery/dtos"
	"notification-server/modules/user-delivery/models"
	"notification-server/modules/user-delivery/repositories"
	webviewRepositories "notification-server/modules/webview-server/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDeliveryService struct {
	repo           *repositories.UserDeliveryRepository
	connectionRepo *connectionRepositories.ConnectionRepository
	webviewRepo    *webviewRepositories.WebViewRepository
}

func NewUserDeliveryService(repo *repositories.UserDeliveryRepository, connectionRepo *connectionRepositories.ConnectionRepository, webviewRepo *webviewRepositories.WebViewRepository) *UserDeliveryService {
	return &UserDeliveryService{
		repo:           repo,
		connectionRepo: connectionRepo,
		webviewRepo:    webviewRepo,
	}
}

func (s *UserDeliveryService) GetUserDeliveryList(ctx context.Context, keyword string, status string, limit int, nextPageToken string) (domain.UserDeliveryResponse, error) {
	userDeliveries, lastID, err := s.repo.GetUserDeliveryList(ctx, keyword, status, limit, nextPageToken)
	if err != nil {
		return domain.UserDeliveryResponse{}, err
	}

	nextPageToken = ""
	if lastID != "" {
		nextPageToken = lastID
	}

	return domain.UserDeliveryResponse{
		Message: "success",
		Code:    200,
		Data: domain.GetUserDeliveryList{
			List:          userDeliveries,
			NextPageToken: nextPageToken,
		},
	}, nil
}

func (s *UserDeliveryService) CreateUserDelivery(ctx context.Context, req dto.CreateUserDelivery) (domain.UserDeliveryResponse, error) {
	exists, err := s.repo.IsUserDeliveryExistsByName(ctx, req.Name)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to check existing User Delivery",
			Code:    500,
			Data:    nil,
		}, err
	}
	if exists {
		return domain.UserDeliveryResponse{
			Message: "User Delivery name already exists",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("user delivery with name '%s' already exists", req.Name)
	}

	objectID := primitive.NewObjectID()
	now := time.Now()

	userDelivery := models.UserDelivery{
		ID:        objectID.Hex(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      req.Name,
		Status:    string(models.StatusInactive),
	}

	err = s.repo.CreateUserDelivery(ctx, &userDelivery)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to create User Delivery",
			Code:    500,
			Data:    nil,
		}, err
	}

	responseData := domain.CreateUserDelivery{ID: userDelivery.ID}
	return domain.UserDeliveryResponse{
		Message: "success",
		Code:    200,
		Data:    responseData,
	}, nil
}

func (s *UserDeliveryService) ChangeUserDeliveryStatus(ctx context.Context, req dto.ChangeUserDeliveryStatus) (domain.UserDeliveryResponse, error) {
	userDelivery, err := s.repo.GetUserDeliveryByID(ctx, req.ID)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "User Delivery not found, no status change applied",
			Code:    500,
			Data:    nil,
		}, err
	}

	if userDelivery.Status == req.Status {
		return domain.UserDeliveryResponse{
			Message: "User Delivery status is already the same as the requested status",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("user delivery with id '%s' already has the requested status '%s'", req.ID, req.Status)
	}

	session, err := s.repo.StartSession(ctx)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "Failed to start transaction",
			Code:    500,
			Data:    nil,
		}, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		objectID, updateErr := s.repo.ChangeUserDeliveryStatus(sessCtx, req.ID, req.Status)
		if updateErr != nil {
			return nil, updateErr
		}

		connections, connErr := s.connectionRepo.GetConnectionByUserDeliveryId(sessCtx, req.ID)
		if connErr != nil {
			return nil, connErr
		}

		for _, conn := range connections {
			if conn.Status == "active" {
				_, err := s.connectionRepo.ChangeConnectionStatus(sessCtx, conn.ID, "inactive")
				if err != nil {
					return nil, err
				}
			} else {
				isActive, err := s.webviewRepo.IsWebviewActive(sessCtx, conn.WebviewServerId)
				if err != nil {
					return nil, err
				}
				if isActive {
					_, err := s.connectionRepo.ChangeConnectionStatus(sessCtx, conn.ID, "active")
					if err != nil {
						return nil, err
					}
				}
			}
		}

		return objectID, nil
	})

	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "Failed to update User Delivery and associated connections",
			Code:    500,
			Data:    nil,
		}, err
	}

	return domain.UserDeliveryResponse{
		Message: "success",
		Code:    200,
		Data:    result,
	}, nil
}

func (s *UserDeliveryService) UpdateUserDeliveryService(ctx context.Context, req dto.UpdateUserDelivery) (domain.UserDeliveryResponse, error) {
	exists, err := s.repo.IsUserDeliveryExistsByID(ctx, req.ID)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to check existing User Delivery",
			Code:    500,
			Data:    nil,
		}, err
	}
	if !exists {
		return domain.UserDeliveryResponse{
			Message: "User Delivery not found",
			Code:    404,
			Data:    nil,
		}, fmt.Errorf("user delivery with id '%s' does not exist", req.ID)
	}

	if existsByName, err := s.repo.IsUserDeliveryExistsByName(ctx, req.Name); err != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to check existing User Delivery by name",
			Code:    500,
			Data:    nil,
		}, err
	} else if existsByName {
		return domain.UserDeliveryResponse{
			Message: "User Delivery name already exists",
			Code:    400,
			Data:    nil,
		}, fmt.Errorf("user delivery with name '%s' already exists", req.Name)
	}

	updateID, updateErr := s.repo.UpdateUserDelivery(ctx, req.ID, req.Name)
	if updateErr != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to update User Delivery",
			Code:    500,
			Data:    nil,
		}, updateErr
	}

	return domain.UserDeliveryResponse{
		Message: "success",
		Code:    200,
		Data:    updateID,
	}, nil
}

func (s *UserDeliveryService) DeleteUserDeliveryService(ctx context.Context, req dto.DeleteUserDelivery) (domain.UserDeliveryResponse, error) {
	exists, err := s.repo.IsUserDeliveryExistsByID(ctx, req.ID)
	if err != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to check existing User Delivery",
			Code:    500,
			Data:    nil,
		}, err
	}
	if !exists {
		return domain.UserDeliveryResponse{
			Message: "User Delivery not found",
			Code:    404,
			Data:    nil,
		}, fmt.Errorf("user delivery with id '%s' does not exist", req.ID)
	}

	deletedID, deleteErr := s.repo.DeleteUserDelivery(ctx, req.ID)
	if deleteErr != nil {
		return domain.UserDeliveryResponse{
			Message: "failed to delete User Delivery",
			Code:    500,
			Data:    nil,
		}, deleteErr
	}

	return domain.UserDeliveryResponse{
		Message: "success",
		Code:    200,
		Data:    deletedID,
	}, nil
}
