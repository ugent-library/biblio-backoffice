# build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o app -v ./main.go

# final stage
FROM alpine:latest
WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
EXPOSE 80
EXPOSE 443
CMD ["/dist/app", "server", "start"]
