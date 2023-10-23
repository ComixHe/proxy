package coreservice

import (
	"errors"

	"github.com/ComixHe/proxy/utils"
	"github.com/black-desk/cgtproxy/pkg/types"
)

func (service *CoreService) RunCoreService() {
	defer close(service.errCh)
	defer service.log.Debug("RunCoreService exit")

	err := service.core.RunCGTProxy(service.ctx)

	var reloadErr *utils.ReloadService
	var cancelErr *utils.CancelByParent
	var sigErr *utils.ErrCancelBySignal

	if errors.As(err, &reloadErr) || errors.As(err, &cancelErr) || errors.As(err, &sigErr) {
		service.log.Debug("Not critical error,ignore:", err)
		err = nil
	}

	service.errCh <- err
}

func (service *CoreService) EventIN() chan<- types.CGroupEvents {
	return service.amCh
}

func (service *CoreService) ErrOut() <-chan error {
	return service.errCh
}

func (service *CoreService) StopCoreService() {
	service.cancel(&utils.CancelByParent{})
}
