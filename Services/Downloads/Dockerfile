FROM golang:1.20.0-alpine3.16 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o /go/bin/downloads_service /go/src/Downloads

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/downloads_service /downloads_service
RUN mkdir /operation_data
ENTRYPOINT /downloads_service
LABEL Name=downloads_service Version=1.0
EXPOSE 6973
CMD ["downloads_service"]