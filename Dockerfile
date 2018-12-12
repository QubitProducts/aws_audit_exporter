FROM golang:1.11

ADD . /src
WORKDIR /src
RUN go install .

ENTRYPOINT ["/go/bin/aws_audit_exporter"]

EXPOSE 9190
