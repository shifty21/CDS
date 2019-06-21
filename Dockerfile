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
#RUN rm -r /usr/bin/himeno

#ADD ./himeno_go/baseline_c /tmp/baseline_c

#RUN cd /tmp/baseline_c && \
#    make && \
#    cp /tmp/baseline_c/himeno /usr/bin/ && \
#    cd ~/
#   rm -r /tmp/mopp-2018-t0-harmonic-progression-sum
