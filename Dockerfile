FROM golang:alpine

COPY . /tmp/app/

RUN set -ex \
    && apk add --no-cache --no-progress unixodbc \
    && apk add --no-cache --no-progress --virtual build-deps \
        gcc \
        make \
        musl-dev \
        unixodbc-dev \
    && mkdir -p /etc/odbc-exporter \
    && cd /tmp/app \
    && go build -o /odbc_exporter cmd/odbc_exporter/main.go \
    && mv configs/test.yml /etc/odbc-exporter/config.yml \
    && apk del build-deps \
    && rm -r /tmp/*

RUN set -ex \
    && apk add --no-cache --no-progress wget keyutils-libs libstdc++ \
    && wget -qO - https://www.exasol.com/support/secure/attachment/111075/EXASOL_ODBC-6.2.9.tar.gz | tar xzf - -C /tmp \
    && mkdir -p /opt/lib/exasol-odbc \
    && mv /tmp/EXASOL_ODBC-6.2.9/lib/linux/x86_64/* /opt/lib/exasol-odbc/ \
    && ln -s /lib/ld-musl-x86_64.so.1 /usr/lib/libresolv.so.2 \
    && apk del wget \
    && rm -r /tmp/EXASOL_ODBC-6.2.9 \
    && printf '[ODBC Drivers]\nEXASolution Driver = Installed\n\n[EXASolution Driver]\nDriver = /opt/lib/exasol-odbc/libexaodbc-uo2214lv2.so\n' > /etc/odbcinst.ini

ENTRYPOINT ["/odbc_exporter"]
