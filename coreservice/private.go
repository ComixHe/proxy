package coreservice

import (
	"github.com/ComixHe/proxy/dbusmon"
	"github.com/black-desk/cgtproxy/pkg/cgtproxy"
	"github.com/black-desk/cgtproxy/pkg/cgtproxy/config"
	"github.com/black-desk/cgtproxy/pkg/nftman"
	"github.com/black-desk/cgtproxy/pkg/nftman/connector"
	"github.com/black-desk/cgtproxy/pkg/routeman"
	"github.com/black-desk/cgtproxy/pkg/types"
	"go.uber.org/zap"
)

func getCGtproxy(cfg *config.Config, eventsChan chan types.CGroupEvents, log *zap.SugaredLogger) (*cgtproxy.CGTProxy, error) {
	root := cfg.CgroupRoot
	bypass := cfg.Bypass

	netlik, err := connector.New()
	if err != nil {
		return nil, err
	}

	nftManager, err := nftman.New(nftman.WithCgroupRoot(root),
		nftman.WithBypass(bypass),
		nftman.WithLogger(log),
		nftman.WithConnFactory(netlik))
	if err != nil {
		return nil, err
	}

	cGroupMonitor, err := dbusmon.New(dbusmon.WithLogger(log), dbusmon.WithChannel(eventsChan))
	if err != nil {
		return nil, err
	}

	routeManager, err := routeman.New(routeman.WithNFTMan(nftManager),
		routeman.WithConfig(cfg),
		routeman.WithCGroupEventChan(cGroupMonitor.Events()),
		routeman.WithLogger(log),
	)
	if err != nil {
		return nil, err
	}

	cgtProxy, err := cgtproxy.New(cgtproxy.WithConfig(cfg),
		cgtproxy.WithLogger(log),
		cgtproxy.WithCGroupMonitor(cGroupMonitor),
		cgtproxy.WithRouteManager(routeManager))

	return cgtProxy, err
}
