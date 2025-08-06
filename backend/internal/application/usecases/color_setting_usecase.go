package usecases

import (
	"context"
	"days/internal/application/dto"
	"days/internal/domain/entities"
	"days/internal/domain/repositories"
	"days/internal/domain/services"
	"days/internal/domain/valueobjects"
	"errors"
)

// ColorSettingUseCase handles color setting-related business operations
type ColorSettingUseCase struct {
	colorSettingRepo repositories.ColorSettingRepository
	userDomainSvc    *services.UserDomainService
}

// NewColorSettingUseCase creates a new ColorSettingUseCase
func NewColorSettingUseCase(
	colorSettingRepo repositories.ColorSettingRepository,
	userDomainSvc *services.UserDomainService,
) *ColorSettingUseCase {
	return &ColorSettingUseCase{
		colorSettingRepo: colorSettingRepo,
		userDomainSvc:    userDomainSvc,
	}
}

// CreateColorSetting creates a new color setting
func (uc *ColorSettingUseCase) CreateColorSetting(ctx context.Context, userID string, req *dto.CreateColorSettingRequest) (*dto.ColorSettingResponse, error) {
	// Validate user exists
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	_, err = uc.userDomainSvc.ValidateUserExists(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Create value objects
	hexColor, err := valueobjects.NewHexColor(req.HexColor)
	if err != nil {
		return nil, err
	}

	// Create color setting entity
	colorSetting, err := entities.NewColorSetting(uid, req.Name, hexColor)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if req.Description != nil {
		colorSetting.UpdateDescription(req.Description)
	}
	if req.SortOrder != nil {
		colorSetting.UpdateSortOrder(*req.SortOrder)
	}

	// Save color setting
	err = uc.colorSettingRepo.Save(ctx, colorSetting)
	if err != nil {
		return nil, err
	}

	return uc.colorSettingToResponse(colorSetting), nil
}

// GetColorSettingsByUserID retrieves all color settings for a user
func (uc *ColorSettingUseCase) GetColorSettingsByUserID(ctx context.Context, userID string) ([]dto.ColorSettingResponse, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	_, err = uc.userDomainSvc.ValidateUserExists(ctx, uid)
	if err != nil {
		return nil, err
	}

	colorSettings, err := uc.colorSettingRepo.FindByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ColorSettingResponse, len(colorSettings))
	for i, cs := range colorSettings {
		responses[i] = *uc.colorSettingToResponse(cs)
	}

	return responses, nil
}

// UpdateColorSetting updates an existing color setting
func (uc *ColorSettingUseCase) UpdateColorSetting(ctx context.Context, userID, settingID string, req *dto.UpdateColorSettingRequest) (*dto.ColorSettingResponse, error) {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	sid, err := valueobjects.NewUserID(settingID)
	if err != nil {
		return nil, err
	}

	// Validate user exists
	_, err = uc.userDomainSvc.ValidateUserExists(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Get existing color setting
	colorSetting, err := uc.colorSettingRepo.FindByID(ctx, sid)
	if err != nil {
		return nil, err
	}
	if colorSetting == nil {
		return nil, errors.New("color setting not found")
	}

	// Verify ownership
	if !colorSetting.BelongsToUser(uid) {
		return nil, errors.New("color setting does not belong to user")
	}

	// Update fields
	if req.Name != nil {
		err = colorSetting.UpdateName(*req.Name)
		if err != nil {
			return nil, err
		}
	}
	if req.HexColor != nil {
		hexColor, err := valueobjects.NewHexColor(*req.HexColor)
		if err != nil {
			return nil, err
		}
		err = colorSetting.UpdateHexColor(hexColor)
		if err != nil {
			return nil, err
		}
	}
	if req.Description != nil {
		colorSetting.UpdateDescription(req.Description)
	}
	if req.SortOrder != nil {
		colorSetting.UpdateSortOrder(*req.SortOrder)
	}

	// Save updated color setting
	err = uc.colorSettingRepo.Update(ctx, colorSetting)
	if err != nil {
		return nil, err
	}

	return uc.colorSettingToResponse(colorSetting), nil
}

// DeleteColorSetting deletes a color setting by ID
func (uc *ColorSettingUseCase) DeleteColorSetting(ctx context.Context, userID, settingID string) error {
	uid, err := valueobjects.NewUserID(userID)
	if err != nil {
		return err
	}

	sid, err := valueobjects.NewUserID(settingID)
	if err != nil {
		return err
	}

	// Validate user exists
	_, err = uc.userDomainSvc.ValidateUserExists(ctx, uid)
	if err != nil {
		return err
	}

	// Get existing color setting
	colorSetting, err := uc.colorSettingRepo.FindByID(ctx, sid)
	if err != nil {
		return err
	}
	if colorSetting == nil {
		return errors.New("color setting not found")
	}

	// Verify ownership
	if !colorSetting.BelongsToUser(uid) {
		return errors.New("color setting does not belong to user")
	}

	return uc.colorSettingRepo.Delete(ctx, sid)
}

// colorSettingToResponse converts a color setting entity to a response DTO
func (uc *ColorSettingUseCase) colorSettingToResponse(cs *entities.ColorSetting) *dto.ColorSettingResponse {
	return &dto.ColorSettingResponse{
		ID:          cs.ID().Value(),
		UserID:      cs.UserID().Value(),
		Name:        cs.Name(),
		HexColor:    cs.HexColor().Value(),
		Description: cs.Description(),
		IsDefault:   cs.IsDefault(),
		SortOrder:   cs.SortOrder(),
		CreatedAt:   cs.CreatedAt(),
		UpdatedAt:   cs.UpdatedAt(),
	}
}
