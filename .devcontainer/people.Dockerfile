FROM ugentlib/people-service:latest

COPY --from=golang:alpine /usr/local/go/ /usr/local/go/
ENV GOPATH=/usr/local/go
ENV PATH="${GOPATH}/bin:${PATH}"
RUN go install github.com/jackc/tern/v2@latest
ENV TERN_CONFIG /tern/people.tern.docker.conf
ENV TERN_MIGRATIONS /src/etc/migrations

VOLUME /tern
VOLUME /src