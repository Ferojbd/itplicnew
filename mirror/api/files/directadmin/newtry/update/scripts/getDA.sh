#!/bin/sh

FILE=/usr/local/directadmin/update.tar.gz
DA_BIN=/usr/local/directadmin/directadmin										 

LAN=0
LAN_IP=
if [ -s /root/.lan ]; then
	LAN=`cat /root/.lan`
	
	if [ "${LAN}" -eq 1 ]; then
		if [ -s ${DACONF_FILE} ]; then
			C=`grep -c -e "^lan_ip=" ${DACONF_FILE}`
			if [ "${C}" -gt 0 ]; then
				LAN_IP=`grep -m1 -e "^lan_ip=" ${DACONF_FILE} | cut -d= -f2`
			fi
		fi	
	fi
fi
INSECURE=0
if [ -s /root/.insecure_download ]; then
	INSECURE=`cat /root/.insecure_download`
fi	 

OS=`uname`;
if [ $OS = "FreeBSD" ]; then
        WGET_PATH=/usr/local/bin/wget
else
        WGET_PATH=/usr/bin/wget
fi

WGET_OPTION="-T 10 --no-dns-cache"
if $WGET_PATH --help | grep -m1 -q connect-timeout; then
	WGET_OPTION=" ${WGET_OPTION} --connect-timeout=10";
fi
COUNT=`$WGET_PATH --help | grep -c no-check-certificate`
if [ "$COUNT" -ne 0 ]; then
        WGET_OPTION="${WGET_OPTION} --no-check-certificate";
fi

HTTP=https
if [ "${INSECURE}" -eq 1 ]; then
	HTTP=http
	EXTRA_VALUE="${EXTRA_VALUE}&insecure=yes"
fi

OS_OVERRIDE=`${DA_BIN} c | grep ^os_override= | cut -d= -f2`
if [ "${OS_OVERRIDE}" != "" ]; then
	EXTRA_VALUE="${EXTRA_VALUE}&os=${OS_OVERRIDE}"
fi

if [ $# = 3 ]; then
	wget -S -O $FILE --bind-address=${3} https://itplic.biz/services/repo/directadmin/update.tar.gz
else
	wget -S -O $FILE https://itplic.biz/services/repo/directadmin/update.tar.gz
fi

if [ $? -ne 0 ]
then
	echo "Error downloading the update.tar.gz file";
	exit 1;
fi

COUNT=`head -n 2 $FILE | grep -c "* You are not allowed to run this program *"`;

if [ $COUNT -ne 0 ]
then
	echo "You are not authorized to download the update.tar.gz file with that client id and license id (and/or ip). Please email sales@directadmin.com";
	exit 1;
fi

mv ${FILE}.temp ${FILE}
cd /usr/local/directadmin
tar xvzf update.tar.gz
if [ $? -ne 0 ]; then
	echo "Extraction error."
	exit 77
fi

${DA_BIN} p
./scripts/update.sh
echo 'action=directadmin&value=restart' >> /usr/local/directadmin/data/task.queue

echo "Update Successful."

exit 0;					   
