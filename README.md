![moobk logo](https://raw.githubusercontent.com/raohwork/moobk/master/moobk.png)

moobk is a CoW filesystem backup tool. Currently only btrfs/zfs are supported.

# Install

- `go get github.com/raohwork/moobk` if you have golang tools installed.
- Download from release page on GitHub.

# Synopsis

Run "moobk help", "moobk help driver", "moobk help repo" to see online docs.

```sh
# show usage
moobk help

# show usage of a command (like snap)
moobk help snap

# take a snapshot of /.droot/@root at localhost, save it to /.droot/snapshot/@root-timestamp at localhost
moobk snap -t btrfs -r local:///.droot/snapshot /.droot/@root

# take a snapshot of /.droot/@root at remote using ssh, save it to /.droot/snapshot/@root-timestamp at remote
# you might have to install moobk on remote host and allow root login with pubkey
moobk snap -t btrfs -r ssh:///.droot/snapshot /.droot/@root

# list all recognized snapshots in /.droot/snapshot at localhost
moobk list -t btrfs -r local:///.droot/snapshot

# list all recognized snapshots in /data/backup at remote host using ssh
moobk list -t btrfs -r ssh://user@ip:port/data/backup

# list all recognized snapshots in /data/backup at remote host using ssh
# you might have to install moobk on remote host and allow passwordless sudo
moobk list -t btrfs -r ssh+sudo://user@ip:port/data/backup

# delete synchronized snapshots in /.droot/snapshot at localhost, according to what exists at remote host
moobk purge -t btrfs -r local:///.droot/snapshot ssh+sudo://user@ip:port/data/backup

# same as above, but reserves up to 2 more snapshot (which means 3 total)
moobk purge -t btrfs -r local:///.droot/snapshot ssh+sudo://user@ip:port/data/backup 2

# same as above, but reserves up to 2 days ago from now (also supports h/w/m for hour/week/month)
moobk purge -t btrfs -r local:///.droot/snapshot ssh+sudo://user@ip:port/data/backup 2d

# transfer snapshots from localhost to remote
moobk sync -t btrfs -r local:///.droot/snapshot ssh+sudo://user@ip:port/data/backup

# transfer matching snapshots from localhost to remote
moobk sync -t btrfs -r local:///.droot/snapshot -name rootfs ssh+sudo://user@ip:port/data/backup

# transfer snapshots from remote to localhost
moobk sync -t btrfs -r ssh+sudo://user@ip:port/data/backup local:///.droot/snapshot

# transfer snapshot from one remote machine to another, with ssh options
# A and B are not directly connected, so it doubles network usage on local machine.
moobk sync -t btrfs -r ssh+sudo://user@1.2.3.4:port/data/backup?ssh_4 ssh+sudo://user@5.6.7.8:port/data/backup?ssh_4
```

# WARNING

- The order of command/options/arguments is extremly strict. Still searching for a good CLI framework to use (stability/learning curve/license are considered)
- It does not know source path from snapshot, you have to remember it by yourself. Write simple shell script for automation is suggested.
- Snapshot naming scheme is fixed.
- It compares snapshots only by their name and timestamp.

# License

GPLv2+

Copyright 

- 2021- Ronmi Ren <ronmi.ren@gmail.com>
