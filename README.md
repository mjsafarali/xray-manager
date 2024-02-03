# Xray Manager

## Documentation

### Installation

```shell
# Install the requirements
apt-get -y update && apt-get -y upgrade
apt-get -y install make wget curl vim git openssl

# Install Docker
wget -O install-docker.sh https://get.docker.com
chmod +x install-docker.sh && ./install-docker.sh
rm -f install-docker.sh

# Install BBR
sudo sh -c 'echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf'
sudo sh -c 'echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf'
sudo sysctl -p
```

```shell
# Install Xray Manager
git clone https://github.com/miladrahimi/xray-manager.git
cd xray-manager
cp configs/main.json configs/main.local.json
docker compose up -d
```

### Web Panel

Access the web panel at the default port 8080. Log in using the username `admin` and password `password`.
* In the `Users` tab, you can create, edit, and view users.
* In the `Servers` tab, add Xray Nodes by specifying their Host (IP), HTTP Port, and HTTP Token.
* In the `Settings` tab, you can modify the Panel password, the Host (IP), the Ports, etc.

### Configuration

```shell
# Modify the web panel and user profile port (default 8080) here.
configs/main.local.json
```

### Update

``` shell
make update
# Execute this each time a new version is released.
```

## Links

* https://github.com/miladrahimi/xray-node
