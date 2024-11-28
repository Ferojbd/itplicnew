#!/bin/sh
#Script to return the main useable device IP address of the box, used for main outbound connections.
#on a LAN, this should match your directadmin.conf lan_ip setting.
#for normal servers, this will likely return your license IP (usually)
#Will also be the default IP that exim sends email through.

IP_BIN=/sbin/ip
if [ ! -e ${IP_BIN} ]; then
	IP_BIN=/usr/sbin/ip
fi
if [ ! -e ${IP_BIN} ]; then
	IP_BIN=/bin/ip
fi

IP=$(ip a | grep inet | grep -m1 brd | awk '{ print $2; };' | cut -d/ -f1)
RET=$?
	
if [ "${IP}" = "" ]; then
	IP=$(ip route get 8.8.8.8 | head -1 | grep -o 'src [^ ]*' | awk '{print $2}')
	RET=$?
fi
echo "${IP}"
exit $RET
