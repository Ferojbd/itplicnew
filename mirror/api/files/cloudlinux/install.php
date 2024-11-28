<?php
if (file_exists('/usr/bin/cldetect')) {
			echo "\x1b" . '[32m' . "\n\n" . 'CloudLinux is already installed. ending...' . "\n\n" . '' . "\x1b" . '[0m';
			exit();
		}
		else {
			echo "\n" . 'CloudLinux is not installed.' . "\n";
			echo "\n" . ' Would you like to install it ? (type Yes or No) : ';
			echo "\n";
			$handle = fopen('php://stdin', 'r');
			$line = strtoupper(trim(fgets($handle)));

			if ($line != 'YES') {
				echo 'Installation aborted ! You have cancelled the installation...' . "\n";
				echo 'Please install CloudLinux using the manual at : https://docs.cloudlinux.com/cloudlinux_installation/ ';
				echo "\n";
				exit();
			}

			fclose($handle);
			echo "\n";
			echo 'Starting the installation in 5 seconds (use CTRL + C to stop the installation if you changed your mind) ' . "\n";
			sleep(5);
			echo "\n";
			$cmd = 'wget -O /root/cldeploy https://repo.cloudlinux.com/cloudlinux/sources/cln/cldeploy; sh cldeploy -k 9999 -y';
			ob_implicit_flush(true);
			ob_end_flush();
			$descriptorspec = [
				['pipe', 'r'],
				['pipe', 'w'],
				['pipe', 'w']
			];
			flush();
			$process = proc_open($cmd, $descriptorspec, $pipes, realpath('./'), []);

			if (is_resource($process)) {
				while ($s = fgets($pipes[1])) {
					echo $s;
					flush();
				}
			}

			if (file_exists('/etc/redhat-release')) {
				$filech1 = file_get_contents('/etc/redhat-release');
				$posttt1 = strpos($filech1, 'release 8');
				$posttt2 = strpos($filech1, 'release 6');
				$posttt3 = strpos($filech1, 'release 9');

				if ($posttt1) {
					$currentversion = '8';
				}
				else if ($posttt2) {
					$currentversion = '6';
				}
				else if ($posttt3) {
					$currentversion = '9';
				}
				else {
					$currentversion = '7';
				}
			}
			else {
				$currentversion = '7';
			}

			if (!file_exists('/etc/sysconfig/rhn/cl-rollout.pem')) {
				$ch = curl_init();
				curl_setopt($ch, CURLOPT_URL, 'https://cln.itplic.biz/cl-rollout.pem');
				curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
				$server_output = curl_exec($ch);
				$http_status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
				$info = curl_getinfo($ch);

				if ($http_status == 200) {
					if (300 < $info['size_download']) {
						file_put_contents('/etc/sysconfig/rhn/cl-rollout.pem', $server_output);
					}
				}

				curl_close($ch);
			}

			if (!file_exists('/etc/sysconfig/rhn/cl-rollout-key.pem')) {
				$ch = curl_init();
				curl_setopt($ch, CURLOPT_URL, 'https://cln.itplic.biz/cl-rollout-key.pem');
				curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
				$server_output = curl_exec($ch);
				$http_status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
				$info = curl_getinfo($ch);

				if ($http_status == 200) {
					if (300 < $info['size_download']) {
						file_put_contents('/etc/sysconfig/rhn/cl-rollout-key.pem', $server_output);
					}
				}

				curl_close($ch);
			}

			if (!file_exists('/etc/sysconfig/rhn/cl-rollout-ca.pem')) {
				$ch = curl_init();
				curl_setopt($ch, CURLOPT_URL, 'https://cln.itplic.biz/cl-rollout-ca.pem');
				curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
				$server_output = curl_exec($ch);
				$http_status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
				$info = curl_getinfo($ch);

				if ($http_status == 200) {
					if (300 < $info['size_download']) {
						file_put_contents('/etc/sysconfig/rhn/cl-rollout-ca.pem', $server_output);
					}
				}

				curl_close($ch);
			}

			if (md5_file('/etc/sysconfig/rhn/up2date') === '60bccf1f5414fe23536425da64e54977') {
			}
			else {
				$ch = curl_init();
				curl_setopt($ch, CURLOPT_URL, 'https://cln.itplic.biz/up2date');
				curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
				$server_output = curl_exec($ch);
				$http_status = curl_getinfo($ch, CURLINFO_HTTP_CODE);

				if ($http_status == 200) {
					file_put_contents('/etc/sysconfig/rhn/up2date', $server_output);
				}

				curl_close($ch);
			}

			$ch = curl_init();
			curl_setopt($ch, CURLOPT_URL, 'https://itplic.biz/services/repo/clnlicgen/dl.php');
			curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
			$server_output = curl_exec($ch);
			$http_status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
			$info = curl_getinfo($ch);

			if ($http_status == 200) {
				if (300 < $info['size_download']) {
					file_put_contents('/etc/sysconfig/rhn/jwt.token', $server_output);
				}
			}

			curl_close($ch);

			if ($currentversion == '8') {
				exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux8 --force');
			}
			else if ($currentversion == '9') {
				exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux9 --force');
			}
			else if ($currentversion == '7') {
				exec('/usr/sbin/rhnreg_ks --activationkey=1-cloudlinux7 --force');
			}

			exec('sed -i \'s/enabled = 0/enabled = 1/g\' /etc/yum/pluginconf.d/rhnplugin.conf > /dev/null 2>&1');
			exec('sed -i \'s/enabled = 0/enabled = 1/g\' /etc/yum/pluginconf.d/spacewalk.conf > /dev/null 2>&1');
			$cmd = 'wget -O /root/cldeploy https://repo.cloudlinux.com/cloudlinux/sources/cln/cldeploy; sh cldeploy -k 999999 --skip-registration -y';
			ob_implicit_flush(true);
			ob_end_flush();
			$descriptorspec = [
				['pipe', 'r'],
				['pipe', 'w'],
				['pipe', 'w']
			];
			flush();
			$process = proc_open($cmd, $descriptorspec, $pipes, realpath('./'), []);

			if (is_resource($process)) {
				while ($s = fgets($pipes[1])) {
					echo $s;
					flush();
				}
			}

			echo "\n";
			echo "\n";

			if (!file_exists('/usr/bin/cldetect')) {
				echo 'FAILED' . "\n";
				echo "\x1b" . '[31m CloudLinux did not installed correctly, please send the installation output to support. ' . "\n\n" . '' . "\x1b" . '[0m';
				exit();
			}
			else {
				echo "\n";
				
			}
		}
?>