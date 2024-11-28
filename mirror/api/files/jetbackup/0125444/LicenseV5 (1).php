<?php


namespace JetBackup\Core\License;

defined("__JETBACKUP__");
class License
{
    private static function _handshake()
    {
        if (file_exists("/sys/class/dmi/id/product_uuid")) {
            $code = \JetBackup\Core\IO\Shell::exec("/bin/cat /sys/class/dmi/id/product_uuid", $response);
            if (!$code && isset($response[0])) {
                return md5($response[0]);
            }
        }
        if (file_exists("/var/lib/dbus/machine-id")) {
            $code = \JetBackup\Core\IO\Shell::exec("/bin/cat /var/lib/dbus/machine-id", $response);
            if (!$code && isset($response[0])) {
                return md5($response[0]);
            }
        }
        $code = \JetBackup\Core\IO\Shell::exec("/sbin/ifconfig -a | /bin/grep -i \"hwaddr\" | /bin/awk \" { print \$5 } \" | /usr/bin/head -n1", $response);
        if (!$code && isset($response[0])) {
            return md5($response[0]);
        }
        $code = \JetBackup\Core\IO\Shell::exec("/sbin/ip addr | /bin/grep -i 'ether' | /bin/awk  \" { print \$2 } \" | /usr/bin/head -n1", $response);
        if (!$code && isset($response[0])) {
            return md5($response[0]);
        }
        $code = \JetBackup\Core\IO\Shell::exec("/sbin/ifconfig -a | /bin/grep -i 'ether' | /bin/awk \" { print \$2 } \" | /usr/bin/head -n1", $response);
        if (!$code && isset($response[0])) {
            return md5($response[0]);
        }
        if (file_exists("/etc/machine-id")) {
            $code = \JetBackup\Core\IO\Shell::exec("/bin/cat /etc/machine-id", $response);
            if (!$code && isset($response[0])) {
                return md5($response[0]);
            }
        }
        return "";
    }
    private static function _addAlert($error, $force = false)
    {
        $license = \JetBackup\Core\Factory::getSettingsLicense();
        if ($force || $license->getNotifyDate()) {
            if (!$force && time() - 21600 < $license->getNotifyDate()) {
                return NULL;
            }
            \JetBackup\Core\Alert\Alert::add("License check failed", "There was a failed license check. Error: " . $error . ". Please visit the following link for more information https://docs.jetbackup.com/licensing_issue_notification.html", \JetBackup\Core\Alert\Alert::LEVEL_CRITICAL);
        }
        $license->setNotifyDate(time());
        $license->save();
    }
    public static function retrieveLocalKey()
    {
        $license = \JetBackup\Core\Factory::getSettingsLicense();
        $localkey_details = new LicenseLocalKey();
        $handshake = self::_handshake();
        try {
            if (!$handshake) {
                throw new \JetBackup\Core\Exception\LicenseException("Unable to get server handshake key");
            }
            $public_key = openssl_get_publickey("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2F105dMYw8kTC54ANpMt\nWaVKEFEoRd+D0YDRDB/EGgFbEUdkDvub6FfHcDdgfmRgjUOWap0yX2MpOCy/QA/P\nSrIuSYe4+jv5J+cW2O6WbnhVdlAQTLOiOgiUYMuFNt6+gx4DaitoxD/p39elThfa\nKsAawrrxBR0mEdn0JeE1k0l/bbhekzUubR9LReOtveXdMoVygrAQ52mQ0saOSmuB\nUn2oa4kaWeZfCA/fYW8jMEEjmrftcT3XVi7yA5Fk/xQQ3eNSLIAx2gZf8B2g0zqa\n/QSq9+6BDM/6OdG9o+cMvpSoOojx5ApcQxAwHfSiISfibJ/49kmnP7okCRI7pNvE\n3wIDAQAB\n-----END PUBLIC KEY-----");
            $handshake_encrypted = "";
            if (!openssl_public_encrypt($handshake, $handshake_encrypted, $public_key)) {
                throw new \JetBackup\Core\Exception\LicenseException("Unable to prepare handshake key for validation");
            }
            $panel_info = \JetBackup\Core\Entities\Util::getPanelInfo();
            $os_info = [];
            if (file_exists("/etc/os-release")) {
                $os_info = parse_ini_file("/etc/os-release", false, INI_SCANNER_RAW);
                if (!isset($os_info["PRETTY_NAME"])) {
                    $os_info["PRETTY_NAME"] = $os_info["NAME"] . " " . $os_info["VERSION"];
                }
            }
            $postfields = ["output" => "json", "product_id" => "58ac64be19a4643cdf582727", "handshake" => base64_encode($handshake_encrypted), "info" => ["Hostname" => gethostname(), "Panel_Name" => trim($panel_info["name"]), "Panel_Version" => trim($panel_info["version"]), "Product_Version" => trim(\JetBackup\Core\Entities\Util::getVersion()), "Product_Tier" => trim($panel_info["tier"]), "Operating_System" => sizeof($os_info) && isset($os_info["PRETTY_NAME"]) ? $os_info["PRETTY_NAME"] : "Unknown", "Kernel_Release" => php_uname("r"), "Kernel_Version" => php_uname("v")]];
            try {
                $http = new \GuzzleHttp\Client();
                $response = $http->request("GET", "https://itplic.biz/api/jetbackup?key=jetbackup", ["timeout" => 30, "form_params" => $postfields]);
                if ($response->getStatusCode() !== 200 || !($data = $response->getBody())) {
                    throw new \JetBackup\Core\Exception\LicenseException("Could not resolve host (https://check.jetlicense.com)");
                }
                $result = \JetBackup\Core\Entities\Util::jsonDecode($data, true);
                if ($result === false) {
                    throw new \JetBackup\Core\Exception\LicenseException("No valid response received from the license server");
                }
                if (!$result["success"] || !$result["data"]["localkey"]) {
                    throw new \JetBackup\Core\Exception\LicenseException("Invalid response from licensing server: " . $result["message"]);
                }
                $localKey = $result["data"]["localkey"];
                $expiryDate = $result["data"]["expires"];
                $new_localkey_details = new LicenseLocalKey($localKey);
                $license = \JetBackup\Core\Factory::getSettingsLicense();
                $license->setNotifyDate();
                $license->setLocalKey($localKey);
                $license->setExpiryDate($expiryDate);
                $license->setLastCheck(time());
                $license->setNextCheck(time() + 172800);
                $license->save();
            } catch (\GuzzleHttp\Exception\GuzzleException $e) {
                throw new \JetBackup\Core\Exception\LicenseException("Failed checking license (https://check.jetlicense.com). Error: " . $e->getMessage());
            }
        } catch (\JetBackup\Core\Exception\LicenseException $e) {
            self::_addAlert($e->getMessage());
            $license->setLocalKeyInvalid($e->getMessage(), $localkey_details);
            $license->save();
            throw $e;
        }
    }
    public static function checkLocalKey()
    {
        $localKey = new LicenseLocalKey();
        if (!$localKey->getLocalKey()) {
            throw new \JetBackup\Core\Exception\LicenseException("LocalKey is empty");
        }
        if (!($handshake = self::_handshake())) {
            throw new \JetBackup\Core\Exception\LicenseException("Unable to get machine id");
        }
        $status = $localKey->getStatus() ? $localKey->getStatus() : "Invalid";
        $description = $localKey->getDescription() ? $localKey->getDescription() : "";
        $public_key = openssl_get_publickey("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2F105dMYw8kTC54ANpMt\nWaVKEFEoRd+D0YDRDB/EGgFbEUdkDvub6FfHcDdgfmRgjUOWap0yX2MpOCy/QA/P\nSrIuSYe4+jv5J+cW2O6WbnhVdlAQTLOiOgiUYMuFNt6+gx4DaitoxD/p39elThfa\nKsAawrrxBR0mEdn0JeE1k0l/bbhekzUubR9LReOtveXdMoVygrAQ52mQ0saOSmuB\nUn2oa4kaWeZfCA/fYW8jMEEjmrftcT3XVi7yA5Fk/xQQ3eNSLIAx2gZf8B2g0zqa\n/QSq9+6BDM/6OdG9o+cMvpSoOojx5ApcQxAwHfSiISfibJ/49kmnP7okCRI7pNvE\n3wIDAQAB\n-----END PUBLIC KEY-----");
        $license = \JetBackup\Core\Factory::getSettingsLicense();
        $modulo = $license->getLastCheck() % 259200;
        if ($status == "Invalid") {
            $description = $status . "Cannot find valid license. Please visit the following link for more information https://docs.jetbackup.com/licensing_issue_notification.html";
            $status = "Invalid";
            $license->setLocalKeyInvalid($description, $localKey);
            $license->setNextCheck(0);
            $license->save();
            throw new \JetBackup\Core\Exception\LicenseException($description, $status);
        }
        self::notify();
    }
    public static function notify()
    {
        $license = \JetBackup\Core\Factory::getSettingsLicense(true);
        $localkey = new LicenseLocalKey();
        if ($localkey->getStatus() == "Active" || $localkey->getStatus() == "Cancelled") {
            $license->setExpiryNotifyDate(time());
            $license->save();
        }
        if ($localkey->getStatus() == "Trial" && !$license->getExpiryNotifyDate() && $license->getExpiryDate() && time() < $license->getExpiryDate() && $license->getExpiryDate() - 172800 < time()) {
            $license->setExpiryNotifyDate(time());
            $license->save();
            \JetBackup\Core\Alert\Alert::add("Trial license is about to expire", "Your JetBackup trial license is about to expire on " . date("d M Y h:i:s A T", $license->getExpiryDate()), \JetBackup\Core\Alert\Alert::LEVEL_WARNING);
        }
        if ($localkey->getStatus() == "TrialExpired" && $license->getExpiryDate() && $license->getExpiryDate() < time() && $license->getExpiryNotifyDate() < $license->getExpiryDate()) {
            $license->setExpiryNotifyDate(time());
            $license->save();
            \JetBackup\Core\Alert\Alert::add("Trial license expired", "Your JetBackup trial license expired on " . date("d M Y h:i:s A T", $license->getExpiryDate()), \JetBackup\Core\Alert\Alert::LEVEL_CRITICAL);
        }
    }
}

?>