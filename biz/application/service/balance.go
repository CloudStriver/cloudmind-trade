package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/convertor"
	balancemapper "github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/balance"
	gentrade "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/google/wire"
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
	oldBalance := convertor.BalanceToBalanceMapper(req.OldBalance)
	oldBalance.ID = oid
	if _, err = s.BalanceMongoMapper.Update(ctx, convertor.BalanceToBalanceMapper(req.Balance), oldBalance); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *BalanceServiceImpl) GetBalance(ctx context.Context, req *gentrade.GetBalanceReq) (resp *gentrade.GetBalanceResp, err error) {
	resp = new(gentrade.GetBalanceResp)
	stock, err := s.BalanceMongoMapper.FindOne(ctx, req.UserId)
	if err != nil {
		return resp, err
	}
	resp.Balance = convertor.BalanceMapperToBalance(stock)
	return resp, nil
}

func (s *BalanceServiceImpl) CreateBalance(ctx context.Context, req *gentrade.CreateBalanceReq) (resp *gentrade.CreateBalanceResp, err error) {
	resp = new(gentrade.CreateBalanceResp)
	oid, _ := primitive.ObjectIDFromHex(req.UserId)
	balance := convertor.BalanceToBalanceMapper(req.Balance)
	balance.ID = oid

	if _, err = s.BalanceMongoMapper.Insert(ctx, balance); err != nil {
		return resp, err
	}
	return resp, nil
}
