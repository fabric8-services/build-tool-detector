FROM centos:7
LABEL maintainer "Devtools <devtools@redhat.com>"
LABEL author "Konrad Kleine <kkleine@redhat.com>"
ENV LANG=en_US.utf8
ARG USE_GO_VERSION_FROM_WEBSITE=0

# Some packages might seem weird but they are required by the RVM installer.
RUN yum install epel-release -y && \
    yum --enablerepo=centosplus --enablerepo=epel install -y \
      findutils \
      git \
      $(test "$USE_GO_VERSION_FROM_WEBSITE" != 1  && echo "golang") \
      make \
      procps-ng \
      tar \
      wget \
      which \
    && yum clean all

RUN echo $USE_GO_VERSION_FROM_WEBSITE
RUN if [[ "$USE_GO_VERSION_FROM_WEBSITE" == 1 ]]; then cd /tmp \
    && wget https://dl.google.com/go/go1.11.3.linux-amd64.tar.gz \
    && echo "b5a64335f1490277b585832d1f6c7f8c6c11206cba5cd3f771dcb87b98ad1a33  go1.11.3.linux-amd64.tar.gz" > checksum \
    && sha256sum -c checksum \
    && tar -C /usr/local -xzf go1.11.3.linux-amd64.tar.gz \
    && rm -f go1.11.3.linux-amd64.tar.gz; \
    fi
ENV PATH=$PATH:/usr/local/go/bin

ENTRYPOINT ["/bin/bash"]