FROM golang:1.7-wheezy

ADD . /go/src/github.com/QubitProducts/aws_audit_exporter
RUN go install github.com/QubitProducts/aws_audit_exporter

ENTRYPOINT ["/go/bin/aws_audit_exporter"]

EXPOSE 9190
