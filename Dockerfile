FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
RUN CGO_ENABLED=0 go build -buildmode=plugin -o bin/nats.so ./libs/store/nats/nats.go
RUN CGO_ENABLED=0 go install ./cli/receiver -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch
# the test program:
COPY --from=app-builder /go/bin/nats.so /nats.so
COPY --from=app-builder /go/src/app/config.toml /config.toml
COPY --from=app-builder /go/bin/receiver /receiver
# the tls certificates:
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/receiver config.toml"]