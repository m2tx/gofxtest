FROM golang:1.25-bookworm AS base
ARG NAME
ARG MODULE
WORKDIR $GOPATH/src/${MODULE}

ENV GOLANG_CI_LINT_VERSION=v2.0.2
ENV GOFLAGS="-mod=mod"
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANG_CI_LINT_VERSION}

ENV NANCY_VERSION=v1.0.52
RUN curl -o /bin/nancy https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-darwin-amd64 && chmod +x /bin/nancy

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

FROM base AS build 
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -v -a -tags netgo -installsuffix cgo -o /${NAME} ./cmd/server/

FROM base AS test
CMD go test -cpu 1 -coverprofile coverage.out -covermode=set -v ./... && \
    grep -v "_mock" coverage.out >> /coverage.out && \
    go tool cover -func /coverage.out -o /coverage.txt && \
    go tool cover -html /coverage.out -o /coverage.html

FROM base AS lint 
RUN golangci-lint run -c ./.golangci.yml

FROM base AS audit 
RUN go list -m all | nancy sleuth

FROM amd64/alpine:latest AS image
ARG NAME
RUN apk --no-cache add curl
COPY --from=build /${NAME} /server
ENTRYPOINT [ "/server" ]