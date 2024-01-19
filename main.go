package main

import (
	"github.com/CloudStriver/cloudmind-trade/provider"
	"github.com/CloudStriver/go-pkg/utils/kitex/middleware"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade/tradeservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"net"
)

func main() {
	klog.SetLogger(log.NewKlogLogger())
	s, err := provider.NewTradeServerImpl()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", s.ListenOn)
	if err != nil {
		panic(err)
	}

	svr := tradeservice.NewServer(
		s,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: s.Name}),
		server.WithMiddleware(middleware.LogMiddleware(s.Name)),
	)

	err = svr.Run()
	if err != nil {
		log.Error(err.Error())
	}
}
