FROM golang:1.11  AS build-env

ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 go install .

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/aws_audit_exporter /]
ENTRYPOINT ["/aws_audit_exporter"]
EXPOSE 9190
