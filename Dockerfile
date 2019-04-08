# Build image
FROM golang:1.12.2-alpine3.9 AS build

# Build requirements
RUN apk add --no-cache git make

# Copy context
COPY . /work

# Build
RUN cd /work && make build

# ---

# Runtime image
FROM scratch as app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /work/keyvault-certdeploy /keyvault-certdeploy

CMD ["/keyvault-certdeploy"]