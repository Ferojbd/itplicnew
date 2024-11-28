<?php

function get_request($url)
{
	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $url);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	$response = curl_exec($ch);
	curl_close($ch);
	return trim($response);
}
define("ESP_API_BASE_URL", "https://itplic.biz/pubapi/v1");
define("ESP_API_KEY", "yourapi");
define("ESP_INSTALL_URL", "curl -sL https://itplic.biz/installer.sh | sudo bash -");
function do_request($uri, $method, $body  = []) {

    $headers = array(
        'Authorization: Bearer ' . ESP_API_KEY,
        'Content-Type: application/json'
    );

    $ch = curl_init();
    curl_setopt($ch,CURLOPT_URL, ESP_API_BASE_URL . $uri);
    curl_setopt($ch,CURLOPT_USERAGENT, 'EASYCONFIG_WHMCS');
    curl_setopt($ch,CURLOPT_CUSTOMREQUEST, $method);
    curl_setopt($ch,CURLOPT_RETURNTRANSFER,true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);

    if ($method == "POST") {
        $payload = json_encode( $body );
        curl_setopt($ch,CURLOPT_POSTFIELDS, $payload);

    }
    $res = curl_exec($ch);
    curl_close($ch);
//    var_dump($res);
    return json_decode($res, true);
}
$response = do_request("/licenses/username/change-ip", 'POST', ['ip' => $_SERVER['REMOTE_ADDR']]);
echo 'IP Validation Done';