// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"github.com/CloudStriver/cloudmind-trade/biz/adaptor"
	"github.com/CloudStriver/cloudmind-trade/biz/application/service"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/balance"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/stock"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/stores/redis"
)

// Injectors from wire.go:

func NewTradeServerImpl() (*adaptor.TradeServerImpl, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	redisRedis := redis.NewRedis(configConfig)
	iBalanceMongoMapper := balance.NewMongoMapper(configConfig)
	balanceServiceImpl := &service.BalanceServiceImpl{
		Config:             configConfig,
		Redis:              redisRedis,
		BalanceMongoMapper: iBalanceMongoMapper,
	}
	iStockMongoMapper := stock.NewMongoMapper(configConfig)
	stockServiceImpl := &service.StockServiceImpl{
		Config:           configConfig,
		Redis:            redisRedis,
		StockMongoMapper: iStockMongoMapper,
	}
	tradeServerImpl := &adaptor.TradeServerImpl{
		Config:         configConfig,
		BalanceService: balanceServiceImpl,
		StockService:   stockServiceImpl,
	}
	return tradeServerImpl, nil
}