package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	balancemapper "github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/balance"
	gentrade "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BalanceService interface {
	UpdateBalance(ctx context.Context, req *gentrade.UpdateBalanceReq) (resp *gentrade.UpdateBalanceResp, err error)
	GetBalance(ctx context.Context, req *gentrade.GetBalanceReq) (resp *gentrade.GetBalanceResp, err error)
	CreateBalance(ctx context.Context, req *gentrade.CreateBalanceReq) (resp *gentrade.CreateBalanceResp, err error)
}

var BalanceSet = wire.NewSet(
	wire.Struct(new(BalanceServiceImpl), "*"),
	wire.Bind(new(BalanceService), new(*BalanceServiceImpl)),
)

type BalanceServiceImpl struct {
	Config             *config.Config
	Redis              *redis.Redis
	BalanceMongoMapper balancemapper.IBalanceMongoMapper
}

func (s *BalanceServiceImpl) UpdateBalance(ctx context.Context, req *gentrade.UpdateBalanceReq) (resp *gentrade.UpdateBalanceResp, err error) {
	resp = new(gentrade.UpdateBalanceResp)
	oid, _ := primitive.ObjectIDFromHex(req.UserId)
	result, err := s.BalanceMongoMapper.Update(ctx, &balancemapper.Balance{
		ID:     oid,
		Flow:   req.Flow,
		Memory: req.Memory,
		Point:  req.Point,
	})
	if err != nil {
		return resp, err
	}
	resp.Ok = result.ModifiedCount != 0
	return resp, nil
}

func (s *BalanceServiceImpl) GetBalance(ctx context.Context, req *gentrade.GetBalanceReq) (resp *gentrade.GetBalanceResp, err error) {
	balance, err := s.BalanceMongoMapper.FindOne(ctx, req.UserId)
	if err != nil {
		return resp, err
	}
	return &gentrade.GetBalanceResp{
		Flow:   *balance.Flow,
		Memory: *balance.Memory,
		Point:  *balance.Point,
	}, nil
}

func (s *BalanceServiceImpl) CreateBalance(ctx context.Context, req *gentrade.CreateBalanceReq) (resp *gentrade.CreateBalanceResp, err error) {
	oid, _ := primitive.ObjectIDFromHex(req.UserId)
	if _, err = s.BalanceMongoMapper.Insert(ctx, &balancemapper.Balance{
		ID:     oid,
		Flow:   lo.ToPtr(s.Config.Balance.DefaultFlow),
		Memory: lo.ToPtr(s.Config.Balance.DefaultMemory),
		Point:  lo.ToPtr(s.Config.Balance.DefaultPoint),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
