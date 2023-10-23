package coreservice

import (
	"context"
	"errors"

	"github.com/ComixHe/proxy/utils"
	"github.com/black-desk/cgtproxy/pkg/cgtproxy"
	"github.com/black-desk/cgtproxy/pkg/cgtproxy/config"
	"github.com/black-desk/cgtproxy/pkg/types"
	"go.uber.org/zap"
)

type CoreService struct {
	log    *zap.SugaredLogger
	core   *cgtproxy.CGTProxy
	config *config.Config
	ctx    context.Context
	cancel context.CancelCauseFunc
	amCh   chan types.CGroupEvents
	errCh  chan error
}

type Opt func(p *CoreService) (ret *CoreService, err error)

func New(opts ...Opt) (ret *CoreService, err error) {

	p := &CoreService{}
	p.amCh = make(chan types.CGroupEvents, 1)
	p.errCh = make(chan error, 1)

	for i := range opts {
		p, err = opts[i](p)
		if err != nil {
			return
		}
	}

	if p.log == nil {
		p.log = utils.GetLogger()
	}

	cgtProxy, err := getCGtproxy(p.config, p.amCh, p.log)
	if err != nil {
		return nil, err
	}

	p.core = cgtProxy
	ret = p

	p.log.Debugw("Create DBus service...")

	return
}

func WithCoreConfig(cfg *config.Config) Opt {
	return func(p *CoreService) (ret *CoreService, err error) {
		if cfg == nil {
			err = errors.New("core config is nil.")
			return
		}
		p.config = cfg
		ret = p
		return
	}
}

func WithCancelCause(ctx context.Context, cancel context.CancelCauseFunc) Opt {
	return func(p *CoreService) (ret *CoreService, err error) {
		p.ctx = ctx
		p.cancel = cancel
		ret = p
		return
	}
}
