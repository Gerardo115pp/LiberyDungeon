FROM golang:1.21.4-alpine3.18 AS builder
RUN apk add --no-cache git gcc musl-dev pkgconfig
RUN apk add --no-cache make nasm yasm pkgconfig
RUN ping -c 1 google.com
RUN git clone https://git.ffmpeg.org/ffmpeg.git ffmpeg && \
    cd ffmpeg && \
    git checkout n4.1 && \
    ./configure && \
    make -j$(nproc) && \
    make install
WORKDIR /go/src
COPY . .
RUN go mod download
RUN go build -o /go/bin/medias_service /go/src/Medias

FROM alpine:3.16.0
RUN apk add --no-cache ca-certificates pkgconfig
RUN apk add --no-cache git gcc musl-dev make nasm yasm
RUN git clone https://git.ffmpeg.org/ffmpeg.git ffmpeg && \
    cd ffmpeg && \
    git checkout n4.1 && \
    ./configure && \
    make -j$(nproc) && \
    make install
COPY --from=builder /go/bin/medias_service /medias_service
RUN mkdir /operation_data
RUN mkdir /pandasmedia_root
ENTRYPOINT /medias_service
LABEL Name=medias_service Version=Dev.1.0
EXPOSE 6971
CMD ["medias_service"]