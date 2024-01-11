FROM golang:1.21.5-alpine AS builder
ENV DOCKER_DNS 8.8.8.8

ARG PROJECT_NAME
ENV PROJECT_NAME=$PROJECT_NAME 
ARG PROJECT_VERSION
ENV PROJECT_VERSION=$PROJECT_VERSION


WORKDIR /$PROJECT_NAME
RUN echo $PROJECT_NAME
COPY go.mod go.mod
COPY main.go main.go
RUN CGO_ENABLED=0 go build -ldflags '-s -w  -X main.appName='$PROJECT_NAME' -X main.appVersion='$PROJECT_VERSION -trimpath -o ./bin/$PROJECT_NAME .
RUN apk update && apk add upx
RUN upx ./bin/$PROJECT_NAME

FROM alpine:latest AS runner
# ENV TZ=Asia/Almaty
ARG PROJECT_NAME
ENV PROJECT_NAME=$PROJECT_NAME
ARG PROJECT_VERSION
ENV PROJECT_VERSION=$PROJECT_VERSION

RUN apk update && apk add tzdata
RUN apk update && apk add mc

WORKDIR /$PROJECT_NAME
COPY --from=builder $PROJECT_NAME/bin/$PROJECT_NAME bin/

EXPOSE $PORT
USER root:root
ENTRYPOINT bin/$PROJECT_NAME