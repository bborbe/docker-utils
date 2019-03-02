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

	"github.com/pkg/errors"

	"github.com/bborbe/docker-utils"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

const defaultMaxAge = 100 * 24 * time.Hour

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := &application{
		registry: docker.Registry{
			Url: "https://registry-1.docker.io",
		},
	}
	flag.StringVar(&app.registry.Username, "username", "", "Registry Username")
	flag.StringVar(&app.registry.Password, "password", "", "Registry Password")
	flag.StringVar(&app.passwordFile, "passwordfile", "", "Password-File")
	flag.DurationVar(&app.maxAge, "max-age", defaultMaxAge, "Max age")

	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	app.printArgs()

	if err := app.validate(); err != nil {
		glog.Exit(err)
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
	registry     docker.Registry
	passwordFile string
	maxAge       time.Duration
}

func (a *application) printArgs() {
	glog.V(0).Infof("Parameter registry.url: %s", a.registry.Url)
	glog.V(0).Infof("Parameter registry.username: %s", a.registry.Username)
	glog.V(0).Infof("Parameter registry.password length: %d", len(a.registry.Password))
	glog.V(0).Infof("Parameter passwordFile: %s", a.passwordFile)
	glog.V(0).Infof("Parameter maxAge: %v", a.maxAge)
}

func (a *application) validate() error {
	return nil
}

func (a *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(a.passwordFile) > 0 {
		if err := a.registry.RegistryPasswordFromFile(a.passwordFile); err != nil {
			return err
		}
	}
	now := time.Now()

	httpClient := docker.NewHttpClient(http.DefaultClient)
	dockerHubClient := docker.NewDockerHubClient(httpClient, a.registry)
	repositories := make(chan docker.DockerHubTagRepository, runtime.NumCPU())
	go func() {
		defer close(repositories)
		if err := dockerHubClient.ListRepositories(ctx, docker.RepositoryName(a.registry.Username), repositories); err != nil {
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
					if now.Sub(date) > a.maxAge {
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
