FROM quay.io/ceph/ceph:v19.2.1

ADD dist/ceph-chaos-monkey_linux_amd64_v1/ceph-chaos-monkey /usr/local/bin/ceph-chaos-monkey
