package utils

import (
	"fmt"
	"os"
)

type ErrCancelBySignal struct {
	os.Signal
}

type ReloadService struct{}
type CancelByParent struct{}

func (e *ErrCancelBySignal) Error() string {
	return fmt.Sprintf("Cancelled by signal (%v).", e.Signal)
}

func (e *ReloadService) Error() string {
	return fmt.Sprintln("Reload By Parents.")
}

func (e *CancelByParent) Error() string {
	return fmt.Sprintln("Cancel By Parent.")
}
