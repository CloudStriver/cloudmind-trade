package provider

import (
	"github.com/CloudStriver/cloudmind-trade/biz/application/service"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/balance"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/stock"
	"github.com/CloudStriver/cloudmind-trade/biz/infrastructure/stores/redis"
	"github.com/google/wire"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.BalanceSet,
	service.StockSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	MapperSet,
)

var MapperSet = wire.NewSet(
	stock.NewMongoMapper,
	balance.NewMongoMapper,
)
