package dbusmon

import (
	"github.com/ComixHe/proxy/utils"
	"github.com/black-desk/cgtproxy/pkg/types"
	. "github.com/black-desk/lib/go/errwrap"
	"go.uber.org/zap"
)

type DBusMonitor struct {
	events chan types.CGroupEvents
	log    *zap.SugaredLogger
}

type Opt func(d *DBusMonitor) (ret *DBusMonitor, err error)

func New(opts ...Opt) (ret *DBusMonitor, err error) {
	defer Wrap(&err, "create DBusMonitor.")
	d := &DBusMonitor{}

	for i := range opts {
		d, err = opts[i](d)
		if err != nil {
			return nil, err
		}
	}

	if d.log == nil {
		d.log = utils.GetLogger()
	}

	if d.events == nil {
		panic("Uninitialized channel.")
	}

	ret = d
	d.log.Debugw("Create DBusMonitor.")
	return
}

func WithLogger(log *zap.SugaredLogger) Opt {
	return func(d *DBusMonitor) (ret *DBusMonitor, err error) {
		d.log = log
		ret = d
		return
	}
}

func WithChannel(eventsChan chan types.CGroupEvents) Opt {
	return func(d *DBusMonitor) (ret *DBusMonitor, err error) {
		d.events = eventsChan
		ret = d
		return
	}
}
