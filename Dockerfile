# build stage
FROM golang:1.27.3-alpine AS build
WORKDIR /build
COPY . .
RUN go get -v ./...
RUN go build -o app -v
RUN GOBIN=/build/ go install github.com/jackc/tern/v2@v2.4.3

# final stage
FROM alpine:3.24
# hooks/build passes a new value on every build, so this layer is never served from cache.
ARG SNAPSHOT_EPOCH
RUN apk upgrade --no-cache

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ARG IMAGE_NAME
ENV SOURCE_BRANCH $SOURCE_BRANCH
ENV SOURCE_COMMIT $SOURCE_COMMIT
ENV IMAGE_NAME $IMAGE_NAME

ENV TERN_CONFIG /dist/tern.docker.conf
ENV TERN_MIGRATIONS /dist/db/migrations

WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "server", "start"]
