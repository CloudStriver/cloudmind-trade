package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	gentrade "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

type StockService interface {
	GetStocks(ctx context.Context, req *gentrade.GetStocksReq) (resp *gentrade.GetStocksResp, err error)
	GetStock(ctx context.Context, req *gentrade.GetStockReq) (resp *gentrade.GetStockResp, err error)
	AddStock(ctx context.Context, req *gentrade.AddStockReq) (resp *gentrade.AddStockResp, err error)
}

var StockSet = wire.NewSet(
	wire.Struct(new(StockServiceImpl), "*"),
	wire.Bind(new(StockService), new(*StockServiceImpl)),
)

type StockServiceImpl struct {
	Config *config.Config
	Redis  *redis.Redis
}

func (s *StockServiceImpl) GetStocks(ctx context.Context, req *gentrade.GetStocksReq) (resp *gentrade.GetStocksResp, err error) {
	resp = new(gentrade.GetStocksResp)

	req.ProductIds = lo.Map(req.ProductIds, func(item string, _ int) string {
		return fmt.Sprintf("stock:%s", item)
	})

	stocks, err := s.Redis.MgetCtx(ctx, req.ProductIds...)
	if err != nil {
		return resp, err
	}

	resp.Stocks = lo.Map[string, int64](stocks, func(item string, _ int) int64 {
		if item == "" {
			return 0
		} else {
			stock, _ := strconv.ParseInt(item, 10, 64)
			return stock
		}
	})
	return resp, nil
}

func (s *StockServiceImpl) GetStock(ctx context.Context, req *gentrade.GetStockReq) (resp *gentrade.GetStockResp, err error) {
	resp = new(gentrade.GetStockResp)
	val, err := s.Redis.GetCtx(ctx, fmt.Sprintf("stock:%s", req.ProductId))
	if err != nil {
		return resp, err
	}
	if val == "" {
		resp.Stock = 0
		return resp, nil
	}
	resp.Stock, _ = strconv.ParseInt(val, 10, 64)
	return resp, nil
}

func (s *StockServiceImpl) AddStock(ctx context.Context, req *gentrade.AddStockReq) (resp *gentrade.AddStockResp, err error) {
	resp = new(gentrade.AddStockResp)
	if _, err := s.Redis.IncrbyCtx(ctx, fmt.Sprintf("stock:%s", req.ProductId), req.Amount); err != nil {
		return resp, err
	}
	return resp, nil
}
