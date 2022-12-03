# Build image
FROM alpine:latest AS build

# Build requirements
RUN apk add --no-cache ca-certificates

# Copy binary
COPY keyvault-certdeploy /keyvault-certdeploy

# ---

# Runtime image
FROM scratch
LABEL org.opencontainers.image.source = "https://github.com/emgag/keyvault-certdeploy"
LABEL org.opencontainers.image.license = "MIT"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /keyvault-certdeploy /keyvault-certdeploy

ENTRYPOINT ["/keyvault-certdeploy"]
CMD []
