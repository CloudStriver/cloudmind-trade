package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	stockmapper "github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/stock"
	gentrade "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockService interface {
	GetStocks(ctx context.Context, req *gentrade.GetStocksReq) (resp *gentrade.GetStocksResp, err error)
	GetStock(ctx context.Context, req *gentrade.GetStockReq) (resp *gentrade.GetStockResp, err error)
	UpdateStock(ctx context.Context, req *gentrade.UpdateStockReq) (resp *gentrade.UpdateStockResp, err error)
	CreateStock(ctx context.Context, req *gentrade.CreateStockReq) (resp *gentrade.CreateStockResp, err error)
}

var StockSet = wire.NewSet(
	wire.Struct(new(StockServiceImpl), "*"),
	wire.Bind(new(StockService), new(*StockServiceImpl)),
)

type StockServiceImpl struct {
	Config           *config.Config
	Redis            *redis.Redis
	StockMongoMapper stockmapper.IStockMongoMapper
}

func (s *StockServiceImpl) CreateStock(ctx context.Context, req *gentrade.CreateStockReq) (resp *gentrade.CreateStockResp, err error) {
	resp = new(gentrade.CreateStockResp)
	oid, _ := primitive.ObjectIDFromHex(req.ProductId)
	if _, err = s.StockMongoMapper.Insert(ctx, &stockmapper.Stock{
		ID:     oid,
		Amount: lo.ToPtr(req.Stock),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *StockServiceImpl) GetStocks(ctx context.Context, req *gentrade.GetStocksReq) (resp *gentrade.GetStocksResp, err error) {
	resp = new(gentrade.GetStocksResp)
	stocks, err := s.StockMongoMapper.FindMany(ctx, req.ProductIds)
	if err != nil {
		return resp, err
	}

	index := make(map[string]int, len(req.ProductIds))
	lo.ForEach[string](req.ProductIds, func(v string, i int) {
		index[v] = i
	})

	resp.Stocks = make([]int64, len(req.ProductIds), len(req.ProductIds))
	lo.ForEach[*stockmapper.Stock](stocks, func(v *stockmapper.Stock, _ int) {
		resp.Stocks[index[v.ID.Hex()]] = *v.Amount
	})
	return resp, nil
}

func (s *StockServiceImpl) GetStock(ctx context.Context, req *gentrade.GetStockReq) (resp *gentrade.GetStockResp, err error) {
	resp = new(gentrade.GetStockResp)
	stock, err := s.StockMongoMapper.FindOne(ctx, req.ProductId)
	if err != nil {
		return resp, err
	}
	resp.Stock = *stock.Amount
	return resp, nil
}

func (s *StockServiceImpl) UpdateStock(ctx context.Context, req *gentrade.UpdateStockReq) (resp *gentrade.UpdateStockResp, err error) {
	resp = new(gentrade.UpdateStockResp)
	oid, _ := primitive.ObjectIDFromHex(req.ProductId)
	if _, err = s.StockMongoMapper.Update(ctx, &stockmapper.Stock{
		ID:     oid,
		Amount: lo.ToPtr(req.Amount),
	}, lo.ToPtr(req.OldAmount)); err != nil {
		return resp, err
	}
	return resp, nil
}
