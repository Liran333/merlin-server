FROM golang:1.21 as BUILDER
ARG MERLIN_GOPROXY=https://goproxy.cn,direct
ARG GH_USER
ARG GH_TOKEN
RUN go env -w GOPROXY=${MERLIN_GOPROXY} && go env -w GOPRIVATE=github.com/openmerlin
RUN echo "machine github.com login ${GH_USER} password ${GH_TOKEN}" > $HOME/.netrc

# build binary
COPY . /go/src/github.com/openmerlin/merlin-server
RUN cd /go/src/github.com/openmerlin/merlin-server && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'" && \
    cd /go/src/github.com/openmerlin/merlin-server/cmd && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update --repo OS --repo update && \
    dnf in -y shadow --repo OS --repo update && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 1000 modelfoundry && \
    useradd -u 1000 -g modelfoundry -s /sbin/nologin -m modelfoundry && \
    echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd && \
    echo "umask 027" >> /root/.bashrc &&\
    echo 'set +o history' >> /root/.bashrc && \
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs && \
    rm -rf /tmp/*

USER modelfoundry
WORKDIR /home/modelfoundry

COPY  --chown=modelfoundry --from=BUILDER /go/src/github.com/openmerlin/merlin-server/merlin-server /home/modelfoundry
COPY  --chown=modelfoundry --from=BUILDER /go/src/github.com/openmerlin/merlin-server/cmd/cmd /home/modelfoundry

ARG MODE=release
RUN chmod 550 /home/modelfoundry/merlin-server && \
    [ ${MODE} == "release" ] && rm /home/modelfoundry/cmd || chmod 550 /home/modelfoundry/cmd && \
    echo "umask 027" >> /home/modelfoundry/.bashrc && \
    echo 'set +o history' >> /home/modelfoundry/.bashrc

ENTRYPOINT ["/home/modelfoundry/merlin-server"]
