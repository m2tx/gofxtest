FROM golang:1.25-bookworm AS base
ARG NAME
ARG MODULE
WORKDIR $GOPATH/src/${MODULE}

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

FROM amd64/alpine:latest AS image
ARG NAME
RUN apk --no-cache add curl
COPY --from=build /${NAME} /server
ENTRYPOINT [ "/server" ]