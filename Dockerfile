# Development image that includes golang tools
# and source code (for running unit tests, scripts, etc)
FROM golang:1.12 as dev

ENV BERGLAS_VERSION=0.2.0
RUN curl -OLs https://storage.googleapis.com/berglas/${BERGLAS_VERSION}/linux_amd64/berglas && \
    mv ./berglas /usr/local/bin/berglas && \
    chmod 0755 /usr/local/bin/berglas

ENV DOCKERIZE_VERSION=v0.6.1
RUN curl -OLs https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz && \
    tar -C /usr/local/bin -xzvf dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz && \
    rm dockerize-linux-amd64-${DOCKERIZE_VERSION}.tar.gz && \
    chmod 0755 /usr/local/bin/dockerize

ENV TINI_VERSION=v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /usr/local/bin/tini
RUN chmod 0755 /usr/local/bin/tini
ENTRYPOINT ["tini", "--"]

ENV GO111MODULE=on \
    GOFLAGS=-mod=vendor

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . .

RUN ./scripts/build
CMD ["berglas", "exec", "--local", "--", "climbcomp", "server"]


# Production image that only contains the built app
FROM gcr.io/distroless/base as prod

RUN mkdir -p /etc/climbcomp

COPY --from=dev /etc/climbcomp/config.yml /etc/climbcomp/config.yml
COPY --from=dev /usr/local/bin/climbcomp /usr/local/bin/climbcomp
COPY --from=dev /usr/local/bin/berglas /usr/local/bin/berglas

CMD ["berglas", "exec", "--local", "--", "climbcomp", "server"]
