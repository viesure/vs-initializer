FROM golang:1.19 as build

ARG VERSION=latest

RUN apt-get -y update; \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    tzdata \
    openssl

WORKDIR /
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=v${VERSION}'" -o /vs-initializer

FROM scratch
LABEL org.opencontainers.image.authors="roland.zagler@viesure.io"

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /vs-initializer /vs-initializer

ENTRYPOINT [ "/vs-initializer" ]
