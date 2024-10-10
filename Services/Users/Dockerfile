FROM golang:1.20.0-alpine3.16 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src/users_service
COPY . .
RUN go mod download
RUN go build -o /go/bin/users_service

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/users_service /users_service
RUN mkdir /operation_data
ENTRYPOINT /users_service
LABEL Name=users_service Version=1.0
EXPOSE 6970
CMD ["users_service"]