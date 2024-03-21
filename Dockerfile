# build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o app -v
RUN GOBIN=/build/ go install github.com/jackc/tern/v2@latest

# final stage
FROM alpine:latest

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ARG IMAGE_NAME
ENV SOURCE_BRANCH $SOURCE_BRANCH
ENV SOURCE_COMMIT $SOURCE_COMMIT
ENV IMAGE_NAME $IMAGE_NAME

ENV TERN_CONFIG /dist/tern.docker.conf
ENV TERN_MIGRATIONS /dist/etc/snapstore/migrations

WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "server", "start"]
