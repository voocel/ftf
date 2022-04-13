FROM golang:1.14.15-alpine AS builder

LABEL stage=gobuilder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct


# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
ADD go.mod .
ADD go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -ldflags="-s -w" -o /app/ftf

# Build a small image
FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

COPY --from=builder /app/ftf /ftf

# Export necessary port
EXPOSE 1234

# Command to run
CMD ["/ftf"]