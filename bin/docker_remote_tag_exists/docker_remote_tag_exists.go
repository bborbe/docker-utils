package main

import (
	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

const (
	PARAMETER_REGISTRY = "registry"
)

var (
	registryPtr = flag.String(PARAMETER_REGISTRY, "", "Docker Registry")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	glog.V(2).Infof("exists image in registry %v", *registryPtr)
	return nil
}
