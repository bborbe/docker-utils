module github.com/bborbe/docker-utils

go 1.26.0

require (
	github.com/actgardner/gogen-avro/v9 v9.2.0
	github.com/bborbe/argument v1.3.2
	github.com/bborbe/flagenv v0.0.0-20181019084341-2956c4545608
	github.com/bborbe/io v0.0.0-20180829202151-54b762caaee8
	github.com/golang/glog v1.2.5
	github.com/google/addlicense v1.2.0
	github.com/incu6us/goimports-reviser/v3 v3.12.3
	github.com/kisielk/errcheck v1.9.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.12.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.39.1
	github.com/pkg/errors v0.9.1
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067
	golang.org/x/vuln v1.1.4
)

require (
	github.com/bborbe/assert v0.0.0-20181116222016-22a6c6341415 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/incu6us/goimports-reviser v0.1.6 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.2 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/telemetry v0.0.0-20260109210033-bd525da824e2 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

exclude cloud.google.com/go v0.26.0

exclude (
	github.com/go-logr/glogr v1.0.0-rc1
	github.com/go-logr/glogr v1.0.0
)

exclude (
	github.com/go-logr/logr v1.0.0-rc1
	github.com/go-logr/logr v1.0.0
)

exclude (
	go.yaml.in/yaml/v3 v3.0.3
	go.yaml.in/yaml/v3 v3.0.4
)

exclude (
	golang.org/x/tools v0.38.0
	golang.org/x/tools v0.39.0
)

exclude (
	k8s.io/api v0.34.0
	k8s.io/api v0.34.1
	k8s.io/api v0.34.2
	k8s.io/api v0.34.3
	k8s.io/api v0.34.4
	k8s.io/api v0.35.0
	k8s.io/api v0.35.1
)

exclude (
	k8s.io/apiextensions-apiserver v0.34.0
	k8s.io/apiextensions-apiserver v0.34.1
	k8s.io/apiextensions-apiserver v0.34.2
	k8s.io/apiextensions-apiserver v0.34.3
	k8s.io/apiextensions-apiserver v0.34.4
	k8s.io/apiextensions-apiserver v0.35.0
	k8s.io/apiextensions-apiserver v0.35.1
)

exclude (
	k8s.io/apimachinery v0.34.0
	k8s.io/apimachinery v0.34.1
	k8s.io/apimachinery v0.34.2
	k8s.io/apimachinery v0.34.3
	k8s.io/apimachinery v0.34.4
	k8s.io/apimachinery v0.35.0
	k8s.io/apimachinery v0.35.1
)

exclude (
	k8s.io/client-go v0.34.0
	k8s.io/client-go v0.34.1
	k8s.io/client-go v0.34.2
	k8s.io/client-go v0.34.3
	k8s.io/client-go v0.34.4
	k8s.io/client-go v0.35.0
	k8s.io/client-go v0.35.1
)

exclude (
	k8s.io/code-generator v0.34.0
	k8s.io/code-generator v0.34.1
	k8s.io/code-generator v0.34.2
	k8s.io/code-generator v0.34.3
	k8s.io/code-generator v0.34.4
	k8s.io/code-generator v0.35.0
	k8s.io/code-generator v0.35.1
)

exclude (
	sigs.k8s.io/structured-merge-diff/v6 v6.0.0
	sigs.k8s.io/structured-merge-diff/v6 v6.1.0
	sigs.k8s.io/structured-merge-diff/v6 v6.2.0
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0
)

replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20250701173324-9bd5c66d9911
