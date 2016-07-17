FROM scratch
MAINTAINER Goutham Veeramachaneni <goutham@boomerangcommerce.com>

COPY pixie pixie
EXPOSE 8080

ENTRYPOINT ["/pixie"]
CMD ["-config", "/config.json"]
