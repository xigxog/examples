## Build time
FROM golang:1.21 as builder

WORKDIR /app

# Cache dependencies.
COPY go.mod .
RUN go mod download

COPY ./ .
RUN CGO_ENABLED=0 go build -o http-srv

## Runtime
FROM alpine
COPY --from=builder /app/http-srv /
ENTRYPOINT [ "/http-srv" ]
