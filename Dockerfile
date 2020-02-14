FROM golang:1.13.4-alpine3.10 as builder
COPY . /go/src/gocrontranslate
WORKDIR /go/src/gocrontranslate
RUN GO111MODULE=on CGO_ENABLED=0 go install

FROM busybox:1.31.0
RUN addgroup -S appuser && adduser -S -s /bin/false -G appuser appuser
COPY --from=builder /go/bin/gocrontranslate /go/bin/gocrontranslate
RUN chown -R appuser:appuser /go/bin
WORKDIR /go/bin
USER appuser:appuser
ENTRYPOINT ["/go/bin/gocrontranslate"]
