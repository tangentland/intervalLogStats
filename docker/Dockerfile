FROM ubuntu:latest

RUN cd ~ \
    apt-get update \
    apt-get -y install wget software-properties-common git build-essential \
    apt-get -y dist-upgrade \
    wget -L "https://golang.org/dl/go1.21.4.linux-amd64.tar.gz" \
    tar -xf go1.21.4.linux-amd64.tar.gz \
    chown root:root ./go \
    mv ./go /usr/local \
    ln -s /usr/local/go/bin/go /usr/bin/go \
    echo "export GOPATH=$HOME/go" >> ~/.bashrc \
    echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> ~/.bashrc \
    source ~/.bashrc \
    git clone https://github.com/tangentland/intervalLogStats@v0.1.0 /opt/intervalLogStats \
    curl --proto '=https' --tlsv1.2 -fL -o /usr/local/bin/ocb \
      https://github.com/open-telemetry/opentelemetry-collector/releases/download/cmd%2Fbuilder%2Fv0.96.0/ocb_0.96.0_linux_amd64 \
    chmod +x /usr/local/bin/ocb \
    cd /opt/intervalLogStats \
    /usr/local/bin/ocb --config ./builder-config.yaml \
    chmod +x /opt/opentel/collector \
    cp  /opt/intervalLogStats/run_collector.sh /opt/opentel/ \
    chmod +x /opt/opentel/run_collector.sh


ENTRYPOINT ["/opt/opentel/run_collector.sh"]


