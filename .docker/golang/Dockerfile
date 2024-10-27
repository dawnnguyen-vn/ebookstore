FROM golang:1.23.2-bookworm

RUN if ! getent group 999 > /dev/null 2>&1; then \
        groupadd -g 999 mygroup; \
    fi

RUN if ! id -u 1000 > /dev/null 2>&1; then \
    useradd -m -u 1000 -g mygroup myuser; \
    fi

RUN apt-get -y update \
    && apt-get install -y zip