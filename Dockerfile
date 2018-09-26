FROM golang:1.11.0-alpine3.8 as builder

WORKDIR /develop
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o download-etcd-snapshot-b2 -ldflags '-w'

FROM gcr.io/etcd-development/etcd:v3.3.9

ENV ETCDCTL_API=3 \
    EBS_B2_DOWNLOAD_RETRY_INTERVAL=5000

COPY entrypoint.sh /usr/local/bin/entrypoint
COPY --from=builder /develop/download-etcd-snapshot-b2 /usr/local/bin/download-etcd-snapshot-b2
RUN chmod +x /usr/local/bin/download-etcd-snapshot-b2 /usr/local/bin/entrypoint

VOLUME [ "/etcd" ]
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=10 CMD [ "etcdctl", "endpoint" , "health"]

CMD ["entrypoint"]
