//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/CloudStriver/cloudmind-trade/biz/adaptor"
	"github.com/google/wire"
)

func NewTradeServerImpl() (*adaptor.TradeServerImpl, error) {
	wire.Build(
		wire.Struct(new(adaptor.TradeServerImpl), "*"),
		AllProvider,
	)
	return nil, nil
}
