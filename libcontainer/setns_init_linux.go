//go:build linux
// +build linux

package libcontainer

import (
	"os"

	"github.com/opencontainers/runc/libcontainer/apparmor"
	"github.com/opencontainers/runc/libcontainer/label"
	"github.com/opencontainers/runc/libcontainer/system"
)

// linuxSetnsInit performs the container's initialization for running a new process
// inside an existing container.
type linuxSetnsInit struct {
	config *initConfig
}

// invoking functions from init_linux.go
// syscall package has a function "Setrlimit" to set rlimits
// needs two integers, soft and hard limit value
// Rlimit is a struct within the syscall package, can be set using &syscall.Rlimit{Max: rlimit.Hard, Cur: rlimit.Soft}
// who is setting the hard/soft values, its coming from config
func (l *linuxSetnsInit) Init() error {
	if err := setupRlimits(l.config.Config); err != nil { // sets hard and soft limit
		return err
	}
	if err := finalizeNamespace(l.config); err != nil { // whitelist current calling process, set user, drop capabilities, chdir
		return err
	}
	if err := apparmor.ApplyProfile(l.config.Config.AppArmorProfile); err != nil {
		return err
	}
	if l.config.Config.ProcessLabel != "" {
		if err := label.SetProcessLabel(l.config.Config.ProcessLabel); err != nil {
			return err
		}
	}
	return system.Execv(l.config.Args[0], l.config.Args[0:], os.Environ())  // is this the command passed in cli while launching the container??
}
