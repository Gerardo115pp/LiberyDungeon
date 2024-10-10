FROM golang:1.23.0-alpine3.20 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o /go/bin/libery_metadata_service /go/src/Metadata

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/libery_metadata_service /libery_metadata_service
RUN mkdir /operation_data
ENTRYPOINT /libery_metadata_service
LABEL Name=libery_metadata_service Version=1.0
EXPOSE 6970
CMD ["libery_metadata_service"]