# ceph-chaos-monkey

This software is designed to train Ceph engineers to recover Ceph clusters in
various ways by interacting with Ceph components and data to trigger errors
in the cluster. Therefore it could damage the data stored within the cluster
and that's why there are some limitations where you can run ceph-chaos-monkey:

* <=10 OSD daemons
* <=10 nodes
* <=500 GB of raw space

These restrictions are hardcoded and cannot be changed in runtime but anyway
if you have such a small clusters with important data please check twice where
you're running ceph-chaos-monkey.
