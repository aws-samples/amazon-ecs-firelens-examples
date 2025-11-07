FROM public.ecr.aws/docker/library/golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum main.go .

RUN go mod download || (go env -w GOPROXY=direct && go mod download)

# Build the application with maximum optimization
# CGO_ENABLED=0 ensures a fully static binary
# -ldflags="-s -w" strips debug information and symbol table
# -trimpath removes build paths for reproducible builds
# -buildmode=exe ensures executable format
RUN CGO_ENABLED=0 go build \
    -a -installsuffix cgo \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o log-generator main.go

FROM public.ecr.aws/docker/library/alpine:3.22

RUN apk add --no-cache curl jq sqlite

RUN echo -e ".headers on\n.mode column\n.width 5 32 12 12 10" > ~/.sqliterc

COPY --from=builder /app/log-generator /log-generator

ENTRYPOINT ["/log-generator"]
