FROM golang:1.20.0-alpine3.16 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src/JD_service
COPY . .
RUN go mod download
RUN go build -o /go/bin/JD_service

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/JD_service /JD_service
RUN mkdir /operation_data
ENTRYPOINT /JD_service
LABEL Name=JD_service Version=1.0
EXPOSE 6970
CMD ["JD_service"]