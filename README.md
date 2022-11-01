TCP server written in Go with bare metal (no Docker) deployment scripts.

# How to build
```bash
make
```

# How to prepare for deployment
This step need to be run on Debian-like Linux distro with `dpkg-deb` package installed

```bash
make package
```

# How to deploy
We assume `192.168.1.10` is the address of a remote host, 
running Debian-like Linux distro and having the SSH port opened.
We also assume we have the remote access to `root` user with authorization provided by ssh-agent.

```bash
./deploy.sh --user root --host 192.168.1.10
```

# How to monitor the deployment

Log in to the machine through SSH.

### Get status
```bash
systemctl status go-tcp-server
```

### Get logs
```bash
journalctl -u go-tcp-server
```
