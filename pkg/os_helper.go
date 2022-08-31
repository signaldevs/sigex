package sigex

import (
	"os/exec"
	"syscall"
)

var osHelperInstance OSHelper

type OSHelper interface {
	LookPath(string) (string, error)
	Exec(string, []string, []string) error
}

type osHelper struct{}

func (o osHelper) LookPath(path string) (string, error) {
	return exec.LookPath(path)
}

func (o osHelper) Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

func GetOSHelper() OSHelper {
	return osHelperInstance
}

func SetOSHelper(helper OSHelper) {
	osHelperInstance = helper
}

func init() {
	SetOSHelper(osHelper{})
}
