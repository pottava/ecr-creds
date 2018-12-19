FROM alpine:3.8

RUN apk add --no-cache ca-certificates

ENV AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_DEFAULT_REGION=us-east-1

ADD dist/linux_amd64 /ecr-creds
RUN chmod +x /ecr-creds
ENTRYPOINT ["/ecr-creds"]
CMD ["--help"]
