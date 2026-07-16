FROM golang:1.26.5 AS gob

ARG VERSION="unknown"

WORKDIR /build

COPY . /build/

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-X main.version=${VERSION}" -a -o api .

FROM alpine:3.23.3

COPY --from=gob /build/api /api

ENTRYPOINT [ "/api" ]
