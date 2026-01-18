# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# 1. Manually create a perfect go.mod file inside the container
RUN echo 'module docker-monitor' > go.mod && \
    echo 'go 1.24' >> go.mod && \
    echo 'replace github.com/docker/docker => github.com/moby/moby v27.5.1+incompatible' >> go.mod

# 2. Copy only the source code
COPY main.go .
COPY index.html .

# 3. Force Go to resolve the dependencies using our replacement map
ENV GOPROXY=direct
RUN go mod tidy

# 4. Build for 64-bit Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker-monitor main.go

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/docker-monitor .
COPY --from=builder /app/index.html .
EXPOSE 8080
USER root 
CMD ["./docker-monitor"]