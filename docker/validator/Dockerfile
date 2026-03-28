FROM golang:1.25 AS build

WORKDIR /src

COPY tools/cli/ /src/

ARG TARGETOS
ARG TARGETARCH

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/stability-flow-validator .

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/stability-flow-validator /usr/local/bin/sfv

ENTRYPOINT ["/usr/local/bin/sfv"]
CMD ["--help"]