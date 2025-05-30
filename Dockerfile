FROM --platform=$BUILDPLATFORM golang:latest AS builder
ARG TARGETARCH
ARG VERSION


WORKDIR /convcommitlint
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-X github.com/coolapso/convcommitlint/cmd.Version=${VERSION}" -a -o convcommitlint

FROM alpine:latest

COPY --from=builder convcommitlint/convcommitlint /usr/bin/convcommitlint
RUN mkdir /data

WORKDIR data
ENTRYPOINT ["/usr/bin/convcommitlint"] 
