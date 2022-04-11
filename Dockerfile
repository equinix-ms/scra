FROM golang:1.17 as builder
RUN mkdir /workdir
WORKDIR /workdir
COPY . /workdir
RUN CGO_ENABLED=0 go build -o app main.go

FROM scratch
COPY --from=builder /workdir/app /app
ENTRYPOINT ["/app", "watch"]