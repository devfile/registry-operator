#
# Copyright Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the manager binary
FROM registry.access.redhat.com/ubi9/go-toolset:1.22.9@sha256:e4193e71ea9f2e2504f6b4ee93cadef0fe5d7b37bba57484f4d4229801a7c063 as builder
ARG TARGETARCH=amd64
USER root

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o manager main.go

ARG ENABLE_WEBHOOKS=true
ENV ENABLE_WEBHOOKS=${ENABLE_WEBHOOKS}

# disable http/2 on the webhook server by default
ARG ENABLE_WEBHOOK_HTTP2=false
ENV ENABLE_WEBHOOK_HTTP2=${ENABLE_WEBHOOK_HTTP2}

# Use ubi-micro as minimal base image to package the manager binary
FROM registry.access.redhat.com/ubi9/ubi-micro:9.6@sha256:955512628a9104d74f7b3b0a91db27a6bbecdd6a1975ce0f1b2658d3cd060b98
WORKDIR /
COPY --from=builder /workspace/manager .
USER 1001

ENTRYPOINT ["/manager"]
