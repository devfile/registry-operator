# Operator SDK v1.28 to v1.39.1 migration

This project uses Operator SDK v1.39.1 with the `go/v4` plugin. Go remains
1.26.0 and the UBI9 Go 1.26.7 builder is retained.

The [SDK upgrade guides](https://sdk.operatorframework.io/docs/upgrading-sdk-version/)
and the linked [go/v3 to go/v4 guide](https://book.kubebuilder.io/migration/manually_migration_guide_gov3_to_gov4)
were applied cumulatively:

| Release | Applicable changes |
| --- | --- |
| 1.29, 1.32 | No migrations. |
| 1.30 | No concrete Go migration in the published guide. |
| 1.31 | Ansible changes do not apply to this Go operator. |
| 1.33 | go/v4 PROJECT metadata, `cmd/main.go`, `internal/controller`, build/import/test paths, admission/v1 tests, and modern Kustomize configuration. |
| 1.34 | Build multi-architecture images directly from Dockerfile, using target architecture arguments. |
| 1.35 | Helm-only changes do not apply. |
| 1.36 | Versioned tool installation and `build-installer`; dependency and envtest versions are superseded by 1.39. |
| 1.37 | Quarkus and hybrid-Helm changes do not apply. |
| 1.38 | Kubernetes libraries 0.30.1, controller-runtime 0.18.4, Ginkgo 2.17.1, Gomega 1.32.0, Kubernetes envtest 1.30.0, Kustomize 5.4.2, and authenticated manager metrics replacing kube-rbac-proxy. |
| 1.39.0 | Kubernetes libraries 0.31.0, controller-runtime 0.19.0, Ginkgo 2.19.0, Gomega 1.33.1, envtest 1.31.0, Kustomize 5.4.3, and optional ingress NetworkPolicies. |
| 1.39.1 | No further migrations; CLI, scorecard images, and bundle metadata use this patch release. |

Scorecard images, release workflow downloads, bundle metadata, and documentation
now target SDK 1.39.1. The obsolete manager ComponentConfig is removed.
Kustomize uses replacements for certificate and conversion-webhook references.

## Compatibility with the retained Go version

Controller-gen remains v0.23.0: the guide's older generator is incompatible with
Go 1.26. The versioned tool installer keeps stable tool paths and selects the
configured version on every invocation. Envtest uses the release-0.19 tool and
the GitHub release asset index instead of the retired GCS endpoint.

This repository has no golangci-lint configuration or targets; migration edits
to existing `.golangci.yml` settings and linter versions therefore do not apply.
The existing formatting, vet, and gosec checks remain in use.

## Metrics

The manager disables metrics by default when run directly. The deployment enables
HTTPS metrics on port 8443 with TokenReview authentication and SubjectAccessReview
authorization. The `metrics-reader` role still grants GET access to `/metrics`;
monitoring service accounts need a binding to that role. The operator service
account has the review permissions through `metrics-auth-role`.

HTTP/2 is disabled by default. `--enable-http2` enables it for both servers;
`ENABLE_WEBHOOK_HTTP2` retains its existing webhook-only override.
`--metrics-secure=false` explicitly selects unauthenticated HTTP metrics for local
use. The existing ServiceMonitor uses a bearer token and skips verification of
the generated self-signed certificate. Production monitoring should use a trusted
certificate and CA verification.

## Optional network policies

`config/network-policy` contains the SDK ingress policy templates. They are not
enabled in the default deployment. To enable them, uncomment the network-policy
resource in `config/default/kustomization.yaml`. Label monitoring namespaces
`metrics=enabled` and adapt the webhook policy to the API server's source network
for the cluster/CNI. Labeling the namespace of a custom resource does not identify
the source of admission requests: those requests come from the API server.
The webhook policy uses the actual pod port 9443, rather than Service port 443.
Validate webhook connectivity before enabling the policies in a cluster.

## Bundle regeneration

Use Operator SDK v1.39.1 when running `make bundle`, for example by setting
`OPERATOR_SDK_CLI` to its binary. The existing SDK/Kustomize issue which drops
`DevfileRegistriesList` and `ClusterDevfileRegistriesList` examples still applies;
restore these two examples in `alm-examples` after regeneration, as described in
the README known issues. The checked-in bundle retains all three examples.

`make build-installer IMG=<operator-image>` writes `dist/install.yaml` and updates
the manager image in `config/manager/kustomization.yaml`, like `make deploy`.
