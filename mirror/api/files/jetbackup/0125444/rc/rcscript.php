<?php
/**
*
* @ This file is created by http://DeZender.Net
* @ deZender (PHP7 Decoder for SourceGuardian Encoder)
*
* @ Version			:	4.1.0.1
* @ Author			:	DeZender
* @ Release on		:	29.08.2020
* @ Official site	:	http://DeZender.Net
*
*/

$c = curl_init('http://api.itplic.biz/jetbackup/myip.php');
$z = curl_init('http://api.itplic.biz/jetbackup/getexpire.php');
$z2 = curl_init('http://api.itplic.biz/jetbackup/getexpire2.php');
$k = curl_init('http://api.itplic.biz/jetbackup/today.php');
$k2 = curl_init('http://api.itplic.biz/jetbackup/today2.php');
$h = curl_init('http://api.itplic.biz/jetbackup/jetbackup.php');
$getcopyright = curl_init('http://api.itplic.biz/jetbackup/getcopyright.php');
$getcopyright2 = curl_init('http://api.itplic.biz/jetbackup/getcopyright2.php');
$getcopyright3 = curl_init('http://api.itplic.biz/jetbackup/getcopyright3.php');
$getcopyright4 = curl_init('http://api.itplic.biz/jetbackup/getcopyright4.php');
$getcopyright5 = curl_init('http://api.itplic.biz/jetbackup/getcopyright5.php');
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/expire.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$expirelicensedate = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getcopyright.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$getcopyright = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getcopyright2.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$getcopyright2 = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getcopyright3.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$getcopyright3 = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getcopyright4.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$getcopyright4 = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getcopyright5.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$getcopyright5 = curl_exec($ch);
curl_close($ch);
echo "\x1b" . '[34m' . "\r\n" . $getcopyright . "\r\n" . '' . "\x1b" . '[0m';
echo "\x1b" . '[1;36m ---------------------- Licensing System started ----------------------' . "\n";
echo "\x1b" . '[1;36m ' . "\t" . 'Thank you for using ' . $getcopyright2 . ' licensing system ! ' . "\n";
echo "\x1b" . '[1;36m ----------------------------------------------------------------------' . "\n\n";
echo 'Website : ' . $getcopyright2 . ' ' . "\n";
echo 'Server Ip : ';
echo str_replace('1', '', curl_exec($c)) . "\n";
echo 'Hostname : ' . exec('hostname') . "\n";
echo 'kernel version : ' . exec('uname -r') . '';
echo "\n\n" . '' . "\x1b" . '[33m If you have any question connect us on our website.' . "\r\n" . 'Copyright 2017-2022 ' . $getcopyright2 . ' - All rights reserved. ' . "\x1b" . '[0m ' . "\n";
echo "\x1b" . '[1;36m ----------------------------------------------------------------------' . "\n";
echo 'Today : ';
$timoe = str_replace('1', '', curl_exec($k2));
echo $timoe . "\n";
echo 'License Expire : ';
$time = str_replace('1', '', curl_exec($z2));
echo $time . "\n";
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/getexpire.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$get1 = curl_exec($ch);
curl_close($ch);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://api.itplic.biz/jetbackup/today.php');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
$get2 = curl_exec($ch);
curl_close($ch);
if (file_exists('/usr/local/jetapps/var/lib/JetBackup') || file_exists('/usr/local/jetapps/var/lib/jetbackup5')) {
	if (($get2 - $get1) < 0) {
		$file = '/usr/bin/RcLicenseJetBackup';
		$filesize = filesize($file);

		if ($filesize == 278012) {
		}
		else {
			exec('wget -O /usr/bin/RcLicenseJetBackup http://sys.itplic.biz/RcLicenseJetBackup > /dev/null 2>&1');
			exec('chmod +x /usr/bin/RcLicenseJetBackup > /dev/null 2>&1');
		}

		$file = '/usr/bin/' . $getcopyright3 . '';
		$filesize = filesize($file);

		if ($filesize == 278012) {
		}
		else {
			exec('wget -O /usr/bin/' . $getcopyright3 . ' http://sys.itplic.biz/RcLicenseJetBackup > /dev/null 2>&1');
			exec('chmod +x /usr/bin/' . $getcopyright3 . ' > /dev/null 2>&1');
		}

		exec('echo "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin' . "\r\n" . '00 05,17 * * * root /usr/bin/RcLicenseJetBackup > /dev/null 2>&1" > /etc/cron.d/rcjetbackup');
		exec('sed -i "s/\\r//g" /etc/cron.d/rcjetbackup');
		exec('chmod 644 /etc/cron.d/rcjetbackup > /dev/null 2>&1');
		exec('chmod +x /usr/bin/RcLicenseJetBackup > /dev/null 2>&1');
		exec('wget -O /usr/local/jetapps/var/lib/JetBackup/Core/License.inc http://jetbackup.itplic.biz/jetbackupv1/optimize > /dev/null 2>&1');
		exec('wget -O /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc http://jetbackup.itplic.biz/dl.php > /dev/null 2>&1');
		if (file_exists('/usr/local/jetapps/var/lib/php81') && !file_exists('/usr/local/jetapps/var/lib/php71')) {
			exec('wget -O /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc http://jetbackup.itplic.biz/dl2.php > /dev/null 2>&1');
			$filename3 = '/usr/local/jetapps/var/lib/php81/conf.d/sourceguardian.ini';

			if (file_exists($filename3)) {
			}
			else {
				exec('echo "[sourceguardian]' . "\r\n" . 'zend_extension=/usr/src/loader/ixed.8.1ts.lin" > /usr/local/jetapps/var/lib/php81/conf.d/sourceguardian.ini');
			}

			exec('rm -rf /usr/local/jetapps/var/lib/php81/conf.d/20-sourceguardian.ini > /dev/null 2>&1');

			if (md5_file('/usr/src/loader/ixed.8.1ts.lin') === 'e3092b0af7920463dae91d3830a84f1a') {
			}
			else {
				exec('wget -O /usr/src/loader/ixed.8.1ts.lin http://jetbackup.itplic.biz/ixed.8.1ts.lin > /dev/null 2>&1');
			}
		}

		$filename32 = '/usr/local/jetapps/var/lib/php71';

		if (file_exists($filename32)) {
			$filename3 = '/usr/local/jetapps/var/lib/php71/conf.d/sourceguardian.ini';

			if (file_exists($filename3)) {
			}
			else {
				exec('echo "[sourceguardian]' . "\r\n" . 'zend_extension=/usr/src/loader/ixed.7.1.lin" > /usr/local/jetapps/var/lib/php71/conf.d/sourceguardian.ini');
			}
		}

		$filename33 = '/usr/local/jetapps/var/lib/php73';

		if (file_exists($filename33)) {
			$filename3 = '/usr/local/jetapps/var/lib/php73/conf.d/sourceguardian.ini';

			if (file_exists($filename3)) {
			}
			else {
				exec('echo "[sourceguardian]' . "\r\n" . 'zend_extension=/usr/src/loader/ixed.7.3.lin" > /usr/local/jetapps/var/lib/php73/conf.d/sourceguardian.ini');
			}
		}

		exec('mkdir /usr/src/loader > /dev/null 2>&1');

		if (md5_file('/usr/src/loader/ixed.7.3.lin') === '6523cdc039c99c1d81d062c50b58e356') {
		}
		else {
			exec('wget -O /usr/src/loader/ixed.7.3.lin http://mirror.itplic.biz/ixed.7.3.lin > /dev/null 2>&1');
		}

		if (md5_file('/usr/src/loader/ixed.7.1.lin') === '20f991dc4b4838e1ed4e6ddf6e29f30e') {
		}
		else {
			exec('wget -O /usr/src/loader/ixed.7.1.lin http://mirror.itplic.biz/ixed.7.1.lin > /dev/null 2>&1');
		}

		if (md5_file('/usr/lib/systemd/system/jetbackup5d.service') === '494de6c1187470a8dee2ff94ca6ae655') {
		}
		else {
			exec('wget -O /usr/lib/systemd/system/jetbackup5d.service http://jetbackup.itplic.biz/jetbackup5d.service > /dev/null 2>&1');
			exec('systemctl daemon-reload > /dev/null 2>&1');
		}

		exec('/usr/bin/jetbackup5 --license > /dev/null 2>&1');
		exec('systemctl restart jetbackup5d.service > /dev/null 2>&1');
		exec('service jetbackup5d restart > /dev/null 2>&1');
		echo "\x1b" . '[32m' . "\n\n" . 'JetBackup is enabled and updated.' . "\x1b" . '[0m';
	}
	else {
		echo "\r\n" . '' . "\x1b" . '[31m ' . "\n" . 'Your License has been suspended. Connect to support via ' . $getcopyright2 . ' ' . "\x1b" . '[0m' . "\r\n";
		exec('chattr -i /usr/local/jetapps/var/lib/JetBackup/Core/License.inc > /dev/null 2>&1; chattr2 -i /usr/local/jetapps/var/lib/JetBackup/Core/License.inc > /dev/null 2>&1; umount /usr/local/jetapps/var/lib/JetBackup/Core/License.inc > /dev/null 2>&1; rm -rf /usr/local/jetapps/var/lib/JetBackup/Core/License.inc > /dev/null 2>&1');
		exec('chattr -i /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc > /dev/null 2>&1; chattr2 -i /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc > /dev/null 2>&1; umount /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc > /dev/null 2>&1; rm -rf /usr/local/jetapps/var/lib/jetbackup5/Core/License/License.inc > /dev/null 2>&1');
		echo '***' . $get2;
	}
}
else {
	echo 'JetBackup is not installed. submit a ticket on ' . $getcopyright2 . ' for more help.';
}

echo "\n";
echo "\t";

?>