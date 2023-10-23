package proxyService

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ComixHe/proxy/config"
	"github.com/black-desk/cgtproxy/pkg/types"
	"github.com/godbus/dbus/v5"
)

type AMProxy1 interface {
	SetProxy(cgroupPath string, state types.CgroupEventType) (e *dbus.Error)
}

type DCCProxy1 interface {
	Reload(config_json string) (e *dbus.Error)
}

func (p *proxyService) SetProxy(cgroupPath string, state types.CgroupEventType) (e *dbus.Error) {

	if !p.mutex.TryLock() {
		e = dbus.MakeFailedError(errors.New("service is busy, please try again later."))
		return
	}

	defer p.mutex.Unlock()
	busName := <-p.busNameCh

	UID, e := p.GetUID(busName)
	if e != nil {
		return
	}
	p.log.Debugw("call SetProxy", "sender", busName, "user", UID)
	userCore := p.cores[UID]

	var eventInfo types.CGroupEvent
	eventInfo.EventType = state
	eventInfo.Path = cgroupPath

	ch := make(chan error)

	events := types.CGroupEvents{
		Events: []types.CGroupEvent{
			{EventType: state, Path: cgroupPath},
		},
		Result: ch,
	}
	p.log.Debugw("info:", "state", state, "path", cgroupPath)
	userCore.EventIN() <- events
	err := <-ch
	if err != nil {
		e = dbus.MakeFailedError(err)
	}

	if p.log == nil {
		return
	}

	p.log.Debugw("call SetProxy:", "path", cgroupPath, "state", state)

	return
}

func (p *proxyService) Reload(config_json string) (e *dbus.Error) {

	if !p.mutex.TryLock() {
		e = dbus.MakeFailedError(errors.New("service is busy, please try again later."))
		return
	}
	defer p.mutex.Unlock()

	userCfg := &config.UserConfig{}
	err := json.Unmarshal([]byte(config_json), userCfg)
	if err != nil {
		e = dbus.MakeFailedError(err)
		return
	}

	busName := <-p.busNameCh

	UID, e := p.GetUID(busName)
	if e != nil {
		return
	}
	p.log.Debugw("call Reload", "sender", busName, "user", UID)

	userCore := p.cores[UID]
	if userCore != nil {
		userCore.StopCoreService()
		err := <-userCore.ErrOut()

		if err != nil {
			e = dbus.MakeFailedError(err)
			p.log.Warn("stop failed,")
			return
		}

		delete(p.cores, UID)
	}

	subCtx, subCancel := context.WithCancelCause(p.ctx)

	//coreservice.New(coreservice.WithCancelCause(subCtx, subCancel),coreservice.WithCoreConfig(cfg *config.Config))

	return
}
