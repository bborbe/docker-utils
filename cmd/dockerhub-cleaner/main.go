package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/bborbe/argument"

	"github.com/pkg/errors"

	"github.com/bborbe/docker-utils"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")

	app := &application{}
	if err := argument.Parse(app); err != nil {
		glog.Exitf("parse args failed: %v", err)
	}

	glog.V(0).Infof("application started")
	if err := app.run(contextWithSig(context.Background())); err != nil {
		glog.Exitf("application failed: %+v", err)
	}
	glog.V(0).Infof("application finished")
	os.Exit(0)
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
		glog.V(2).Infof("cancled")
	}()

	return ctxWithCancel
}

type application struct {
	Url          string        `required:"true" arg:"url" default:"https://registry-1.docker.io" usage:"Registry Url"`
	Username     string        `required:"true" arg:"username" usage:"Registry Username"`
	Password     string        `arg:"password" usage:"Registry Password" display:"length"`
	PasswordFile string        `arg:"passwordfile" usage:"Password-File"`
	MaxAge       time.Duration `required:"true" arg:"max-age" usage:"Max age" default:"2400h"`
}

func (a *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	registry := docker.Registry{
		Url:      a.Url,
		Username: a.Username,
		Password: a.Password,
	}

	if len(a.PasswordFile) > 0 {
		if err := registry.RegistryPasswordFromFile(a.PasswordFile); err != nil {
			return err
		}
	}
	now := time.Now()

	httpClient := docker.NewHttpClient(http.DefaultClient)
	dockerHubClient := docker.NewDockerHubClient(httpClient, registry)
	repositories := make(chan docker.DockerHubTagRepository, runtime.NumCPU())
	go func() {
		defer close(repositories)
		if err := dockerHubClient.ListRepositories(ctx, docker.RepositoryName(registry.Username), repositories); err != nil {
			glog.Warningf("read repositories failed: %v", err)
			cancel()
		}
	}()
	type entry struct {
		Tag  docker.TagName
		Repo docker.RepositoryName
	}
	var wg sync.WaitGroup
	for repository := range repositories {
		select {
		case <-ctx.Done():
			return nil
		default:
			var list []entry
			tags := make(chan docker.DockerHubTag, runtime.NumCPU())
			go func() {
				defer close(tags)
				if err := dockerHubClient.ListTags(ctx, repository.RepositoryName(), tags); err != nil {
					glog.Warningf("list tags failed: %v", err)
					cancel()
				}
			}()
			for tag := range tags {
				select {
				case <-ctx.Done():
					return nil
				default:
					date, err := time.Parse(time.RFC3339Nano, tag.LastUpdated)
					if err != nil {
						return errors.Wrapf(err, "parse date %s failed", tag.LastUpdated)
					}
					if now.Sub(date) > a.MaxAge {
						list = append(list, entry{
							Tag:  tag.Tag,
							Repo: repository.RepositoryName(),
						})
					}
				}
			}
			for _, rm := range list {
				wg.Add(1)
				go func(rm entry) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
						if err := dockerHubClient.DeleteTag(ctx, rm.Repo, rm.Tag); err != nil {
							cancel()
							glog.Warningf("delete failed: %v", err)
						}
						fmt.Printf("deleted %s:%s\n", rm.Repo, rm.Tag)
					}
				}(rm)
			}
		}
	}
	wg.Wait()
	return nil
}
