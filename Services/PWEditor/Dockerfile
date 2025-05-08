FROM golang:1.20.0-alpine3.16 AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /go/src/pw_editor
COPY . .
RUN go mod download
RUN go build -o /go/bin/pw_editor

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/pw_editor /pw_editor
RUN mkdir /operation_data
ENTRYPOINT /pw_editor
LABEL Name=pw_editor Version=1.0
EXPOSE 6970
CMD ["pw_editor"]