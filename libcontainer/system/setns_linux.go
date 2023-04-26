package system

import (
	"fmt"
	"runtime"
	"syscall"
)

// Via http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/commit/?id=7b21fddd087678a70ad64afc0f632e0f1071b092
//
// We need different setns values for the different platforms and arch
// We are declaring the macro here because the SETNS syscall does not exist in th stdlib

//  The setns() system call allows the calling thread to move into different namespaces
// lsof - list of open files. This command will show the open files and its file descriptors which are integers
// But all file descriptors are CHR – Character special file (or character device file) and lives purely in memory.
// These are not regular files
var setNsMap = map[string]uintptr{
	"linux/386":     346,
	"linux/arm64":   268,
	"linux/amd64":   308,
	"linux/arm":     375,
	"linux/ppc":     350,
	"linux/ppc64":   350,
	"linux/ppc64le": 350,
	"linux/s390x":   339,
}

//Package "runtime" contains operations that interact with Go's runtime system, such as functions to control goroutines
// GOOS is the running program's operating system target: one of darwin, freebsd, linux, and so on. 
// To view possible combinations of GOOS and GOARCH, run "go tool dist list".
// GOARCH is the running program's architecture target: one of 386, amd64, arm, s390x, and so on.
/* Sprint verbs:
%v	the value in a default format
	when printing structs, the plus flag (%+v) adds field names
%#v	a Go-syntax representation of the value
%T	a Go-syntax representation of the type of the value
%%	a literal percent sign; consumes no value
%s	the uninterpreted bytes of the string or slice
%q	a double-quoted string safely escaped with Go syntax
*/
// getting the namespace of current OS and adding into slice
var sysSetns = setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]

func SysSetns() uint32 {
	return uint32(sysSetns)
}

func Setns(fd uintptr, flags uintptr) error {
	ns, exists := setNsMap[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)]
	if !exists {
		return fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	_, _, err := syscall.RawSyscall(ns, fd, flags, 0)
	if err != 0 {
		return err
	}
	return nil
}
