package convertor

import (
	balancemapper "github.com/CloudStriver/cloudmind-trade/biz/infrastructure/mapper/balance"
	gentrade "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
)

func BalanceMapperToBalance(in *balancemapper.Balance) *gentrade.Balance {
	return &gentrade.Balance{
		Flow:   in.Flow,
		Memory: in.Memory,
		Point:  in.Point,
	}
}

func BalanceToBalanceMapper(in *gentrade.Balance) *balancemapper.Balance {
	return &balancemapper.Balance{
		Flow:   in.Flow,
		Memory: in.Memory,
		Point:  in.Point,
	}
}
