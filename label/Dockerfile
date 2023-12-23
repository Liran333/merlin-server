FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/openmerlin/merlin-server
COPY . .
RUN cd label && GO111MODULE=on CGO_ENABLED=0 go build -a -o merlin-label .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    dnf install -y git && \
    groupadd -g 1000 merlin && \
    useradd -u 1000 -g merlin -s /sbin/nologin -m merlin && \
    echo "umask 027" >> /home/merlin/.bashrc && \
    echo 'set +o history' >> /home/merlin/.bashrc && \
    echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd && \
    echo 'set +o history' >> /root/.bashrc && \
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs && rm -rf /tmp/* && \
    mkdir /opt/app -p && chmod 700 /opt/app && chown 1000:1000 /opt/app


COPY --chown=merlin --from=BUILDER /go/src/github.com/openmerlin/merlin-server/label/merlin-label /opt/app/merlin-label

RUN chmod 550 /opt/app/merlin-label

USER merlin

WORKDIR /opt/app/

ENTRYPOINT ["/opt/app/merlin-label"]