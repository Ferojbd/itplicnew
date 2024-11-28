#!/bin/sh

umask 022;
OS=`uname`;
M1=data.installatron.com;

PROXY="";
PROXYUSER="";
for w; do if [ "$w" = "--proxy" ]; then PROXY="-1"; elif [ "$PROXY" = "-1" ]; then PROXY=$w; elif [ "$w" = "--proxy-user" ]; then PROXYUSER="-1"; elif [ "$PROXYUSER" = "-1" ]; then PROXYUSER=$w; fi; done;

INSECURE=0;
for w; do if [ "$w" = "--no-check-certificate" ]; then INSECURE=1; break; fi; done;

if [ -e /var/installatron ]; then if [ ! -e /var/installatron/custompanel ]; then echo "Error: Installatron Plugin is installed on this system."; exit; fi; fi;

mkdir -p /var/installatron/custompanel

if [ -e /usr/local/installatron/bin/run ]; then
	PHP=/usr/local/installatron/bin/run
elif [ -e /usr/local/installatron/php/bin/php ]; then
	PHP=/usr/local/installatron/php/bin/php
elif [ -e /usr/local/bin/php ]; then
	PHP=/usr/local/bin/php
elif [ -e /usr/bin/php ]; then 
	PHP=/usr/bin/php
elif [ -e /usr/bin/yum ]; then
	if [ -e /home/ec2-user ]; then
		amazon-linux-extras install -y php7.2
	fi
	#RHRL: subscription-manager repos --enable rhel-server-rhscl-7-eus-rpms
	#CentOS: yum -y install centos-release-scl
	yum -y install unzip
	yum -y install php-common php-cli php-gd php-mbstring php-mysqli php-pdo php-xml php-intl php-zip php-json
	PHP=/usr/bin/php
elif [ -f /usr/bin/apt-get ]; then
	apt-get -y install unzip
	apt-get -y install php-common php-cli php-gd php-mbstring php-mysqlnd php-sqlite3 php-xml php-intl php-zip php-curl php-bcmath php-soap || apt-get -y install php5-common php5-cli php5-gd php5-mysqlnd php5-sqlite php5-intl
	PHP=/usr/bin/php
else
	echo "Error: PHP not installed (and YUM/apt-get not supported).";
	echo
	echo "Execute the commands at the following URL to install an instance of PHP, and then re-execute this script."
	echo "http://installatron.com/docs/admin/troubleshooting#nophp";
	echo
	echo "Contact Installatron Support with any questions:"
	echo "https://secure.installatron.com/contact"
	exit;
fi

if [ "$OS" = "FreeBSD" ]; then
	FETCHER="fetch -o "
else
	if [ ! -e /usr/bin/curl ] && [ ! -e /bin/curl ]; then
		if [ -e /usr/bin/yum ]; then
			yum -y install curl
		elif [ -f /usr/bin/apt-get ]; then
			apt-get -y install curl
		fi
	fi
	if [ "$PROXYUSER" != "" ]; then
		FETCHER="curl --proxy $PROXY --proxy-user $PROXYUSER -o "
	elif [ "$PROXY" != "" ]; then
		FETCHER="curl --proxy $PROXY -o "
	else
		FETCHER="curl -o "
	fi
fi

if [ ! -e /usr/local/installatron/etc/php.ini ]; then
	mkdir -p /usr/local/installatron/etc
	$FETCHER /usr/local/installatron/etc/php.ini https://$M1/php.ini
fi

$FETCHER /usr/local/installatron/etc/repair https://$M1/repair
$PHP -n -c /usr/local/installatron/etc/php.ini -q /usr/local/installatron/etc/repair -- $* 2>/dev/null || $PHP -q /usr/local/installatron/etc/repair -- $*
