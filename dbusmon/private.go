package dbusmon

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/black-desk/cgtproxy/pkg/types"
)

const cgRoot = "/sys/fs/cgroup"

func (d *DBusMonitor) walkFn(events *[]types.CGroupEvent) func(path string, dir fs.DirEntry, err error) error {
	return func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				d.log.Debug(
					"Cgroup had been removed.",
					"path", path,
				)
				err = nil
			}
			d.log.Errorw(
				"Errors occurred while first time going through cgroupfs.",
				"path", path,
				"error", err,
			)
			err = nil
		}

		if !dir.IsDir() {
			return nil
		}

		path = strings.TrimRight(path, "/")

		if path == cgRoot {
			return nil
		}

		*events = append(*events, types.CGroupEvent{
			Path:      path,
			EventType: types.CgroupEventTypeNew,
		})

		return nil
	}
}

func (d *DBusMonitor) walkThroughCgroupFs() (ret []types.CGroupEvent, err error) {
	events := []types.CGroupEvent{}

	err = filepath.WalkDir(cgRoot, d.walkFn(&events))
	if err != nil {
		return nil, err
	}

	ret = events
	return
}

func (m *DBusMonitor) send(ctx context.Context, cgEvents types.CGroupEvents) (err error) {
	for i := range cgEvents.Events {
		path := strings.TrimRight(cgEvents.Events[i].Path, "/")
		cgEvents.Events[i].Path = path
	}

	cnt := len(cgEvents.Events)

	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case m.events <- cgEvents:
		m.log.Debugw("Cgroup events sent.",
			"size", cnt,
		)
	}

	return
}
