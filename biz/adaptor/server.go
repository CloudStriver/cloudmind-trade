package adaptor

import (
	"context"
	"github.com/CloudStriver/cloudmind-trade/biz/application/service"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
)

type TradeServerImpl struct {
	*config.Config
	BalanceService service.BalanceService
	StockService   service.StockService
}

func (s *TradeServerImpl) UpdateBalance(ctx context.Context, req *trade.UpdateBalanceReq) (resp *trade.UpdateBalanceResp, err error) {
	return s.BalanceService.UpdateBalance(ctx, req)
}

func (s *TradeServerImpl) CreateBalance(ctx context.Context, req *trade.CreateBalanceReq) (resp *trade.CreateBalanceResp, err error) {
	return s.BalanceService.CreateBalance(ctx, req)
}

func (s *TradeServerImpl) GetBalance(ctx context.Context, req *trade.GetBalanceReq) (resp *trade.GetBalanceResp, err error) {
	return s.BalanceService.GetBalance(ctx, req)
}

func (s *TradeServerImpl) GetStock(ctx context.Context, req *trade.GetStockReq) (resp *trade.GetStockResp, err error) {
	return s.StockService.GetStock(ctx, req)
}

func (s *TradeServerImpl) GetStocks(ctx context.Context, req *trade.GetStocksReq) (resp *trade.GetStocksResp, err error) {
	return s.StockService.GetStocks(ctx, req)
}

func (s *TradeServerImpl) AddStock(ctx context.Context, req *trade.AddStockReq) (resp *trade.AddStockResp, err error) {
	return s.StockService.AddStock(ctx, req)
}
