# This Dockerfile requires DOCKER_BUILDKIT=1 to be build.
# We do not use syntax header so that we do not have to wait
# for the Dockerfile frontend image to be pulled.
FROM golang:1.20-alpine3.18 AS build

RUN apk --update add make git gcc musl-dev tzdata && \
  adduser -D -H -g "" -s /sbin/nologin -u 1000 user
COPY . /go/src/regex2json
WORKDIR /go/src/regex2json
RUN \
  make build-static && \
  mv regex2json /go/bin/regex2json

FROM alpine:3.18 AS debug
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/bin/regex2json /
ENTRYPOINT ["/regex2json"]

FROM scratch AS production
RUN --mount=from=busybox:1.36-musl,src=/bin/,dst=/bin/ ["/bin/mkdir", "-m", "1755", "/tmp"]
COPY --from=build /etc/services /etc/services
COPY --from=build /etc/protocols /etc/protocols
# Apart from the USER statement, the rest is the same as for the debug image.
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/bin/regex2json /
USER user:user
ENTRYPOINT ["/regex2json"]
