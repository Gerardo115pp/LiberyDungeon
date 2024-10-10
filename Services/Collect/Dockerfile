FROM python:3.11.7-alpine3.19 
# FROM python:3.10.13-slim-bookworm
RUN apk update && apk upgrade
RUN apk add --no-cache linux-headers
# RUN apt-get update && apt-get install -y \
#     build-essential \
#     libssl-dev \
#     libffi-dev \
#     python3-dev \
#     cargo \
#     gcc \
#     musl-dev \
#     make \
#     ca-certificates \
#     && rm -rf /var/lib/apt/lists/*
RUN pip3 install --upgrade pip
RUN pip3 install --upgrade setuptools
RUN apk add --no-cache ca-certificates build-base rust cargo python3-dev libc-dev libffi-dev openssl-dev g++ gcc musl-dev make
RUN pip install --upgrade pip

WORKDIR /app

COPY ./Collect/requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY ./Shared /app/Shared
COPY ./Collect /app/Collect

ENV PYTHONPATH "${PYTHONPATH}:/app/Shared/python_shared"

RUN mkdir /app/operation_data

EXPOSE 6972

WORKDIR /app/Collect

CMD ["uvicorn", "collect_service:app", "--host", "0.0.0.0", "--port", "6972"]