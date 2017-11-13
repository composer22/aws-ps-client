FROM alpine:latest
ARG release_tag
LABEL Name=aws-ps-client \
      Version=${release_tag}

ENV AWS_PS_RELEASE=${release_tag}
# Params for search.
ENV AWS_PS_KEY=
ENV AWS_PS_PATH=
ENV AWS_PS_VERSION=

RUN apk upgrade --update \
  && apk --update-cache update \
  && apk add --update bash curl \
  && rm -rf /var/cache/apk/* \
  && mkdir -p /usr/local/docker/aws-ps-client \
  && curl -f -o /usr/local/docker/aws-ps-client/aws-ps-client \
      -L https://github.com/composer22/aws-ps-client/releases/download/${release_tag}/aws-ps-client-linux-amd64 \
  && chmod +x /usr/local/docker/aws-ps-client/aws-ps-client

# Boostrap and config files
COPY ./docker-entrypoint.sh  /usr/local/docker/aws-ps-client/docker-entrypoint.sh
COPY ./examples/aws-ps-client.yaml  /usr/local/docker/aws-ps-client/.aws-ps-client.yaml

WORKDIR /usr/local/docker/aws-ps-client/
CMD []
ENTRYPOINT ["/usr/local/docker/aws-ps-client/docker-entrypoint.sh"]
