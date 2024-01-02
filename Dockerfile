FROM openeuler/openeuler:23.09 as BUILDER
RUN dnf update -y --repo OS --repo update && \
    dnf install -y golang --repo OS --repo update && \
    go env -w GOPROXY=https://goproxy.cn,direct

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
    echo 'set +o history' >> /root/.bashrc && \
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs && \
    rm -rf /tmp/*

USER modelfoundry
WORKDIR /home/modelfoundry

COPY  --chown=modelfoundry --from=BUILDER /go/src/github.com/openmerlin/merlin-server/merlin-server /home/modelfoundry
COPY  --chown=modelfoundry --from=BUILDER /go/src/github.com/openmerlin/merlin-server/cmd/cmd /home/modelfoundry

RUN chmod 550 /home/modelfoundry/merlin-server && \
    chmod 550 /home/modelfoundry/cmd && \
    echo "umask 027" >> /home/modelfoundry/.bashrc && \
    echo 'set +o history' >> /home/modelfoundry/.bashrc

ENTRYPOINT ["/home/modelfoundry/merlin-server"]
