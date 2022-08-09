FROM golang:1.17 as buid_service
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

WORKDIR /service
RUN go build -ldflags "-X main.buid=${BUILD_REF}"

FROM alpine:3.16
ARG BUILD_REF
RUN addgroup -g 1000 -S service && \
    adduser -u 1000 -h /service -G service -S service
COPY --from=buid_service --chown=service:service /service /service
WORKDIR /service
USER service
CMD ["./ultimateservice"]