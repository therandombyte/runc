package system

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// look in /proc to find the process start time so that we can verify
// that this pid has started after ourself
// The most common numeric conversions are Atoi (string to int) and Itoa (int to string).
func GetProcessStartTime(pid int) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(data), " ")
	// the starttime is located at pos 22
	// from the man page
	//
	// starttime %llu (was %lu before Linux 2.6)
	// (22)  The  time the process started after system boot.  In kernels before Linux 2.6, this
	// value was expressed in jiffies.  Since Linux 2.6, the value is expressed in  clock  ticks
	// (divide by sysconf(_SC_CLK_TCK)).
	return parts[22-1], nil // starts at 1
}

/*
/proc is very special in that it is also a virtual filesystem. 
It's sometimes referred to as a process information pseudo-file system. 
It doesn't contain 'real' files but runtime system information (e.g. system memory, devices mounted, hardware configuration, etc).
 For this reason it can be regarded as a control and information centre for the kernel. 
 In fact, quite a lot of system utilities are simply calls to files in this directory. 
 For example, 'lsmod' is the same as 'cat /proc/modules' 
 while 'lspci' is a synonym for 'cat /proc/pci'. 
 By altering files located in this directory you can even read/change kernel parameters (sysctl) while the system is running.

The most distinctive thing about files in this directory is the fact that all of them have a file size of 0, 
with the exception of kcore, mtrr and self.
*/
