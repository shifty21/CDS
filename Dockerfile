FROM tudinfse/cds_server



ADD ./mopp-2018-t0-harmonic-progression-sum /tmp/mopp-2018-t0-harmonic-progression-sum
ADD ./amd_go /tmp/amd_go
ADD ./himeno_go /tmp/himeno_go
ADD ./cds_server.json /etc/

# Build code and copy into container system paths
RUN apt update && \
    apt install -y build-essential golang-go
RUN cd /tmp/mopp-2018-t0-harmonic-progression-sum && \
    make && \
    cp /tmp/mopp-2018-t0-harmonic-progression-sum/harmonic-progression-sum /usr/bin/ && \
    cd ~/
RUN cd /tmp/amd_go && \
    make && \
    cp /tmp/amd_go/bin/amd /usr/bin/ && \
    cd ~/
RUN cd /tmp/himeno_go && \
    make && \
    cp /tmp/himeno_go/bin/himeno /usr/bin/ && \
    cd ~/
ADD ./amd_go/baseline_c/amd /usr/bin/
#   rm -r /tmp/mopp-2018-t0-harmonic-progression-sum
