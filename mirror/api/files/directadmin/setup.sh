#!/bin/bash

###############################################################################
# setup.sh
# DirectAdmin  setup.sh  file  is  the  first  file  to  download  when doing a
# DirectAdmin Install.  If  you  are unable to run this script with
# ./setup.sh  then  you probably need to set it's permissions.  You can do this
# by typing the following:
#
# chmod 755 setup.sh
#
# after this has been done, you can type ./setup.sh to run the script.
#
###############################################################################

color_reset=$(printf '\033[0m')
color_green=$(printf '\033[32m')
color_red=$(printf '\033[31m')

echogreen () {
	echo "[setup.sh] ${color_green}$*${color_reset}"
}

echored () {
	echo "[setup.sh] ${color_red}$*${color_reset}"
}


if [ "$(id -u)" != "0" ]; then
	echored "You must be root to execute the script. Exiting."
	exit 1
fi

#Global variables
DA_CHANNEL=${DA_CHANNEL:="current"}
DA_PATH=/usr/local/directadmin
DACONF=${DA_PATH}/conf/directadmin.conf
DA_SCRIPTS="${DA_PATH}/scripts"

SETUP_TXT="${DA_SCRIPTS}/setup.txt"

SYSTEMDDIR=/etc/systemd/system

export DEBIAN_FRONTEND=noninteractive
export DEBCONF_NOWARNINGS=yes

case "${1}" in
	--help|help|\?|-\?|h)
		echo ""
		echo "Usage: $0 <license_key>"
		echo ""
		echo "or"
		echo ""
		echo "Usage: DA_CHANNEL=\"beta\" $0 <license_key>"
		echo ""
		echo "You may use the following environment variables to pre-define the settings:"
		echo "       DA_CHANNEL : Release channel: alpha, beta, current, stable"
		echo "         DA_EMAIL : Default email address"
		echo "    DA_ADMIN_USER : Default admin account user name"
		echo "DA_ADMIN_PASSWORD : Default admin account password"
		echo "      DA_HOSTNAME : Hostname to use for installation"
		echo "       DA_ETH_DEV : Network device"
		echo "           DA_NS1 : pre-defined ns1"
		echo "           DA_NS2 : pre-defined ns2"
		echo ""
		echo "Just set any of these environment variables to non-empty value (for example, DA_SKIP_CSF=true) to:"
		echo "            DA_SKIP_FASTEST : do not check for fastest server"
		echo "                DA_SKIP_CSF : skip installation of CFS firewall"
		echo "      DA_SKIP_MYSQL_INSTALL : skip installation of MySQL/MariaDB"
		echo "         DA_SKIP_SECURE_PHP : skip disabling insecure PHP functions automatically"
		echo "        DA_SKIP_CUSTOMBUILD : skip all the CustomBuild actions"
		echo " DA_INTERACTIVE_CUSTOMBUILD : run interactive CustomBuild installation if DA_SKIP_CUSTOMBUILD is unset"
		echo " DA_FOREGROUND_CUSTOMBUILD  : run CustomBuild installation in foreground DA_SKIP_CUSTOMBUILD is unset"
		echo ""
		echo "To customize any CustomBuild options, we suggest using environment variables: https://docs.directadmin.com/getting-started/installation/overview.html#running-the-installation-with-predefined-options"
		echo ""
		exit 0
		;;
esac

if ! command -v curl > /dev/null; then
	echogreen "Installing dependencies..."
	if [ -e /etc/debian_version ]; then
		apt-get --quiet --yes update
		apt-get --quiet --quiet --yes install curl
	else
		yum --quiet --assumeyes install curl
	fi
fi

if ! command -v curl > /dev/null; then
	echored "Please make sure 'curl' tool is available on your system and try again."
	exit 1
fi

HOST=""
if [ -n "${DA_HOSTNAME}" ]; then
	HOST="${DA_HOSTNAME}"
elif [ -s "/root/.use_hostname" ]; then
	HOST="$(head -n 1 < /root/.use_hostname)"
fi

ADMIN_USER=""
if [ -n "${DA_ADMIN_USER}" ]; then
	ADMIN_USER="${DA_ADMIN_USER}"
fi

ADMIN_PASS=""
if [ -n "${DA_ADMIN_PASSWORD}" ]; then
	ADMIN_PASS="${DA_ADMIN_PASSWORD}"
fi

EMAIL=""
if [ -n "${DA_EMAIL}" ]; then
	EMAIL="${DA_EMAIL}"
elif [ -s /root/.email.txt ]; then
	EMAIL=$(head -n 1 < /root/.email.txt)
fi

NS1=""
if [ -n "${DA_NS1}" ]; then
	NS1="${DA_NS1}"
elif [ -s /root/.ns1.txt ]; then
	NS1=$(head -n1 < /root/.ns1.txt)
fi

NS2=""
if [ -n "${DA_NS2}" ]; then
	NS2="${DA_NS2}"
elif [ -s /root/.ns2.txt ]; then
	NS2=$(head -n1 < /root/.ns2.txt)
fi

if [ $# -eq 0 ]; then
	LK=""
	until [ "${#LK}" -eq 44 ]; do
		printf "Please enter your License Key: "
		read -r LK
	done
	DA_INTERACTIVE_CUSTOMBUILD=true
elif [ "$1" = "auto" ] || [ $# -ge 4 ]; then
	if [ -e /root/.skip_get_license ]; then
		LK="skipped"
	else
		LK=$(curl --silent --location https://www.directadmin.com/clients/my_license_info.php | grep -m1 '^license_key=' | cut -d= -f2,3)
	fi
	if [ -z "${LK}" ]; then
		for ip_address in $(ip -o addr | awk '!/^[0-9]*: ?lo|link\/ether/ {print $4}' | cut -d/ -f1 | grep -v ^fe80); do {
			LK=$(curl --silent --connect-timeout 20 --interface "${ip_address}" --location https://www.directadmin.com/clients/my_license_info.php | grep -m1 '^license_key=' | cut -d= -f2,3)
			if [ -n "${LK}" ]; then
				break
			fi
		};
		done
	fi
	case "$2" in
		alpha|beta|current|stable)
			DA_CHANNEL="$2"
	esac
	if [ -z "${LK}" ]; then
		echo "Unable to detect your license key, please re-run setup.sh with LK provided as the argument."
		exit 1
	fi
	if [ $# -ge 4 ]; then
		HOST=$3
	fi
else
	LK="$1"
fi

###############################################################################
set -e

echo ""
echogreen "Welcome to DirectAdmin installer!"
echo ""
echogreen "Using these parameters for the installation:"
echo "                License Key: ${LK}"
echo "                 DA_CHANNEL: ${DA_CHANNEL}"
echo "                   DA_EMAIL: ${EMAIL}"
echo "             DA_ADMIN_USER : ${ADMIN_USER}"
echo "         DA_ADMIN_PASSWORD : ${ADMIN_PASS}"
echo "                DA_HOSTNAME: ${HOST}"
echo "                 DA_ETH_DEV: ${DA_ETH_DEV}"
echo "                     DA_NS1: ${NS1}"
echo "                     DA_NS2: ${NS2}"
echo "            DA_SKIP_FASTEST: ${DA_SKIP_FASTEST:-no}"
echo "                DA_SKIP_CSF: ${DA_SKIP_CSF:-no}"
echo "      DA_SKIP_MYSQL_INSTALL: ${DA_SKIP_MYSQL_INSTALL:-no}"
echo "         DA_SKIP_SECURE_PHP: ${DA_SKIP_SECURE_PHP:-no}"
echo "        DA_SKIP_CUSTOMBUILD: ${DA_SKIP_CUSTOMBUILD:-no}"
echo " DA_INTERACTIVE_CUSTOMBUILD: ${DA_INTERACTIVE_CUSTOMBUILD:-no}"
echo "  DA_FOREGROUND_CUSTOMBUILD: ${DA_FOREGROUND_CUSTOMBUILD:-no}"
echo ""

echogreen "Starting installation..."

if [ -e ${DACONF} ]; then
	echo ""
	echo ""
	echo "*** DirectAdmin already exists ***"
	echo "    Press Ctrl-C within the next 10 seconds to cancel the install"
	echo "    Else, wait, and the install will continue, but will destroy existing data"
	echo ""
	echo ""
	sleep 10
fi

if [ -e /usr/local/cpanel ]; then
        echo ""
        echo ""
        echo "*** CPanel exists on this system ***"
        echo "    Press Ctrl-C within the next 10 seconds to cancel the install"
        echo "    Else, wait, and the install will continue overtop (as best it can)"
        echo ""
        echo ""
        sleep 10
fi

echo "* Installing pre-install packages ....";
if [ -e "/etc/debian_version" ]; then
	apt-get --quiet --yes update || true
	apt-get -y install \
		patch diffutils perl tar zip unzip curl \
		openssl quota logrotate rsyslog zstd git \
		procps file e2fsprogs xfsprogs hostname \
		iproute2 cron ca-certificates dnsutils \
		mime-support `# provides /etc/mime.types, used by DA web service, can be replaced with media-types on Debian >= 11` \
		python3 debianutils
else
	yum -y install \
		patch diffutils perl tar zip unzip curl \
		openssl quota logrotate rsyslog zstd git \
		procps-ng file e2fsprogs xfsprogs hostname \
		iproute cronie ca-certificates bind-utils \
		mailcap `# provides /etc/mime.types, used by DA web service` \
		python3 which
fi
echo "*";
echo "*****************************************************";
echo "";

###############################################################################
###############################################################################

# We now have all information gathered, now we need to start making decisions

#######
# Ok, we're ready to go.
if [ -e "/etc/debian_version" ] && [ -e /etc/apparmor.d ]; then
	mkdir -p /etc/apparmor.d/disable
	for aa_file in /etc/apparmor.d/*; do
		if [ -f "$aa_file" ]; then
			ln -s "$aa_file" /etc/apparmor.d/disable/ 2>/dev/null || true
			if [ -x /sbin/apparmor_parser ]; then
				/sbin/apparmor_parser -R "$aa_file" 2>/dev/null || true
			fi
		fi
	done
fi

if [ -n "${DA_SKIP_MYSQL_INSTALL}" ]; then
	export mysql_inst=no
fi


###############################################################################

getLicense() {
	if [ -e /root/.skip_get_license ]; then
		echo "/root/.skip_get_license exists. Not downloading license"
		return
	fi

	mkdir -p "${DA_PATH}/conf"
	echo "$1" > "${DA_PATH}/conf/license.key"
	chmod 600 "${DA_PATH}/conf/license.key"
}

# Helper function to detect static network configs without DNS servers, Hetzner
# installer is known to create such configurations
fix_static_network_without_dns() {
	if ! command -v nmcli >/dev/null; then
		return
	fi

	local conn
	conn=$(nmcli -f NAME -m tabular -t connection show --active || true)
	if [ "$(wc -l <<< "${conn}")" -ne 1 ]; then
		# we do not support multi-iface configurations
		return
	fi
	if [ "$(nmcli -f ipv4.method -m tabular -t connection show "${conn}")" != "manual" ]; then
		# DNS will be received via DHCP
		return
	fi
	if [ -n "$(nmcli -f ipv4.dns -m tabular -t connection show "${conn}")" ]; then
		# Static DNS servers are configured we are good
		return
	fi

	# We know server has one network interface with static network
	# configuration and without any DNS servers configured. It might be
	# working now because /etc/resolv.conf is not yet touched by
	# NetowrkManager but as soon as NM reconfigures the interfaces (for
	# example afer reboot) server will become semi-non functional because
	# there are not DNS servers configured. We pro actively set Google and
	# CloudFlare DNS as a fallback.
	nmcli connection modify "${conn}" +ipv4.dns 8.8.8.8,1.1.1.1 || true
}


if mount | grep -m1 -q '^/var'; then
	echo "*** You have /var partition.  The databases, emails and logs will use this partition. *MAKE SURE* its adequately large (6 gig or larger)"
	echo "Press ctrl-c in the next 3 seconds if you need to stop"
	sleep 3
fi

if [ -e /etc/logrotate.d ]; then
	cp $DA_SCRIPTS/directadmin.rotate /etc/logrotate.d/directadmin
	chmod 644 /etc/logrotate.d/directadmin
fi

mkdir -p /var/log/httpd/domains
chmod 710 /var/log/httpd/domains
chmod 710 /var/log/httpd

ULTMP_HC=/usr/lib/tmpfiles.d/home.conf
if [ -s ${ULTMP_HC} ]; then
	#Q /home 0755 - - -
	if grep -m1 -q '^Q /home 0755 ' ${ULTMP_HC}; then
		perl -pi -e 's#^Q /home 0755 #Q /home 0711 #' ${ULTMP_HC};
	fi
fi

mkdir -p /var/www/html
chmod 755 /var/www/html

cp -f ${DA_SCRIPTS}/directadmin.service ${SYSTEMDDIR}/
cp -f ${DA_SCRIPTS}/directadmin-userd@.service ${SYSTEMDDIR}/
cp -f ${DA_SCRIPTS}/directadmin-userd@.socket ${SYSTEMDDIR}/

cp -f ${DA_SCRIPTS}/startips.service ${SYSTEMDDIR}/
chmod 644 ${SYSTEMDDIR}/startips.service

systemctl daemon-reload
systemctl enable directadmin.service
systemctl enable startips.service

${DA_SCRIPTS}/fstab.sh
${DA_SCRIPTS}/cron_deny.sh

getLicense "$LK"
fix_static_network_without_dns

cp -f ${DA_SCRIPTS}/redirect.php /var/www/html/redirect.php

OLD_ADMIN=$(grep -m 1 '^adminname=' ${SETUP_TXT} 2> /dev/null | cut -d= -f2)
if [ -n "${OLD_ADMIN}" ]; then
	if getent passwd "${OLD_ADMIN}" > /dev/null 2>&1; then
		userdel -r "${OLD_ADMIN}" 2>/dev/null
	fi
	rm -rf "${DA_PATH}/data/users/${OLD_ADMIN}"
fi

#moved here march 7, 2011
mkdir -p /etc/cron.d
cp -f ${DA_SCRIPTS}/directadmin_cron /etc/cron.d/
chmod 600 /etc/cron.d/directadmin_cron
chown root /etc/cron.d/directadmin_cron
		
#CentOS/RHEL bits
if [ ! -s /etc/debian_version ]; then
	systemctl daemon-reload
	systemctl enable crond.service
	systemctl restart crond.service
fi

${DA_PATH}/directadmin install  	 \
	"--adminname=${ADMIN_USER}" 	 \
	"--adminpass=${ADMIN_PASS}" 	 \
	"--update-channel=${DA_CHANNEL}" \
	"--email=${EMAIL}"          	 \
	"--hostname=${HOST}"        	 \
	"--network-dev=${DA_ETH_DEV}"  	 \
	"--ns1=${NS1}"              	 \
	"--ns2=${NS2}"              	 \
	|| exit 1

if ! ${DA_PATH}/directadmin build install; then
	echored "Failed to configure CustomBuild"
	exit 1
fi

echo ""
echo "System Security Tips:"
echo "  https://docs.directadmin.com/operation-system-level/securing/general.html#basic-system-security"
echo ""

if [ ! -s $DACONF ]; then
	echo "";
	echo "*********************************";
	echo "*";
	echo "* Cannot find $DACONF";
	echo "* Please see this guide:";
	echo "* https://docs.directadmin.com/directadmin/general-usage/troubleshooting-da-service.html#directadmin-not-starting-cannot-execute-binary-file";
	echo "*";
	echo "*********************************";
	exit 1;
fi

if ! systemctl restart directadmin.service; then
	echored "Failed to start directadmin service, please make sure you have a valid license"
	if [ ! -e /root/.skip_get_license ]; then
		systemctl --no-pager status directadmin.service
		exit 1
	fi
fi

${DA_PATH}/directadmin taskq --syslog || true

if [ -e /etc/aliases ]; then
	if ! grep -q diradmin /etc/aliases; then
		echo "diradmin: :blackhole:" >> /etc/aliases
	fi
fi

if [ -s ${DACONF} ]; then
	echo ""
	echo "DirectAdmin should be accessible now";
	echo "If you cannot connect to the login URL, then it is likely that a firewall is blocking port 2222. Please see:"
	echo "  https://docs.directadmin.com/directadmin/general-usage/troubleshooting-da-service.html#cannot-connect-to-da-on-port-2222"
fi

if [ -z "${DA_SKIP_CUSTOMBUILD}" ]; then
	if [ -n "${DA_INTERACTIVE_CUSTOMBUILD}" ]; then
		${DA_PATH}/directadmin build create_options
	elif [ -z "${DA_SKIP_SECURE_PHP}" ]; then
		${DA_PATH}/directadmin build set secure_php yes
	fi
	${DA_PATH}/directadmin build lego
	if ${DA_PATH}/directadmin c | grep -q '^ssl=0' && ${DA_PATH}/directadmin c | grep -q '^ssl_port=0'; then
		${DA_PATH}/scripts/letsencrypt.sh request "$(hostname)" >/dev/null 2>&1 && systemctl restart directadmin.service
	fi
	# Install CustomBuild
	if [ ! -e /root/.skip_csf ] && [ -z "${DA_SKIP_CSF}" ]; then
		${DA_PATH}/directadmin build set csf yes
	fi
	if [ -z "${DA_FOREGROUND_CUSTOMBUILD}" ]; then
		if ! resp=$(curl --fail --silent --insecure --data '{"command":["all","d"]}' "$(${DA_PATH}/directadmin api-url)/api/custombuild/run"); then
			echored "Failed to start CustomBuild installer, please run command 'da build all d' manually"
		else
			logfile=$(grep -o 'custombuild.*log' <<< "${resp}")
			echo "CustomBuild installation has started, you may check the progress using the following command: tail -f /var/log/directadmin/${logfile}"
			echogreen "You will receive a message in the DirectAdmin panel when background installation finalizes."
		fi
	else
		CB_FIRST_BUILD=yes ${DA_PATH}/directadmin build all d
	fi
fi

echo ""
echo "The following information has been set:"
echo "Admin username: $(grep ^adminname= "${SETUP_TXT}" | cut -d= -f2)"
echo "Admin password: $(grep ^adminpass= "${SETUP_TXT}" | cut -d= -f2)"
echo "Admin email: $(grep ^email= "${SETUP_TXT}" | cut -d= -f2)"
echo ""
echo ""
echo "Server Hostname: $(grep ^hostname= ${SETUP_TXT} | cut -d= -f2)"
echo ""
echogreen "To login now, follow this URL: $(${DA_PATH}/directadmin login-url)"

printf \\a
sleep 1
printf \\a
sleep 1
printf \\a

exit 0
