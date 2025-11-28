// file: internal/data/history.go
package data

import (
	"context"
	"errors"
	"yuyan/internal/biz"
	data "yuyan/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// historyRepo 是 biz.HistoryRepo 接口的实现
type historyRepo struct {
	data *Data
	log  *log.Helper
}

// NewHistoryRepo 是构造函数，用于创建 historyRepo 实例
func NewHistoryRepo(data *Data, logger log.Logger) biz.HistoryRepo {
	return &historyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// toDataHistory 将业务层的 History 模型转换为数据层的模型
func (r *historyRepo) toDataHistory(b *biz.History) *data.History {
	if b == nil {
		return nil
	}
	return &data.History{
		ID:         b.ID,
		UserID:     b.UserID,
		FileName:   b.FileName,
		CreatureID: b.CreatureID,
		ResultText: b.ResultText,
		// IdentifyTime 会被 GORM 的 autoCreateTime 自动处理
	}
}

// toBizHistory 将数据层的 History 模型转换为业务层的模型
func (r *historyRepo) toBizHistory(d *data.History) *biz.History {
	if d == nil {
		return nil
	}
	return &biz.History{
		ID:           d.ID,
		UserID:       d.UserID,
		FileName:     d.FileName,
		CreatureID:   d.CreatureID,
		ResultText:   d.ResultText,
		IdentifyTime: d.IdentifyTime,
	}
}

// CreateHistory 创建新的历史记录
func (r *historyRepo) CreateHistory(ctx context.Context, history *biz.History) (*biz.History, error) {
	// 1. 将业务模型转换为数据模型
	dataHistory := r.toDataHistory(history)

	// 2. 使用 GORM 创建记录
	if err := r.data.db.WithContext(ctx).Create(dataHistory).Error; err != nil {
		r.log.Errorf("create history failed: %v", err)
		return nil, err
	}

	// 3. 将创建后的数据模型（包含ID）转换回业务模型返回
	r.log.Infof("created history successfully, id: %d, user_id: %d", dataHistory.ID, dataHistory.UserID)
	return r.toBizHistory(dataHistory), nil
}

// GetHistoryByID 根据 ID 获取历史记录
func (r *historyRepo) GetHistoryByID(ctx context.Context, id int64) (*biz.History, error) {
	// 创建一个数据模型实例来接收查询结果
	var dataHistory data.History
	// 使用 GORM 查询
	if err := r.data.db.WithContext(ctx).First(&dataHistory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("history with id %d not found", id)
		} else {
			r.log.Errorf("get history by id %d failed: %v", id, err)
		}
		return nil, err
	}

	// 将数据模型转换为业务模型返回
	return r.toBizHistory(&dataHistory), nil
}

// GetHistoryByUserID 根据用户ID获取其所有历史记录
func (r *historyRepo) GetHistoryByUserID(ctx context.Context, userID int64) ([]*biz.History, error) {
	// 创建一个数据模型切片来接收查询结果
	var dataHistories []data.History
	// 使用 GORM 的 Where 子句查询
	if err := r.data.db.WithContext(ctx).Where("user_id = ?", userID).Find(&dataHistories).Error; err != nil {
		r.log.Errorf("get histories by user_id %d failed: %v", userID, err)
		return nil, err
	}

	// 将数据模型切片转换为业务模型切片返回
	bizHistories := make([]*biz.History, 0, len(dataHistories))
	for _, dh := range dataHistories {
		bizHistories = append(bizHistories, r.toBizHistory(&dh))
	}

	return bizHistories, nil
}
