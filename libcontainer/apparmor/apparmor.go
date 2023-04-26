//go:build apparmor && linux
// +build apparmor,linux

package apparmor

// #cgo LDFLAGS: -lapparmor
// #include <sys/apparmor.h>
// #include <stdlib.h>
import "C"
import (
	"io/ioutil"
	"os"
	"unsafe"
)

// 1) Check if apparmor is enabled
// 2) Apply apparmor profile

func IsEnabled() bool {
	if _, err := os.Stat("/sys/kernel/security/apparmor"); err == nil && os.Getenv("container") == "" {
		if _, err = os.Stat("/sbin/apparmor_parser"); err == nil {
			buf, err := ioutil.ReadFile("/sys/module/apparmor/parameters/enabled")
			return err == nil && len(buf) > 1 && buf[0] == 'Y'
		}
	}
	return false
}

func ApplyProfile(name string) error {
	if name == "" {
		return nil
	}
	// Go string to C string
	// The C string is allocated in the C heap using malloc.
	// It is the caller's responsibility to arrange for it to be
	// freed, such as by calling C.free (be sure to include stdlib.h
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	/*
	The aa_change_onexec() function is like the aa_change_profile() function except it
       specifies that the profile transition should take place on the next exec instead of
       immediately.
	*/
	if _, err := C.aa_change_onexec(cName); err != nil {
		return err
	}
	return nil
}

/*
An AppArmor profile applies to an executable program; 
if a portion of the program needs different access permissions 
than other portions, the program can "change profile" to a
different profile
*/
