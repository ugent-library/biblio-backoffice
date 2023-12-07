# build stage
FROM golang:1.20-alpine AS build
WORKDIR /build
COPY . .
RUN go install github.com/jackc/tern/v2@latest \
    && cp $GOPATH/bin/tern /usr/local/bin/

# final stage
FROM alpine:latest
WORKDIR /migrations
COPY --from=build /build/db/migrations .
COPY --from=build /usr/local/bin/tern  /usr/local/bin/
COPY --from=build /build/docker/tern.sh .
RUN chmod +x /migrations/tern.sh

ENV PGDATABASE $PGDATABASE

CMD "./tern.sh"