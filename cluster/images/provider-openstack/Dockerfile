FROM alpine:3.22
RUN apk --no-cache add ca-certificates bash

LABEL org.opencontainers.image.title="provider-openstack"
LABEL org.opencontainers.image.description="Crossplane provider for OpenStack management"
LABEL org.opencontainers.image.vendor="rossigee"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/rossigee/provider-openstack"

ARG TARGETOS
ARG TARGETARCH

COPY --chmod=0755 _output/bin/${TARGETOS}_${TARGETARCH}/provider /usr/local/bin/provider

USER 65532
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/provider"]