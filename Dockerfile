FROM --platform=$BUILDPLATFORM golang:1.17 as builder
RUN mkdir /workdir
WORKDIR /workdir
COPY . /workdir
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o app -v main.go

FROM scratch
COPY --from=builder /workdir/app /app
ENTRYPOINT ["/app", "watch"]