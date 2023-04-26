//go:build linux
// +build linux

package apparmor

import (
	"io"
	"os"
	"text/template"
)

// 1) Generate a base apparmor template and send it out to  any writerstream
// 2) Add tunable variables and abstractions into the profile if required

type data struct {
	Name         string
	Imports      []string
	InnerImports []string
}

const baseTemplate = `
{{range $value := .Imports}}
{{$value}}
{{end}}

profile {{.Name}} flags=(attach_disconnected,mediate_deleted) {
{{range $value := .InnerImports}}
  {{$value}}
{{end}}

  network,
  capability,
  file,
  umount,

  deny @{PROC}/sys/fs/** wklx,
  deny @{PROC}/sysrq-trigger rwklx,
  deny @{PROC}/mem rwklx,
  deny @{PROC}/kmem rwklx,
  deny @{PROC}/sys/kernel/[^s][^h][^m]* wklx,
  deny @{PROC}/sys/kernel/*/** wklx,

  deny mount,

  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/efi/efivars/** rwklx,
  deny /sys/kernel/security/** rwklx,
}
`
// Package template implements data-driven templates for generating textual output.
// Templates are executed by applying them to a data structure.
// io.Writer: you can write directly into files or other streams, such as http connections
//           Advantage: if you are passing a Writer interface, you can pass anything which implements Write
//           that is not only a file but a http.ResponseWriter, for example, or stdout os.Stdout, without changing the struct methods.
func generateProfile(out io.Writer) error {
	compiled, err := template.New("apparmor_profile").Parse(baseTemplate)
	if err != nil {
		return err
	}
	data := &data{
		Name: "docker-default",
	}
	if tunablesExists() {
		data.Imports = append(data.Imports, "#include <tunables/global>")
	} else {
		data.Imports = append(data.Imports, "@{PROC}=/proc/")
	}
	if abstractionsExists() {
		data.InnerImports = append(data.InnerImports, "#include <abstractions/base>")
	}
	if err := compiled.Execute(out, data); err != nil {
		return err
	}
	return nil
}

// check if the tunables/global exist
/*
tunables: those provide variables that can be used in templates, 
		for example if you want a custom dir considered as it would be a home directory 
		you could modify /etc/apparmor.d/tunables/home which defines the base path rules use for home directories
		
		The tunables directory (/etc/apparmor.d/tunables) contains global variable definitions. 
		When used in a profile, these variables expand to a value that can be changed without 
		changing the entire profile. Add all the tunables definitions that should be available to 
		every profile to /etc/apparmor.d/tunables/global.
*/
func tunablesExists() bool {
	_, err := os.Stat("/etc/apparmor.d/tunables/global")
	return err == nil
}

// check if abstractions/base exist
/*
Abstractions are includes that are grouped by common application tasks. 
These tasks include access to authentication mechanisms, access to name service routines, 
common graphics requirements, and system accounting. Files listed in these abstractions 
are specific to the named task. Programs that require one of these files usually also 
require other files listed in the abstraction file (depending on the local configuration 
	and the specific requirements of the program). Find abstractions in /etc/apparmor.d/abstractions.
*/
func abstractionsExists() bool {
	_, err := os.Stat("/etc/apparmor.d/abstractions/base")
	return err == nil
}
