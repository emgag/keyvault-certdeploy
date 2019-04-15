# Build image
FROM alpine:latest AS build

# Build requirements
RUN apk add --no-cache ca-certificates

# Copy binary
COPY dist/linux_amd64/keyvault-certdeploy /

# ---

# Runtime image
FROM scratch
LABEL maintainer="Matthias Blaser <mb@emgag.com>"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /keyvault-certdeploy /keyvault-certdeploy

ENTRYPOINT ["/keyvault-certdeploy"]
CMD []