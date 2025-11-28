package data

import (
	"context"
	"errors"
	"yuyan/internal/biz"
	data "yuyan/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) toDataUser(b *biz.User) *data.User {
	if b == nil {
		return nil
	}
	return &data.User{
		ID:           b.ID,
		Username:     b.Username,
		PasswordHash: b.PasswordHash,
		CreateTime:   b.CreateAt,
		LastLogin:    b.LastLoginAt,
	}

}

func (r *userRepo) toBizUser(d *data.User) *biz.User {
	if d == nil {
		return nil
	}
	return &biz.User{
		ID:           d.ID,
		Username:     d.Username,
		PasswordHash: d.PasswordHash,
		CreateAt:     d.CreateTime,
		LastLoginAt:  d.LastLogin,
	}
}

// CreateUser 创建新用户
func (r *userRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	// 1. 将业务模型转换为数据模型
	dataUser := r.toDataUser(user)
	// GORM 会自动处理 ID (自增) 和 CreateTime (自动创建时间)

	// 2. 使用 GORM 创建用户
	if err := r.data.db.WithContext(ctx).Create(dataUser).Error; err != nil {
		r.log.Errorf("create user failed: %v", err)
		return nil, err
	}

	// 3. 将创建后的数据模型（包含ID）转换回业务模型返回
	r.log.Infof("created user successfully, id: %d, username: %s", dataUser.ID, dataUser.Username)
	return r.toBizUser(dataUser), nil
}

// DeleteUser 根据 ID 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	// 使用 GORM 删除用户
	result := r.data.db.WithContext(ctx).Delete(&data.User{}, id) // &User{} 指定表名，id 指定主键

	if err := result.Error; err != nil {
		r.log.Errorf("delete user by id %d failed: %v", id, err)
		return err
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		r.log.Warnf("attempted to delete user with id %d, but it does not exist", id)
		// 可以返回一个特定的“未找到”错误，或者直接返回 nil，取决于业务需求
		// 这里我们返回 gorm 的 ErrRecordNotFound 以保持一致性
		return gorm.ErrRecordNotFound
	}

	r.log.Infof("deleted user successfully, id: %d", id)
	return nil
}

// GetUser 根据 ID 获取用户
func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	// 创建一个数据模型实例来接收查询结果
	var dataUser data.User
	// 使用 GORM 查询用户
	if err := r.data.db.WithContext(ctx).First(&dataUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("user with id %d not found", id)
		} else {
			r.log.Errorf("get user by id %d failed: %v", id, err)
		}
		return nil, err
	}

	// 将数据模型转换为业务模型返回
	return r.toBizUser(&dataUser), nil
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	// 创建一个数据模型实例来接收查询结果
	var dataUser data.User
	// 使用 GORM 的 Where 子句查询
	if err := r.data.db.WithContext(ctx).Where("username = ?", username).First(&dataUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("user with username %s not found", username)
		} else {
			r.log.Errorf("get user by username %s failed: %v", username, err)
		}
		return nil, err
	}

	// 将数据模型转换为业务模型返回
	return r.toBizUser(&dataUser), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	dataUser := r.toDataUser(user)

	// 使用 Updates 更新指定字段，更安全
	result := r.data.db.WithContext(ctx).Model(&data.User{}).Where("user_id = ?", dataUser.ID).Updates(map[string]interface{}{
		"last_login": dataUser.LastLogin,
	})

	if err := result.Error; err != nil {
		r.log.Errorf("update user by id %d failed: %v", dataUser.ID, err)
		return nil, err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("attempted to update user with id %d, but it does not exist", dataUser.ID)
		return nil, gorm.ErrRecordNotFound
	}

	// 返回更新后的用户信息
	return r.GetUser(ctx, dataUser.ID)
}
