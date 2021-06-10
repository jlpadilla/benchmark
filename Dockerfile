FROM registry.ci.openshift.org/open-cluster-management/builder:go1.16-linux AS builder

WORKDIR /go/src/github.com/jlpadilla/benchmark
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o main main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.3

COPY --from=builder /go/src/github.com/jlpadilla/benchmark/main /bin/main

ENV USER_UID=1001 \
    GOGC=25

USER ${USER_UID}
ENTRYPOINT ["/bin/main"]