

network 1 -> NAT
network2 -> Host-only adpater


### step 0

- port forward on NAT
- host port = 2222
- guest port = 22
- ssh name@127.0.0.1 -p 2222 

### step 1
`sudo nano /etc/netplan/01-netcfg.yaml`

```yml
network:
  version: 2
  renderer: networkd
  ethernets:
    # NAT adapter (usually first network interface)
    enp0s3:
      dhcp4: yes
      dhcp6: no

    # Host-only adapter (usually second network interface)
    enp0s8:
      dhcp4: no
      addresses:
        - 192.168.56.10/24  # Static IP for host-only network
      routes:
        - to: 192.168.56.0/24
          via: 192.168.56.1  # Host-only network gateway
```

`sudo chmod 600 /etc/netplan/01-netcfg.yaml`

`sudo netplan apply`

sudo apt install -y zsh python3-pip neovim net-tools nmap tcpdump

https://docs.docker.com/engine/install/ubuntu/
https://ohmyz.sh
http://www.lazyvim.org/installation


ref:
- https://prometheus.io/docs/introduction/first_steps/
- https://prometheus.io/docs/tutorials/getting_started/
```

wget https://github.com/prometheus/prometheus/releases/download/v3.0.0-rc.0/prometheus-3.0.0-rc.0.linux-amd64.tar.gz
useradd --no-create-home --shell /bin/false prometheus
mkdir /etc/prometheus
mkdir /var/lib/prometheus
chown prometheus:prometheus /etc/prometheus
chown prometheus:prometheus /var/lib/prometheus

tar -xzf prometheus-*.linux-amd64.tar.gz
cd prometheus-*.linux-amd64/
cp prometheus /usr/local/bin/
cp promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus
chown prometheus:prometheus /usr/local/bin/promtool
cp -r consoles /etc/prometheus

cp -r console_libraries /etc/prometheus

chown -R prometheus:prometheus /etc/prometheus/consoles
chown -R prometheus:prometheus /etc/prometheus/console_libraries

```

nano /etc/prometheus/prometheus.yml

```yml
global:
scrape_interval: 10s
scrape_configs:
- job_name: Prometheus_master
static_configs:
- targets: ['localhost:9090']
```

nano /etc/systemd/system/prometheus.service

```yml
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target
[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
--config.file /etc/prometheus/prometheus.yml \
--storage.tsdb.path /var/lib/prometheus/ \
--web.console.templates=/etc/prometheus/consoles \
--web.console.libraries=/etc/prometheus/console_libraries
[Install]
WantedBy=multi-user.target
```


```
sudo 


systemctl daemon-reload
sudo 


systemctl enable --now prometheus
sudo 


systemctl status prometheus
sudo 


ufw --add-port=9090/tcp --permanent
sudo 


ufw--reload
```