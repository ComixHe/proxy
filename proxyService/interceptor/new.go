package interceptor

import (
	"github.com/ComixHe/proxy/utils"
	"go.uber.org/zap"
)

type IncomingInter struct {
	busNameCh chan string
	log       *zap.SugaredLogger
}

type Opt func(p *IncomingInter) (ret *IncomingInter, err error)

func New(opts ...Opt) (ret *IncomingInter, err error) {
	i := &IncomingInter{}
	i.busNameCh = make(chan string, 1)

	for index := range opts {
		i, err = opts[index](i)
		if err != nil {
			return
		}
	}

	if i.log == nil {
		i.log = utils.GetLogger()
	}

	ret = i
	return
}

func WithLogger(log *zap.SugaredLogger) Opt {
	return func(p *IncomingInter) (ret *IncomingInter, err error) {
		p.log = log
		ret = p
		return
	}
}
