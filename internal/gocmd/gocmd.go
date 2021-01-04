/*
Package gocmd is an abstraction layer over the currently installed Go system binary.
Keep in mind, the binary `shore` was compiled with Go, but after it has been compiled into a binary file, Golang is no longer in the picture!
*/
package gocmd

import (
	"os"
	"os/exec"
)

// GoWrapper - An interface for communicating with the currently installed Go binary.
type GoWrapper interface {
	Init(name string) (string, error)
	Get(packages []string) (string, error)
	Vendor() (string, error)
	Version() (string, error)
}

type GoCmd struct {
	Dir string
	Env []string
}

func NewGoCmd(dir string) *GoCmd {
	return &GoCmd{
		Dir: dir,
		Env: append(os.Environ(), "GOPROXY=direct", "GOPRIVATE=*"),
	}
}

func (g *GoCmd) Init(name string) (string, error) {
	goCmd := exec.Command("go", "mod", "init", name)
	return g.call(goCmd)
}

func (g *GoCmd) Get(packages []string) (string, error) {
	// Translates to `go get package1 package2 package3`
	goCmd := exec.Command("go", append([]string{"get"}, packages...)...)
	return g.call(goCmd)
}

func (g *GoCmd) Vendor() (string, error) {
	goCmd := exec.Command("go", "mod", "vendor")
	return g.call(goCmd)
}

func (g *GoCmd) Version() (string, error) {
	goCmd := exec.Command("go", "version")
	return g.call(goCmd)
}

func (g *GoCmd) call(goCmd *exec.Cmd) (string, error) {
	goCmd.Env = g.Env
	goCmd.Dir = g.Dir

	stdout, stderr := goCmd.CombinedOutput()
	return string(stdout), stderr
}
