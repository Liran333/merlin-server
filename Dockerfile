FROM openeuler/openeuler:23.09 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

# build binary
COPY . /go/src/github.com/openmerlin/merlin-server
RUN cd /go/src/github.com/openmerlin/merlin-server && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"
# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 1000 openmerlin && \
    useradd -u 1000 -g openmerlin -s /sbin/nologin -m openmerlin

RUN echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd

RUN echo 'set +o history' >> /root/.bashrc
RUN sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs
RUN rm -rf /tmp/*

USER openmerlin
WORKDIR /home/openmerlin

COPY  --chown=openmerlin --from=BUILDER /go/src/github.com/openmerlin/merlin-server/merlin-server /home/openmerlin

RUN chmod 550 /home/openmerlin/merlin-server

RUN echo "umask 027" >> /home/openmerlin/.bashrc
RUN echo 'set +o history' >> /home/openmerlin/.bashrc

ENTRYPOINT ["/home/openmerlin/merlin-server"]
