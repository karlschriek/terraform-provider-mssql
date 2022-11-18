# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app
COPY . ./

RUN go mod tidy

RUN make build

## Deploy
FROM ubuntu:20.04

WORKDIR /app

COPY --from=build /app/terraform-provider-mssql ./terraform-provider-mssql
