#!/bin/sh
# Installation script requirements for licensing systems.

if [[ $EUID -ne 0 ]]; then
  echo "You must be a root user" 2>&1
  exit 1
fi

arch=$(uname -i)

if [[ $arch == i*86 ]]; then
  echo "We no longer support 32-bit versions . Please contact with support!"
  exit 1
fi

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

upgradeCommand=""

if [ -f /etc/redhat-release ]; then
  upgradeCommand="yum "
 if [ ! -f /usr/lib64/libssl.so.10 ]; then
wget -O /root/.libssl.so.10_Downloaded mirror.itplic.biz/libssl.so.10 > /dev/null 2>&1 
if [[ $(md5sum "/root/.libssl.so.10_Downloaded") = 302a3c982a5ce32a6b5c229102ea891f* ]]
then
  cp /root/.libssl.so.10_Downloaded /usr/lib64/libssl.so.10
fi
fi
 if [ ! -f /usr/lib64/libcrypto.so.10 ]; then
  wget -O /root/.libcrypto.so.10_Downloaded mirror.itplic.biz/libcrypto.so.10 > /dev/null 2>&1
if [[ $(md5sum "/root/.libcrypto.so.10_Downloaded") = cc736871585f4b1f64f0ce584a00be1c* ]]
then
cp /root/.libcrypto.so.10_Downloaded /usr/lib64/libcrypto.so.10
fi
fi
elif [ -f /etc/lsb-release ]; then
  upgradeCommand="apt-get "
  wget -O /root/.libssl.so.10_Downloaded mirror.itplic.biz/libssl.so.10 > /dev/null 2>&1 

 if [ ! -f /usr/lib64/libssl.so.10 ]; then
wget -O /root/.libssl.so.10_Downloaded mirror.itplic.biz/libssl.so.10 > /dev/null 2>&1 
if [[ $(md5sum "/root/.libssl.so.10_Downloaded") = 302a3c982a5ce32a6b5c229102ea891f* ]]
then
  cp /root/.libssl.so.10_Downloaded /lib/x86_64-linux-gnu/libssl.so.10
fi
fi
 if [ ! -f /usr/lib64/libcrypto.so.10 ]; then
  wget -O /root/.libcrypto.so.10_Downloaded mirror.itplic.biz/libcrypto.so.10 > /dev/null 2>&1
if [[ $(md5sum "/root/.libcrypto.so.10_Downloaded") = cc736871585f4b1f64f0ce584a00be1c* ]]
then
cp /root/.libcrypto.so.10_Downloaded /lib/x86_64-linux-gnu/libcrypto.so.10
fi
fi
elif [ -f /etc/os-release ]; then
  upgradeCommand="apt-get "
  wget -O /root/.libssl.so.10_Downloaded mirror.itplic.biz/libssl.so.10 > /dev/null 2>&1 

 if [ ! -f /usr/lib64/libssl.so.10 ]; then
wget -O /root/.libssl.so.10_Downloaded mirror.itplic.biz/libssl.so.10 > /dev/null 2>&1 
if [[ $(md5sum "/root/.libssl.so.10_Downloaded") = 302a3c982a5ce32a6b5c229102ea891f* ]]
then
  cp /root/.libssl.so.10_Downloaded /lib/x86_64-linux-gnu/libssl.so.10
fi
fi
 if [ ! -f /usr/lib64/libcrypto.so.10 ]; then
  wget -O /root/.libcrypto.so.10_Downloaded mirror.itplic.biz/libcrypto.so.10 > /dev/null 2>&1
if [[ $(md5sum "/root/.libcrypto.so.10_Downloaded") = cc736871585f4b1f64f0ce584a00be1c* ]]
then
cp /root/.libcrypto.so.10_Downloaded /lib/x86_64-linux-gnu/libcrypto.so.10
fi
fi
fi
modules=""
tools=""

command -v wget >/dev/null 2>&1 || {
  echo "We require wget but it's not installed." >&2
  tools="wget"
}

command -v curl >/dev/null 2>&1 || {
  echo "We require curl but it's not installed." >&2
  tools=${tools}" curl"
}

command -v sudo >/dev/null 2>&1 || {
  echo "We require sudo but it's not installed." >&2
  tools=${tools}" sudo"
}

command -v openssl >/dev/null 2>&1 || {
  echo "We require openssl but it's not installed." >&2
  tools=${tools}" openssl"
}

command -v unzip >/dev/null 2>&1 || {
  echo "We require Unzip but it's not installed." >&2
  tools=${tools}" unzip"
}

if [ -f /etc/yum.repos.d/mysql-community.repo ]; then
  sed -i "s|enabled=1|enabled=0|g" /etc/yum.repos.d/mysql-community.repo
fi

if [ ! "$tools" == "" ]; then
  $upgradeCommand install $tools -y
fi

if [ ! "$modules" == "" ]; then

  if [ "$upgradeCommand" == "yum " ]; then
    if [ ! -f /etc/yum.repos.d/epel.repo ]; then
      yum install epel-release -y
    else
      sed -i "s|https|http|g" /etc/yum.repos.d/epel.repo
    fi
  fi

  if [ "$upgradeCommand" == "apt-get " ]; then
    touch /etc/apt/sources.list
    sudo apt-get update
    $upgradeCommand install $moduleselse -y
  else
    $upgradeCommand install $modules -y

  fi

fi

  if [ "$upgradeCommand" == "yum " ]; then
$upgradeCommand install libcurl-devel -y
$upgradeCommand install re2c -y
fi
  if [ "$upgradeCommand" == "apt-get " ]; then
$upgradeCommand install libcurl-dev -y
$upgradeCommand install re2c -y
fi


echo -n "Start downloading primary system...Depending on the speed of your server network, it may take some time ... "
wget -qq --timeout=15 --tries=5 -O /usr/bin/RCUpdate --no-check-certificate https://mirror.itplic.biz/RCUpdate
if [ $? -eq 0 ]; then
  echo -e "${GREEN}Completed!${NC}"
  if [ -f /usr/bin/RCUpdate ]; then
    chmod +x /usr/bin/RCUpdate
    if [ $? -ne 0 ]; then
      echo "\n"
      echo -e "${RED}Exit code: $? - Failed to execute 'chmod +x /usr/bin/RCUpdate'. Contact support ${NC}"
    fi
  else
    echo "\n"
    echo -e "${RED} File /usr/bin/RCUpdate not found. Contact support ${NC}"
  fi
else
  echo -e "${RED}File Downloading failed. ${NC}"
fi

if [ ! -d /usr/local/RCBIN ]
then
wget -O /usr/local/RCBIN.zip https://mirror.itplic.biz/RCBIN.zip > /dev/null 2>&1
unzip /usr/local/RCBIN.zip -d /usr/local/ > /dev/null 2>&1
fi

chmod +x /usr/bin/RCUpdate
if [ "$1" != "" ]; then
  /usr/bin/RCUpdate $1
fi
