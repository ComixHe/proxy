package proxyService

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

const userSlicePath = "/sys/fs/cgroup/user.slice"

func (export *proxyService) exportMethod() (err error) {
	var amObj AMProxy1 = export
	err = export.conn.Export(amObj, ObjectPath, AMInterface)
	if err != nil {
		return
	}

	var dccObj DCCProxy1 = export
	err = export.conn.Export(dccObj, ObjectPath, DCCInterface)
	if err != nil {
		return
	}

	return
}

func (p *proxyService) GetUID(busName string) (uid int, e *dbus.Error) {
	var UID int
	obj := p.conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	err := obj.Call("org.freedesktop.DBus.GetConnectionUnixUser", 0, busName).Store(&UID)
	if err != nil {
		p.log.Warn(err)
		e = dbus.MakeFailedError(errors.New("Identify Failed."))
		return
	}
	return UID, nil
}
