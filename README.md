# ceph-chaos-monkey

This software is designed to train Ceph engineers to recover Ceph clusters in
various ways by interacting with Ceph components and data to trigger errors
in the cluster. Therefore it could damage the data stored within the cluster
and that's why there are some limitations where you can run ceph-chaos-monkey:

* <=10 OSD daemons
* <=500 GB of raw space

These restrictions are hardcoded and cannot be changed in runtime but anyway
if you have such a small clusters with important data please check twice where
you're running ceph-chaos-monkey.

## Usage

```shell
usage: ceph-chaos-monkey [<flags>] <command> [<args> ...]

Ceph Chaos Monkey


Flags:
  --[no-]help                    Show context-sensitive help (also try --help-long and --help-man).
  --[no-]trace                   set verbosity level to trace
  --ceph-binary="/usr/bin/ceph"  path to the ceph binary
  --rados-binary="/usr/bin/rados"
                                 path to the rados binary

Commands:
help [<command>...]
    Show help.

run --fuss-interval=FUSS-INTERVAL --game-duration=GAME-DURATION
    run the game

version
    print version and exit
```

ceph-chaos-monkey distributed as a container image so you could simply update
to it via `ceph orch upgrade`.
