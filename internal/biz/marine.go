// file: internal/biz/marine.go
package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

// 业务错误码
const (
	CreatureNotFound      = "creature_not_found"
	CreatureAlreadyExists = "creature_already_exists"
	HistoryNotFound       = "history_not_found"
	InnternalError        = "InternalError"
)

// MarineCreature 海洋生物的业务模型
type MarineCreature struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreateAt    time.Time `json:"create_at"`
}

// History 历史记录的业务模型
type History struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	FileName     string    `json:"file_name"`
	CreatureID   int       `json:"creature_id"`
	ResultText   string    `json:"result_text"`
	IdentifyTime time.Time `json:"identify_time"`
}

// MarineCreatureRepo 定义海洋生物数据访问的接口
type MarineCreatureRepo interface {
	CreateCreature(ctx context.Context, creature *MarineCreature) (*MarineCreature, error)
	GetCreature(ctx context.Context, id int) (*MarineCreature, error)
	GetCreatureByName(ctx context.Context, name string) (*MarineCreature, error)
	UpdateCreature(ctx context.Context, creature *MarineCreature) (*MarineCreature, error)
	DeleteCreature(ctx context.Context, id int) error
}

// HistoryRepo 定义历史记录数据访问的接口
type HistoryRepo interface {
	CreateHistory(ctx context.Context, history *History) (*History, error)
	GetHistoryByID(ctx context.Context, id int64) (*History, error)
	GetHistoryByUserID(ctx context.Context, userID int64) ([]*History, error)
}

// MarineUseCase 海洋生物相关的业务逻辑集合
type MarineUseCase struct {
	creatureRepo MarineCreatureRepo
	historyRepo  HistoryRepo
}

// NewMarineUseCase 创建 MarineUseCase
func NewMarineUseCase(creatureRepo MarineCreatureRepo, historyRepo HistoryRepo) *MarineUseCase {
	return &MarineUseCase{
		creatureRepo: creatureRepo,
		historyRepo:  historyRepo,
	}
}

// --- MarineCreature 相关的业务逻辑 ---

// AddCreature 添加新的海洋生物
func (uc *MarineUseCase) AddCreature(ctx context.Context, name, description, imageURL string) (*MarineCreature, error) {
	// 1. 检查生物是否已存在
	existing, err := uc.creatureRepo.GetCreatureByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}
	if existing != nil {
		return nil, errors.BadRequest(CreatureAlreadyExists, "该生物已存在")
	}

	// 2. 创建业务模型
	newCreature := &MarineCreature{
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		CreateAt:    time.Now(),
	}

	// 3. 调用仓储层创建
	return uc.creatureRepo.CreateCreature(ctx, newCreature)
}

// UpdateCreature 更新海洋生物信息
func (uc *MarineUseCase) UpdateCreature(ctx context.Context, id int, name, description, imageURL string) (*MarineCreature, error) {
	// 1. 检查生物是否存在
	_, err := uc.creatureRepo.GetCreature(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(CreatureNotFound, "生物未找到")
		}
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}

	// 2. 构建更新模型
	updatedCreature := &MarineCreature{
		ID:          id,
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
	}

	// 3. 调用仓储层更新
	return uc.creatureRepo.UpdateCreature(ctx, updatedCreature)
}

// DeleteByCreatureId 删除海洋生物
func (uc *MarineUseCase) DeleteByCreatureId(ctx context.Context, id int) error {
	// 1. 检查生物是否存在
	_, err := uc.creatureRepo.GetCreature(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NotFound(CreatureNotFound, "生物未找到")
		}
		return errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}

	// 2. 调用仓储层删除
	return uc.creatureRepo.DeleteCreature(ctx, id)
}

// GetCreatureByCreatureId 根据 ID 获取海洋生物
func (uc *MarineUseCase) GetCreatureByCreatureId(ctx context.Context, id int) (*MarineCreature, error) {
	creature, err := uc.creatureRepo.GetCreature(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(CreatureNotFound, "生物未找到")
		}
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}
	return creature, nil
}

// --- History 相关的业务逻辑 ---

// AddHistory 添加一条识别历史记录
func (uc *MarineUseCase) AddHistory(ctx context.Context, userID int64, fileName string, creatureID int, resultText string) (*History, error) {
	// 1. 验证 creatureID 是否有效
	_, err := uc.creatureRepo.GetCreature(ctx, creatureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(CreatureNotFound, "关联的生物ID无效")
		}
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}

	// 2. 创建业务模型
	newHistory := &History{
		UserID:       userID,
		FileName:     fileName,
		CreatureID:   creatureID,
		ResultText:   resultText,
		IdentifyTime: time.Now(),
	}

	// 3. 调用仓储层创建
	return uc.historyRepo.CreateHistory(ctx, newHistory)
}

// GetHistoryByHistoryId 根据 ID 获取历史记录
func (uc *MarineUseCase) GetHistoryByHistoryId(ctx context.Context, historyID int64) (*History, error) {
	history, err := uc.historyRepo.GetHistoryByID(ctx, historyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(HistoryNotFound, "历史记录未找到")
		}
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}
	return history, nil
}

// GetHistoryByUser 根据用户ID获取其所有历史记录
func (uc *MarineUseCase) GetHistoryByUser(ctx context.Context, userID int64) ([]*History, error) {
	// 这个方法通常不需要复杂的错误处理，因为“没有记录”是一个有效的空结果
	histories, err := uc.historyRepo.GetHistoryByUserID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer(InnternalError, fmt.Sprintf("database error: %v", err))
	}
	return histories, nil
}
