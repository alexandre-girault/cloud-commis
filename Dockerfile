FROM golang:1.22-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build \ 
  -ldflags="-X 'cloud-commis/config.Version=$(VERSION)' -X 'cloud-commis/config.BuildDate=$(shell date +%Y-%m-%d_%H:%M:%S)'" \
-o /bin/cloudcommis main.go



FROM alpine:latest

COPY --from=builder /bin/cloudcommis /bin/cloudcommis
EXPOSE 8080
ENTRYPOINT ["/bin/cloudcommis"]