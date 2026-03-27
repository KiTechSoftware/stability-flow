FROM golang:1.25 AS build

WORKDIR /src

COPY tools/cli/ /src/

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/stability-flow-validator .

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/stability-flow-validator /usr/local/bin/stability-flow-validator

ENTRYPOINT ["/usr/local/bin/stability-flow-validator"]
CMD ["--help"]
