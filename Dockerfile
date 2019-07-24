###############
# BUILDER     #
###############
FROM golang:1.12.6-alpine3.9 AS builder
RUN apk --no-cache add ca-certificates git mercurial
RUN git config --global url."https://everyminddev:VJN45xFcjXudwBbTG6vU@bitbucket.org/everymind/".insteadOf https://bitbucket.org/everymind/
COPY . /src
WORKDIR /src
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/gronos ./app/*.go
RUN cp $(go env | grep GOROOT | sed 's/GOROOT=//g' | sed 's/"//g')/lib/time/zoneinfo.zip /app/zoneinfo.zip

###############
# FINAL IMAGE #
###############
FROM scratch
ENV ZONEINFO=./zoneinfo.zip
COPY --from=builder /app/gronos .
COPY --from=builder /app/zoneinfo.zip .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt