# Builder image
FROM golang:1.9.2-alpine3.7 as builder

RUN apk add --no-cache \
      ca-certificates tzdata gnupg coreutils bash \
      git findutils make gcc musl-dev

ADD . /go/src/github.com/horizon-games/dgame-server
WORKDIR /go/src/github.com/horizon-games/dgame-server
RUN make dist

# TODO: check that subsequent builds are very small.. builder should stay in tact..
# and each new release should be very small in filesize.

# Runner image
FROM alpine:3.7

RUN apk add --no-cache ca-certificates

# Bin
COPY --from=builder /go/src/github.com/horizon-games/dgame-server/bin/* /usr/bin/

EXPOSE 8000

CMD ["/usr/bin/wsrelay", "-config=/etc/wsrelay.conf"]
