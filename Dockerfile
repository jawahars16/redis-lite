FROM golang:1.21
WORKDIR /app
COPY . /app
RUN make build

FROM scratch
COPY --from=0 /app/bin/redis-lite redis-lite
ENTRYPOINT ["./redis-lite"]