// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/bborbe/docker-utils"
)

var (
	registryPtr     = flag.String("registry", "", "Registry")
	usernamePtr     = flag.String("username", "", "Username")
	passwordPtr     = flag.String("password", "", "Password")
	passwordFilePtr = flag.String("passwordfile", "", "Password-File")
	repositoryPtr   = flag.String("repository", "", "Repository")
	tagPtr          = flag.String("tag", "", "Tag")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := do(context.Background()); err != nil {
		glog.Exitf("%+v", err)
	}
}

func do(ctx context.Context) error {
	registry := &docker.Registry{
		Url:      *registryPtr,
		Username: *usernamePtr,
		Password: *passwordPtr,
	}
	if len(*passwordFilePtr) > 0 {
		if err := registry.RegistryPasswordFromFile(*passwordFilePtr); err != nil {
			return err
		}
	}
	glog.V(2).Infof("use registry %v, repo %v and tag %v", registry, docker.RepositoryName(*repositoryPtr), docker.TagName(*tagPtr))
	client := docker.NewV2Client(docker.NewHttpClient(http.DefaultClient), *registry)
	if err := client.DeleteTag(ctx, docker.RepositoryName(*repositoryPtr), docker.TagName(*tagPtr)); err != nil {
		return errors.Wrap(err, "delete tag exists failed")
	}
	fmt.Printf("tag deleted\n")
	return nil
}
