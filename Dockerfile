FROM golang:1.24-alpine AS build

WORKDIR /app

# Cache module download layer
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download

# Copy the rest of the source code and build
COPY . .
RUN go build -mod=mod -o /app/server .


FROM alpine:3.20

WORKDIR /app

# Copy compiled binary from build stage
COPY --from=build /app/server ./server

ENV PORT=5000
EXPOSE 5000

CMD ["./server"]