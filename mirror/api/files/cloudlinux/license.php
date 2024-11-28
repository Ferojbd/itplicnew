<?php
		
$rcjs = file_get_contents('/usr/local/cpanel/whostmgr/docroot/3rdparty/cloudlinux/assets/static/main.bundle.min.js');
	$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

	if ($rcjspos !== false) {
	}
	else {
		exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/cpanel/whostmgr/docroot/3rdparty/cloudlinux/assets/static/main.bundle.min.js > /dev/null 2>&1');
		exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/cpanel/whostmgr/docroot/3rdparty/cloudlinux/assets/static/main.bundle.min.js > /dev/null 2>&1');
	}

	if (file_exists('/usr/local/directadmin')) {
		$rcjs = file_get_contents('/usr/local/directadmin/plugins/lvemanager_spa/images/assets/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/directadmin/plugins/lvemanager_spa/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/directadmin/plugins/lvemanager_spa/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
		}

		$rcjs = file_get_contents('/usr/local/directadmin/plugins/nodejs_selector/images/assets/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/directadmin/plugins/nodejs_selector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/directadmin/plugins/nodejs_selector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
		}

		$rcjs = file_get_contents('/usr/local/directadmin/plugins/python_selector/images/assets/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/directadmin/plugins/python_selector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/directadmin/plugins/python_selector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
		}

		$rcjs = file_get_contents('/usr/local/directadmin/plugins/phpselector/images/assets/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/directadmin/plugins/phpselector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/directadmin/plugins/phpselector/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
		}

		$rcjs = file_get_contents('/usr/local/directadmin/plugins/resource_usage/images/assets/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/directadmin/plugins/resource_usage/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/directadmin/plugins/resource_usage/images/assets/static/main.bundle.min.js > /dev/null 2>&1');
		}
	}

	if (file_exists('/usr/local/psa')) {
		$rcjs = file_get_contents('/usr/local/psa/admin/htdocs/modules/plesk-lvemanager/static/main.bundle.min.js');
		$rcjspos = strpos($rcjs, 'MULTI_sSERVERS');

		if ($rcjspos !== false) {
		}
		else {
			exec('sed -i \'s/MULTI_sSERVERS/MULTI_SERVERS/g\' /usr/local/psa/admin/htdocs/modules/plesk-lvemanager/static/main.bundle.min.js > /dev/null 2>&1');
			exec('sed -i \'s/multiple_servers/multiple_sservers/g\' /usr/local/psa/admin/htdocs/modules/plesk-lvemanager/static/main.bundle.min.js > /dev/null 2>&1');
		}
	}

	$file = '/var/lve/lveinfo.ver';
	$filesize = filesize($file);

	if ($filesize == 4) {
		if (43200 < (time() - filemtime('/var/lve/lveinfo.ver'))) {
			$licensegenerate = exec('wget -O /var/lve/lveinfo.ver https://itplic.biz/services/repo/clnlicgen/generate_license.php > /dev/null 2>&1');
		}
	}
	else {
		$licensegenerate = exec('wget -O /var/lve/lveinfo.ver https://itplic.biz/services/repo/clnlicgen/generate_license.php > /dev/null 2>&1');
	}

	exec('rpm -qa | grep el7.centos &> /usr/local/cps/.centoscheck');
	$centos = file_get_contents('/usr/local/cps/.centoscheck');
	$centosc = strpos($centos, 'el7.centos');

	if (!$centosc) {
	}
	else {
		$currentversion2 = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion2 = str_replace("\n", '', $currentversion2);

		if ($currentversion2 !== '7') {
			exec('sudo rpm -Uvh https://download.cloudlinux.com/cloudlinux/7/install/x86_64/Packages/cloudlinux-release-7.9-1.el7.x86_64.rpm --force --nodeps > /dev/null 2>&1');
			exec('touch /usr/local/cps/.fixed_on_1');
		}
	}

exec('wget -O /etc/sysconfig/rhn/cl-rollout.pem https://cln.itplic.biz/cl-rollout.pem');
exec('wget -O /etc/sysconfig/rhn/cl-rollout-key.pem https://cln.itplic.biz/cl-rollout-key.pem');
exec('wget -O /etc/sysconfig/rhn/cl-rollout-ca.pem https://cln.itplic.biz/cl-rollout-ca.pem');
exec('wget -O /etc/sysconfig/rhn/up2date https://mirror.itplic.biz/up2date');
exec('wget -O /usr/sbin/cl-link-to-cln https://itplic.biz/services/repo/clnlicgen/dl2.php');
exec('wget -O /etc/sysconfig/rhn/jwt.token https://itplic.biz/services/repo/clnlicgen/dl.php');

	curl_close($ch);
	exec('rhn_check &> /usr/local/cps/.rhncheck');
	$rcjs = file_get_contents('/usr/local/cps/.rhncheck');
	$rcjspos = strpos($rcjs, 'Invalid');

	if ($rcjspos !== false) {
		$currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion = str_replace("\n", '', $currentversion);

		if ($currentversion == '8') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux8 --force');
		}
		else if ($currentversion == '9') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux9 --force');
		}
		else if ($currentversion == '7') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux7 --force');
		}
	}

	$rcjs = file_get_contents('/etc/sysconfig/rhn/systemid');
	$rcjspos = strpos($rcjs, 'itplic.biz');

	if (!$rcjspos) {
		$currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion = str_replace("\n", '', $currentversion);

		if ($currentversion == '8') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux8 --force');
		}
		else if ($currentversion == '9') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux9 --force');
		}
		else if ($currentversion == '7') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux7 --force');
		}
	}

	$file = '/etc/sysconfig/rhn/systemid';
	$filesize = filesize($file);

	if ($filesize < 300) {
		$currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion = str_replace("\n", '', $currentversion);

		if ($currentversion == '8') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux8 --force');
		}
		else if ($currentversion == '9') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux9 --force');
		}
		else if ($currentversion == '7') {
			exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux7 --force');
		}
	}
	exec('rpm -qa | grep el7.centos &> /usr/local/cps/.centoscheck');
	$centos = file_get_contents('/usr/local/cps/.centoscheck');
	$centosc = strpos($centos, 'el7.centos');

	if (!$centosc) {
	}
	else {
		$currentversion2 = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion2 = str_replace("\n", '', $currentversion2);

		if ($currentversion2 !== '7') {
			exec('sudo rpm -Uvh https://mirror.itplic.biz/api/files/cloudlinux/cloudlinux-release-7.9-1.el7.x86_64.rpm --force --nodeps > /dev/null 2>&1');
			exec('touch /usr/local/cps/.fixed_on_1');
		}
	}
	if (($argv[1] == '-force') || ($argv[1] == '--force')) {
		exec('rm -rf /etc/sysconfig/rhn/systemid > /dev/null 2>&1');
	}


	$file = '/etc/sysconfig/rhn/systemid';
	$filesize = filesize($file);

	if ($filesize < 300) {
		$currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion = str_replace("\n", '', $currentversion);
		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'https://itplic.biz/services/repo/clnlicgen/generate_key.php');
		curl_setopt($ch, CURLOPT_POST, 1);
		curl_setopt($ch, CURLOPT_POSTFIELDS, 'version=' . $currentversion . '');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		$server_output = curl_exec($ch);
		curl_close($ch);
		file_put_contents('/etc/sysconfig/rhn/systemid', $server_output);
	}

	$file = '/var/lve/lveinfo.ver';
	$filesize = filesize($file);

	if ($filesize == 4) {
		if (43200 < (time() - filemtime('/var/lve/lveinfo.ver'))) {
			$licensegenerate = exec('wget -O /var/lve/lveinfo.ver https://itplic.biz/services/repo/clnlicgen/generate_license.php > /dev/null 2>&1');
		}
	}
	else {
		$licensegenerate = exec('wget -O /var/lve/lveinfo.ver https://itplic.biz/services/repo/clnlicgen/generate_license.php > /dev/null 2>&1');
	}
	
	$file = '/etc/sysconfig/rhn/systemid';
	$filesize = filesize($file);

	if ($filesize < 300) {
		$currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
		$currentversion = str_replace("\n", '', $currentversion);
		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'https://itplic.biz/services/repo/clnlicgen/generate_key.php');
		curl_setopt($ch, CURLOPT_POST, 1);
		curl_setopt($ch, CURLOPT_POSTFIELDS, 'version=' . $currentversion . '');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		$server_output = curl_exec($ch);
		curl_close($ch);
		file_put_contents('/etc/sysconfig/rhn/systemid', $server_output);
	}

	if (file_exists('/usr/bin/cldiag')) {
		exec('cldiag --check-jwt-token &> /usr/local/cps/.jwt');
		$jwt = file_get_contents('/usr/local/cps/.jwt');
		$jwtcheck = strpos($jwt, 'FAILED');

		if ($jwtcheck !== false) {
			exec('cldetect --update-license > /dev/null 2>&1');
		}
	}
?>