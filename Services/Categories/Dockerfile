FROM golang:1.23.0-alpine3.20 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o /go/bin/categories_service /go/src/Categories

FROM alpine:3.16.0
# RUN echo "Building Categories Service from platform $BUILDPLATFORM for $TARGETPLATFORM"
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin/categories_service /categories_service
RUN mkdir /operation_data
RUN mkdir /pandasmedia_root
ENTRYPOINT /categories_service
LABEL Name=categories_service Version=1.0
EXPOSE 6970
CMD ["categories_service"]