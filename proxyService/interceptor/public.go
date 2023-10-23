package interceptor

import (
	"strconv"

	"github.com/ComixHe/proxy/proxyService"
	"github.com/godbus/dbus/v5"
)

func (i *IncomingInter) GetBusName(msg *dbus.Message) {
	if msg.Type != dbus.TypeMethodCall {
		return
	}

	raw := msg.Headers[dbus.FieldInterface].String()
	inter, _ := strconv.Unquote(raw)

	if inter == proxyService.AMInterface || inter == proxyService.DCCInterface {
		busName, _ := strconv.Unquote(msg.Headers[dbus.FieldSender].String())
		i.busNameCh <- busName
	}
}

func (i *IncomingInter) BusNameChan() <-chan string {
	return i.busNameCh
}
