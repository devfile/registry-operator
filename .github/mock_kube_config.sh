#!/bin/bash
# copied script from the release-service: https://github.com/redhat-appstudio/release-service-bundles/blob/main/.github/scripts/mock_kube_config.sh
printf 'Creating mock kube config for controller tests\n'
mkdir -vp $HOME/.kube || true
cat <<-EOF > $HOME/.kube/config
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: _
  name: _
contexts:
- context:
    cluster: _
  name: _
current-context: _
EOF