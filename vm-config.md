
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