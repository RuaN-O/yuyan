// file: internal/biz/user.go
package biz

import (
	"context"
	"fmt"

	"time"
	"yuyan/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	UserNotFound      = 701
	UserAlreadyExists = 702
	PasswordError     = 703
	InternalError     = 704
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string // 注意：这个字段不应该对外暴露
	CreateAt     time.Time
	LastLoginAt  time.Time
}

// UserRepo 仓储接口
type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error) // 新增：用于更新最后登录时间
}

type UserUseCase struct {
	jwt  jwt.JwtRepo
	repo UserRepo
}

func NewUserUseCase(repo UserRepo, jwt jwt.JwtRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		jwt:  jwt,
	}
}

// generateToken 生成 JWT Token

// Register 用户注册，返回 JWT Token
// 修改签名，直接接收 username 和 password
func (uc *UserUseCase) Register(ctx context.Context, username, password string) (string, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // 注意：这里需要引入 gorm.io/gorm
		return "", errors.InternalServer("DB_ERROR", fmt.Sprintf("database error: %v", err))
	}
	if existingUser != nil {
		return "", errors.BadRequest("user already exists", "用户已存在")
	}

	// 2. 对密码进行哈希 (从 service 层移过来的逻辑)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.InternalServer("failed to hash password", "密码加密失败")
	}

	// 3. 创建用户模型
	newUser := &User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		CreateAt:     time.Now(),  // 在 biz 层设置创建时间
		LastLoginAt:  time.Time{}, // 注册时，最后登录时间为空
	}

	// 4. 调用仓储层创建用户
	createdUser, err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		return "", errors.InternalServer("failed to create user", "创建用户失败")
	}

	// 5. 生成 Token

	token, err := uc.jwt.GenerateToken(createdUser.ID, createdUser.Username)
	if err != nil {
		return "", errors.InternalServer("failed to generate token", "生成Token失败")
	}

	return token, nil
}

// Login 用户登录，返回 JWT Token
func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, error) {
	// 1. 根据用户名查找用户
	user, err := uc.repo.GetUserByUsername(ctx, username)
	// 2. 修正错误处理逻辑
	if err != nil {
		// 统一处理“未找到”错误，避免暴露用户是否存在
		return "", errors.Unauthorized("user not found or password incorrect", "用户不存在或密码错误")
	}
	if user == nil { // 这部分其实是冗余的，因为 repo 层找不到会返回 error，但保留更健壮
		return "", errors.Unauthorized("user not found or password incorrect", "用户不存在或密码错误")
	}

	// 3. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// 移除不安全的日志，只记录警告
		// log.Printf("数据库密码:%v,传参加密后密码:%v", user.PasswordHash, password)
		// 密码错误，返回一个统一的、模糊的错误信息
		return "", errors.Unauthorized("user not found or password incorrect", "用户不存在或密码错误")
	}

	// 4. 更新最后登录时间
	user.LastLoginAt = time.Now()
	_, err = uc.repo.UpdateUser(ctx, user)
	if err != nil {
		// 即使更新失败，也不应该阻止登录，但应该记录日志
		log.Errorf("failed to update last login time for user %d: %v", user.ID, err)
	}

	// 5. 生成 Token
	token, err := uc.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", errors.InternalServer("failed to generate token", "生成Token失败")
	}

	return token, nil
}
