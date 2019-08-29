# Development image that includes golang tools
# and source code (for running unit tests, scripts, etc)
FROM golang:1.12 as dev

ENV GO111MODULE=on \
    GOFLAGS=-mod=vendor

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . .

RUN go build -v \
    -o /usr/local/bin/climbcomp \
    .
CMD ["/usr/local/bin/climbcomp"]


# Production image that only contains the built app
FROM gcr.io/distroless/base as prod

COPY --from=dev /usr/local/bin/climbcomp /usr/local/bin/climbcomp

CMD ["/usr/local/bin/climbcomp"]
