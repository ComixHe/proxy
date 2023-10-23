package proxyService

import (
	"context"
	"errors"
	"sync"

	"github.com/ComixHe/proxy/coreservice"
	"github.com/ComixHe/proxy/utils"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"go.uber.org/zap"
)

const (
	intro = `
<node>
	<interface name="org.desktopspec.deepin.AMProxy1">
		<method name="SetProxy">
			<arg name="cgroup_path" direction="in" type="s"/>
			<arg name="state" direction="in" type="y"/>
		</method>
	</interface>
	<interface name="org.desktopspec.deepin.DCCProxy1">
		<method name="Reload">
			<arg name="config" direction="in" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`

	ObjectPath   = "/org/desktopspec/deepin/Proxy1"
	AMInterface  = "org.desktopspec.deepin.AMProxy1"
	DCCInterface = "org.desktopspec.deepin.DCCProxy1"
)

type proxyService struct {
	conn      *dbus.Conn
	log       *zap.SugaredLogger
	mutex     sync.Mutex
	ctx       context.Context
	cores     map[int]*coreservice.CoreService
	busNameCh <-chan string
}

type Opt func(p *proxyService) (ret *proxyService, err error)

func New(opts ...Opt) (ret *proxyService, err error) {

	p := &proxyService{}
	p.mutex = sync.Mutex{}

	for i := range opts {
		p, err = opts[i](p)
		if err != nil {
			return
		}
	}

	if p.log == nil {
		p.log = utils.GetLogger()
	}

	if p.conn == nil || !p.conn.Connected() {
		err = errors.New("couldn't connect to DBus.")
		return
	}

	ret = p

	p.log.Debugw("Create DBus service...")

	return
}

func (service *proxyService) RunDBusService(ctx context.Context) (err error) {
	defer service.log.Debug("stop DBusService")

	err = service.exportMethod()
	if err != nil {
		return err
	}

	err = service.conn.Export(introspect.Introspectable(intro), ObjectPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return
	}

	reply, err := service.conn.RequestName("org.desktopspec.deepin.Proxy1", dbus.NameFlagDoNotQueue)
	if err != nil {
		return
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		err = errors.New("name already taken.")
		return
	}

	service.log.Debugw("start DBusService...")

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		}
	}

	<-ctx.Done()

	for k, v := range service.cores {
		service.log.Info("stop proxy for user ", k)
		v.StopCoreService()
	}

	return
}

func WithBusNameChannel(ch <-chan string) Opt {
	return func(p *proxyService) (ret *proxyService, err error) {
		p.busNameCh = ch
		ret = p
		return
	}
}

func WithContext(context context.Context) Opt {
	return func(p *proxyService) (ret *proxyService, err error) {
		p.ctx = context
		ret = p
		return
	}
}

func WithLogger(log *zap.SugaredLogger) Opt {
	return func(p *proxyService) (ret *proxyService, err error) {
		p.log = log
		ret = p
		return
	}
}

func WithDBusConnection(conn *dbus.Conn) Opt {
	return func(p *proxyService) (ret *proxyService, err error) {
		if conn == nil || !conn.Connected() {
			panic("DBus is disconnected.")
		}
		p.conn = conn
		ret = p
		return
	}
}
