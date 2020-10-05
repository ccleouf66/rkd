module rkd

go 1.13

require (
	github.com/containers/image/v5 v5.6.0
	github.com/gofrs/flock v0.7.1
	github.com/google/go-github/v32 v32.1.0
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli v1.22.4
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200922070232-aee5d888a860 // indirect
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.3.0
	rsc.io/letsencrypt v0.0.3 // indirect
)
replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
replace github.com/containerd/containerd v1.3.4 => github.com/containerd/containerd v1.4.1
