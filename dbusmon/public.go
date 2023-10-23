package dbusmon

import (
	"context"

	"github.com/black-desk/cgtproxy/pkg/types"
	. "github.com/black-desk/lib/go/errwrap"
)

func (d *DBusMonitor) Events() <-chan types.CGroupEvents {
	return d.events
}

func (d *DBusMonitor) RunCGroupMonitor(ctx context.Context) (err error) {
	defer Wrap(&err, "running filesystem watcher")
	defer close(d.events)

	d.log.Info("CGroupMonitor is initializing...")
	var events types.CGroupEvents
	events.Events, err = d.walkThroughCgroupFs()

	if err != nil {
		d.log.Warn("failed to initialize CGroupMonitor:", err)
		return err
	}

	d.log.Info("CGroupMonitor has been initialized successfully.")
	err = d.send(ctx, events)
	if err != nil {
		return err
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case eventInfo := <-d.events:
			d.log.Debugw("New DBus call arrived.")
			err = d.send(ctx, eventInfo)
		}
	}

	<-ctx.Done()
	err = context.Cause(ctx)
	return
}
