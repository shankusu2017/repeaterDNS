#!/bin/bash
#设置主网卡
ip route add default via 192.168.0.1 dev eth0 table 10
ip route add 192.168.0.0/24 dev eth0 table 10
ip rule add from 192.168.0.236 table 10
#设置副网卡
ip route add default via 192.168.1.1 dev eth1 table 20
ip route  add 192.168.1.0/24 dev eth1 table 20
ip rule add from 192.168.1.132 table 20
