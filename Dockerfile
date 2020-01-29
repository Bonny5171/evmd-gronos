###############
# BUILDER     #
###############
FROM golang:1.13.7-alpine3.11 AS builder
LABEL maintainer="Roberto Besser <roberto.besser@everymind.com.br>"
RUN apk --no-cache add ca-certificates git mercurial
RUN git config --global url."https://everyminddev:VJN45xFcjXudwBbTG6vU@bitbucket.org/everymind/".insteadOf https://bitbucket.org/everymind/
COPY . /src
WORKDIR /src
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org/
ENV GOPRIVATE=bitbucket.org/everymind
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/gronos ./*.go
RUN cp $(go env | grep GOROOT | sed 's/GOROOT=//g' | sed 's/"//g')/lib/time/zoneinfo.zip /app/zoneinfo.zip

###############
# FINAL IMAGE #
###############
FROM scratch
ENV ZONEINFO=./zoneinfo.zip
COPY --from=builder /app/gronos .
COPY --from=builder /app/zoneinfo.zip .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt