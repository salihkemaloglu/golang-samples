FROM golang:onbuild

RUN adduser --disabled-password --gecos '' api
USER api

COPY . /app
WORKDIR /app

