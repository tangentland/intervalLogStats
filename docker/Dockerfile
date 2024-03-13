FROM ubuntu:latest

RUN cd /root
RUN apt-get update
RUN apt-get -y dist-upgrade
RUN apt-get -y install wget curl inetutils-ping software-properties-common git build-essential
RUN wget -L "https://golang.org/dl/go1.21.4.linux-amd64.tar.gz"
RUN tar -xf go1.21.4.linux-amd64.tar.gz
RUN rm go1.21.4.linux-amd64.tar.gz
RUN chown root:root ./go
RUN mv ./go /usr/local
RUN ln -s /usr/local/go/bin/go /usr/bin/go
RUN mkdir /opt/gocode
RUN export GOPATH=/opt
RUN export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
RUN git clone --depth 1 --branch v0.1.0 https://github.com/tangentland/intervalLogStats /opt/intervalLogStats
RUN curl --proto '=https' --tlsv1.2 -fL -o /usr/local/bin/ocb \
      https://github.com/open-telemetry/opentelemetry-collector/releases/download/cmd%2Fbuilder%2Fv0.96.0/ocb_0.96.0_linux_amd64
RUN chmod +x /usr/local/bin/ocb
RUN cd /opt/intervalLogStats
RUN /usr/local/bin/ocb --config /opt/intervalLogStats/builder-config.yaml
RUN chmod +x /opt/opentel/collector
RUN cp  /opt/intervalLogStats/run_collector.sh /opt/opentel/
RUN chmod +x /opt/opentel/run_collector.sh


ENTRYPOINT ["/opt/opentel/run_collector.sh"]


