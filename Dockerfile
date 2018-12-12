FROM golang:1.9

ADD . /go/src/github.com/QubitProducts/aws_audit_exporter
RUN go install github.com/QubitProducts/aws_audit_exporter

ENTRYPOINT ["/go/bin/aws_audit_exporter"]

EXPOSE 9190
