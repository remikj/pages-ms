FROM golang:1.18.3 as builder

WORKDIR /workspace

# Download modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the go sources
COPY src src

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o pages-ms ./src/main.go

FROM gcr.io/distroless/static:latest

WORKDIR /
USER 1000:1000

COPY --from=builder /workspace/pages-ms .

ENTRYPOINT ["/pages-ms"]
