<?php
date_default_timezone_set("Europe/London");
$mac = isset($_GET['mac']) ? $_GET['mac'] : "";
$ip = isset($_GET['ip']) ? $_GET['ip'] : "";
$core = isset($_GET['core']) ? $_GET['core'] : "";

$serial = "+aGa-KS9u-hGtj-OYz8";
$date = new DateTime();
$rdate = $date->format('Y-m-d H:i:s');
$tnow = $date->getTimestamp();
$tsd = date('Y-m-d H:i:s', $tnow);
$mdate = $date->modify('+30 days');
$expires = $mdate->getTimestamp();
$expires_date = date('Y-m-d H:i:s', $expires);
$daysleft = ($expires - $tnow) / 3600 / 24;
    $reefpubkey = <<<EOD
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+oT6o3HT4VcknRthdU3J
mjApaHu+c+s7uRZ0s2bWRgOfArOKAxKhhvzk5KKMUQ+j/KheKCuSJitp0yj8D+Kx
YeghgP/k1Oijqv5HyAKq7jJZl7y0aEzZA1ipMY6pEMkj591rqs88vTE43neMYC4W
mKbmgxrhNXN2w85xp/fMPLnP3TUd8/Hx5D2pMouY49dYcy8OXbFJQyz8QV3aOQpT
132TVszZcqmIL7LaANPfQCseOV6O9xbVWq9w/sEca2vRy+Y5bnJnSgZH9EjQIhiu
WAcg+AyLPyho0UNUvU4LA1iHiuMF0Qj4OCt7djKz9BWQ3vIvZynp3F7MQ3jqmjQe
9wIDAQAB
-----END PUBLIC KEY-----
EOD;

    $privkey = <<<EOD
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQD6hPqjcdPhVySd
G2F1TcmaMCloe75z6zu5FnSzZtZGA58Cs4oDEqGG/OTkooxRD6P8qF4oK5ImK2nT
KPwP4rFh6CGA/+TU6KOq/kfIAqruMlmXvLRoTNkDWKkxjqkQySPn3Wuqzzy9MTje
d4xgLhaYpuaDGuE1c3bDznGn98w8uc/dNR3z8fHkPakyi5jj11hzLw5dsUlDLPxB
Xdo5ClPXfZNWzNlyqYgvstoA099AKx45Xo73FtVar3D+wRxra9HL5jlucmdKBkf0
SNAiGK5YByD4DIs/KGjRQ1S9TgsDWIeK4wXRCPg4K3t2MrP0FZDe8i9nKencXsxD
eOqaNB73AgMBAAECggEAVyxw3vEcDpy5Q+YkOqJv1bSOpCnzMvkXFifnQgo147Uc
3t7E1p7kEEnuCFU7yHVN1fxsj0PiHmAa+fyYAZsbqHsTNYVWBnRRh43mQoYTHsHs
hS2IBYdAOLbgYTtIP2wOj39wGMb2VstLA7bL5SgEeugQ7GwtE+Fy4V30FMPYkQRT
PNbloj9V5TMWgkaalWL4aLP+QMfKXJuZBfb6PyvUeUPEYdp1S8eUl9OyG+5bG7tz
5dVVQYTOcVE0UyHnWu3KnNgFRGOTzA0EfS5EYkIXL7GJeN9oIo5PH9TlFA/mj9O7
t+K1YQWgXjGqZH+kDz6eByoLjHo1ACYSo9Lez4XjAQKBgQD/a6hT1luLzp4BytjZ
O8vV1PNUjLAuMyQ9hTxRm1h8SGuBe8E9YefI3CdMEMnIB8ayRB4E8+wIT0K/kkK+
qgO0+uK8iy/asWeeDmh+wSKkujyyG2VB6vi0HSDNtPz01QYbNjfpTWoigxose4DL
0tGNfDJ4bSBC/TPEzw0KPdbddwKBgQD7FnmfQSxK71gx6ezdfoPlXLo0jG007AXa
XbN/dEhWJiRh90HRALShtWggjFytVAA/n/f0DRcFmQ0CjwhBkeMVQSZ88HOovjGo
LZ8xGWda50bVWN295V+WMNzFxdyDTuJXL5ud+Ivb8ISfYu+BspLWHYUNjdB5dc2/
Jtuhpm8qgQKBgQDH9Ba1cyT9oMWPb6YtAaPEBU5sjSrLMBwZ0Tj0ReGSgfsvRZt0
mzWhx783zBi68GN7YNoDVJUduDbv0+dObbgzMQjZQzk2QhV05aCmQjoFrQohAFNX
tEP4dKkegKZaYH3ERcClcoY4+FtAIXsllSeZVHYKUpuj9aZWVyTFNL4FZwKBgBH7
3B4x9tAvMGvyy0paA2xsJdIZtMCznv+y8mZQl9XDyZtSsF4d5NIoQhsCsqifeZ0V
Ahdy0JFQEwR55id8IX2mOvF772zIopnfGqXTofl60zH4uXkecqg5O7bWoyKshb2k
5Up9QNcx9O3NkkYB2k6Hsr3zyFjKvT/Rsq1zVEcBAoGBAN3JLhX9Zy3zCHa2SpND
zS92bnmZ4Fu2cCxexFYMQS5vPiGaGrPB+g7sL6Z3HOUI4zhk9ztTSl/oYGhKiANO
6bLRieCNHVMI14VEQ0XFGrHZd1BuKXvH6nEhfd6MLvxMPcJ51iZdem+o0aX5FMfO
ETCZGdUqibAeatm1gafjRe29
-----END PRIVATE KEY-----
EOD;

        $license_data = <<<EOF
Serial: $serial
Seq: 2
Expires: $expires
PROC: $core
Upg: $expires
TS: $tnow
F: 1
SIP: $ip
MAC: $mac

EOF;
 $readprikey = openssl_get_privatekey($privkey);
        openssl_private_encrypt($license_data, $license_encrypted, $readprikey, OPENSSL_PKCS1_PADDING);
header('Content-Type: application/license-key');
header("Content-Disposition: attachment; filename=license.key");
header("Content-length: " . strlen($license_encrypted));
die($license_encrypted);
