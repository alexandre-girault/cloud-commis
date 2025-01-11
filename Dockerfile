ARG TARGET_ARCH


FROM --platform=${TARGET_ARCH} alpine:latest
ARG TARGET_ARCH
COPY bin/cloudcommis-linux-${TARGET_ARCH} /bin/cloudcommis
EXPOSE 8080
ENTRYPOINT ["/bin/cloudcommis"]