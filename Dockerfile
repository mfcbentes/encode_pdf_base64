FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    libaio1 \
    gcc \
    libc6-dev \
    wget \
    unzip \
    tar \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/oracle && \
    wget https://download.oracle.com/otn_software/linux/instantclient/1923000/instantclient-basic-linux.x64-19.23.0.0.0dbru.zip -O /opt/oracle/instantclient.zip && \
    unzip /opt/oracle/instantclient.zip -d /opt/oracle && \
    rm /opt/oracle/instantclient.zip && \
    echo "/opt/oracle/instantclient_19_23" > /etc/ld.so.conf.d/oracle-instantclient.conf && \
    ldconfig

ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient_19_23
ENV PATH=/opt/oracle/instantclient_19_23:$PATH

RUN wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz && \
    rm go1.23.4.linux-amd64.tar.gz

ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /app

COPY . .

EXPOSE 8081

RUN go mod tidy

RUN go build -o app main.go

CMD ["./app"]