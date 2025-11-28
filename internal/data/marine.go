package data

import (
	"context"
	"errors"
	"yuyan/internal/biz"
	data "yuyan/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// marineCreatureRepo 是 biz.MarineCreatureRepo 接口的实现
// 它负责与数据库进行交互，操作 marine_creature 表
type marineCreatureRepo struct {
	data *Data
	log  *log.Helper
}

// NewMarineCreatureRepo 是构造函数，用于创建 marineCreatureRepo 实例
// 注意：这里假设您已经在 internal/biz 目录下定义了 MarineCreatureRepo 接口
func NewMarineCreatureRepo(data *Data, logger log.Logger) biz.MarineCreatureRepo {
	return &marineCreatureRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// toDataCreature 将业务层的 MarineCreature 模型转换为数据层的模型
func (r *marineCreatureRepo) toDataCreature(b *biz.MarineCreature) *data.MarineCreature {
	if b == nil {
		return nil
	}
	return &data.MarineCreature{
		// ID 是自增的，创建时不需要设置
		Name:        b.Name,
		Description: b.Description,
		ImageURL:    b.ImageURL,
		// CreateTime 会被 GORM 的 autoCreateTime 自动处理
	}
}

// toBizCreature 将数据层的 MarineCreature 模型转换为业务层的模型
func (r *marineCreatureRepo) toBizCreature(d *data.MarineCreature) *biz.MarineCreature {
	if d == nil {
		return nil
	}
	return &biz.MarineCreature{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		ImageURL:    d.ImageURL,
		CreateAt:    d.CreateTime, // 假设业务层模型字段名为 CreateAt
	}
}

// CreateCreature 创建新的海洋生物记录
func (r *marineCreatureRepo) CreateCreature(ctx context.Context, creature *biz.MarineCreature) (*biz.MarineCreature, error) {
	// 1. 将业务模型转换为数据模型
	dataCreature := r.toDataCreature(creature)

	// 2. 使用 GORM 创建记录
	if err := r.data.db.WithContext(ctx).Create(dataCreature).Error; err != nil {
		r.log.Errorf("create marine creature failed: %v", err)
		return nil, err
	}

	// 3. 将创建后的数据模型（包含ID）转换回业务模型返回
	r.log.Infof("created marine creature successfully, id: %d, name: %s", dataCreature.ID, dataCreature.Name)
	return r.toBizCreature(dataCreature), nil
}

// GetCreature 根据 ID 获取海洋生物
func (r *marineCreatureRepo) GetCreature(ctx context.Context, id int) (*biz.MarineCreature, error) {
	// 创建一个数据模型实例来接收查询结果
	var dataCreature data.MarineCreature
	// 使用 GORM 查询
	if err := r.data.db.WithContext(ctx).First(&dataCreature, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("marine creature with id %d not found", id)
		} else {
			r.log.Errorf("get marine creature by id %d failed: %v", id, err)
		}
		return nil, err
	}

	// 将数据模型转换为业务模型返回
	return r.toBizCreature(&dataCreature), nil
}

// GetCreatureByName 根据名称获取海洋生物
func (r *marineCreatureRepo) GetCreatureByName(ctx context.Context, name string) (*biz.MarineCreature, error) {
	// 创建一个数据模型实例来接收查询结果
	var dataCreature data.MarineCreature
	// 使用 GORM 的 Where 子句查询
	if err := r.data.db.WithContext(ctx).Where("name = ?", name).First(&dataCreature).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("marine creature with name %s not found", name)
		} else {
			r.log.Errorf("get marine creature by name %s failed: %v", name, err)
		}
		return nil, err
	}

	// 将数据模型转换为业务模型返回
	return r.toBizCreature(&dataCreature), nil
}

// UpdateCreature 更新海洋生物信息
func (r *marineCreatureRepo) UpdateCreature(ctx context.Context, creature *biz.MarineCreature) (*biz.MarineCreature, error) {
	dataCreature := r.toDataCreature(creature)
	dataCreature.ID = creature.ID // 确保ID被设置用于更新

	// 使用 Updates 更新指定字段，更安全。它会忽略零值字段。
	result := r.data.db.WithContext(ctx).Model(&data.MarineCreature{}).Where("creature_id = ?", dataCreature.ID).Updates(dataCreature)

	if err := result.Error; err != nil {
		r.log.Errorf("update marine creature by id %d failed: %v", dataCreature.ID, err)
		return nil, err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("attempted to update marine creature with id %d, but it does not exist", dataCreature.ID)
		return nil, gorm.ErrRecordNotFound
	}

	r.log.Infof("updated marine creature successfully, id: %d", dataCreature.ID)
	// 返回更新后的用户信息
	return r.GetCreature(ctx, dataCreature.ID)
}

// DeleteCreature 根据 ID 删除海洋生物
func (r *marineCreatureRepo) DeleteCreature(ctx context.Context, id int) error {
	// 使用 GORM 删除记录
	result := r.data.db.WithContext(ctx).Delete(&data.MarineCreature{}, id) // &MarineCreature{} 指定表名，id 指定主键

	if err := result.Error; err != nil {
		r.log.Errorf("delete marine creature by id %d failed: %v", id, err)
		return err
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		r.log.Warnf("attempted to delete marine creature with id %d, but it does not exist", id)
		return gorm.ErrRecordNotFound
	}

	r.log.Infof("deleted marine creature successfully, id: %d", id)
	return nil
}
