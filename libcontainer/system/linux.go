//go:build linux
// +build linux

package system

import (
	"os/exec"
	"syscall"
	"unsafe"
)

// 1) Get and Set parent death signal
// 2) Get and Set capabilities of calling process
// 3) Set the controlling terminal of calling process
// 4) Allow execution of a cmd/executable using execve
type ParentDeathSignal int

func (p ParentDeathSignal) Restore() error {
	if p == 0 {
		return nil
	}
	current, err := GetParentDeathSignal()
	if err != nil {
		return err
	}
	if p == current {
		return nil
	}
	return p.Set()
}

func (p ParentDeathSignal) Set() error {
	return SetParentDeathSignal(uintptr(p))
}

// looks for the command/executable in the path
// if found, call the execve syscall that will run the executable
func Execv(cmd string, args []string, env []string) error {
	name, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	return syscall.Exec(name, args, env)
}

// This is the signal that the calling process will get when its parent dies
// Ex: The child process dies because it receives PR_SET_PDEATHSIG signal when the parent thread dies
// PR_SET_PDEATHSIG is preserved across exec(2) but gets cleared after fork(2)
func SetParentDeathSignal(sig uintptr) error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_PDEATHSIG, sig, 0); err != 0 {
		return err
	}
	return nil
}

// Return the current value of the parent process death signal, in the location pointed to by (int *) arg2
func GetParentDeathSignal() (ParentDeathSignal, error) {
	var sig int
	_, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_GET_PDEATHSIG, uintptr(unsafe.Pointer(&sig)), 0)
	if err != 0 {
		return -1, err
	}
	return ParentDeathSignal(sig), nil
}

// keep the capabilities of calling process
// execve checks this??
func SetKeepCaps() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_KEEPCAPS, 1, 0); err != 0 {
		return err
	}

	return nil
}

// clears the capabilities of calling process, strip em naked?? no syscalls allowed??
// prctl() manipulates various aspects of the behavior of the calling thread or process
// Ex: PR_SET_KEEPCAPS: Sets the state of the calling thread's "keep capabilities" flag
//     Linux has capabilities list: https://man7.org/linux/man-pages/man7/capabilities.7.html
func ClearKeepCaps() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_PRCTL, syscall.PR_SET_KEEPCAPS, 0, 0); err != 0 {
		return err
	}

	return nil
}

// setctty: set controlling terminal
// IOCTL is a syscall to the device driver
//       ioctl() manipulates the underlying device parameters of special files
//       Many operating characteristics of character special files (e.g. terminals) may be controlled with ioctl() requests
// TIOCSCTTY is a syscall to make the given terminal as controlling terminal
// of calling process
// uintptr is an integer representation of a memory address in golang
func Setctty() error {
	if _, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, 0, uintptr(syscall.TIOCSCTTY), 0); err != 0 {
		return err
	}
	return nil
}

/*
TTYs are devices that enable typing (type, typewriter) from a distance (tele).
In a modern operating system (OS), the concept applies directly. 
Linux uses a device file to represent a virtual TTY, which enables interaction 
with the OS by handling input (usually a keyboard) and output (usually a screen)
*/
