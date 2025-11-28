// file: internal/service/marine.go
package service

import (
	"context"

	pb "yuyan/api/marine/v1"
	"yuyan/internal/biz"
	"yuyan/internal/middleware/auth" // 引入 auth 包来获取 userID

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
)

// MarineServiceService 是 marine service 的实现
type MarineServiceService struct {
	pb.UnimplementedMarineServiceServer
	uc  *biz.MarineUseCase
	log *log.Helper
}

// NewMarineServiceService 创建 MarineServiceService 实例
func NewMarineServiceService(uc *biz.MarineUseCase, logger log.Logger) *MarineServiceService {
	return &MarineServiceService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// --- MarineCreature 相关的 RPC 方法 ---

func (s *MarineServiceService) AddCreature(ctx context.Context, req *pb.AddCreatureRequest) (*pb.AddCreatureResponse, error) {
	creature, err := s.uc.AddCreature(ctx, req.GetName(), req.GetDescription(), req.GetImageUrl())
	if err != nil {
		return nil, err // Kratos 会自动处理错误转换为 gRPC status
	}
	return &pb.AddCreatureResponse{
		CreatureId: uint32(creature.ID),
	}, nil
}

func (s *MarineServiceService) UpdateCreature(ctx context.Context, req *pb.UpdateCreatureRequest) (*pb.UpdateCreatureResponse, error) {
	_, err := s.uc.UpdateCreature(ctx, int(req.GetId()), req.GetCreature().Name, req.GetCreature().Description, req.GetCreature().ImageUrl)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCreatureResponse{
		Creature: req.Creature,
	}, nil
}

func (s *MarineServiceService) DeleteByCreatureId(ctx context.Context, req *pb.DeleteByCreatureIdRequest) (*empty.Empty, error) {
	err := s.uc.DeleteByCreatureId(ctx, int(req.Creatureid))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *MarineServiceService) GetCreatureByCreatureId(ctx context.Context, req *pb.GetCreatureByCreatureIdRequest) (*pb.GetCreatureByCreatureIdResponse, error) {
	creature, err := s.uc.GetCreatureByCreatureId(ctx, int(req.Creatureid))
	if err != nil {
		return nil, err
	}
	return &pb.GetCreatureByCreatureIdResponse{
		Creature: &pb.MarineCreature{
			CreatureId:  uint32(creature.ID),
			Name:        creature.Name,
			ImageUrl:    creature.ImageURL,
			Description: creature.Description,
			CreateTime:  creature.CreateAt.Unix(),
		},
	}, nil
}

// --- History 相关的 RPC 方法 ---

func (s *MarineServiceService) AddHistory(ctx context.Context, req *pb.AddHistoryRequest) (*pb.AddHistoryResponse, error) {
	// 从 context 中获取 user_id
	userID, ok := auth.FromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("USER_NOT_FOUND", "用户未认证")
	}

	history, err := s.uc.AddHistory(ctx, userID, req.GetFileName(), int(req.GetCreatureId()), req.GetResultText())
	if err != nil {
		return nil, err
	}
	return &pb.AddHistoryResponse{
		HistoryId: uint64(history.ID),
	}, nil
}

func (s *MarineServiceService) GetHistoryByHistoryId(ctx context.Context, req *pb.GetHistoryByHistoryIdRequest) (*pb.GetHistoryByHistoryIdResponse, error) {
	history, err := s.uc.GetHistoryByHistoryId(ctx, int64(req.Historyid))
	if err != nil {
		return nil, err
	}
	userid, ok := ctx.Value("user_id").(int64)
	if !ok {
		userid = 0
	}
	return &pb.GetHistoryByHistoryIdResponse{
		History: &pb.History{
			HistoryId:    uint64(history.ID),
			UserId:       uint64(userid),
			FileName:     history.FileName,
			CreatureId:   uint32(history.CreatureID),
			ResultText:   history.ResultText,
			IdentifyTime: history.IdentifyTime.Unix(),
		},
	}, nil
}

func (s *MarineServiceService) GetHistoryByUser(ctx context.Context, req *empty.Empty) (*pb.GetHistoryByUserResponse, error) {
	// 从 context 中获取 user_id
	userID, ok := auth.FromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("USER_NOT_FOUND", "用户未认证")
	}

	histories, err := s.uc.GetHistoryByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为 pb 切片
	pbHistories := make([]*pb.History, 0, len(histories))
	for _, h := range histories {
		pbHistories = append(pbHistories, &pb.History{
			HistoryId:    uint64(h.ID),
			UserId:       uint64(h.UserID),
			FileName:     h.FileName,
			CreatureId:   uint32(h.CreatureID),
			ResultText:   h.ResultText,
			IdentifyTime: h.IdentifyTime.Unix(),
		})
	}

	return &pb.GetHistoryByUserResponse{
		Histories: pbHistories,
	}, nil
}
