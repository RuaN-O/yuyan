// file: internal/service/user.go
package service

import (
	"context"
	"errors"
	v1 "yuyan/api/user/v1"
	"yuyan/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
	uc *biz.UserUseCase
}

func NewUserService(uc *biz.UserUseCase) *UserService {
	return &UserService{uc: uc}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	if req.Password == "" || req.Username == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	// 直接调用 biz 层，传递原始的用户名和密码
	// 注意：这里 biz.Register 的签名需要修改
	token, err := s.uc.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterReply{
		Token: token, // 返回 biz 层生成的 token
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("用户名和密码不能为空")
	}

	// 直接调用 biz 层，传递原始的用户名和密码
	token, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &v1.LoginReply{
		Token: token, // 返回 biz 层生成的 token
	}, nil
}

// GetUser ... (暂时保持不变)
func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserReply, error) {
	return nil, nil
}
