FROM tudinfse/cds_server



ADD ./mopp-2018-t0-harmonic-progression-sum /tmp/mopp-2018-t0-harmonic-progression-sum
ADD ./cds_server.json /etc/

# Build code and copy into container system paths
RUN apt update && \
    apt install -y build-essential
RUN cd /tmp/mopp-2018-t0-harmonic-progression-sum && \
    make && \
    cp /tmp/mopp-2018-t0-harmonic-progression-sum/harmonic-progression-sum /usr/bin/ && \
    cd ~/
#   rm -r /tmp/mopp-2018-t0-harmonic-progression-sum
