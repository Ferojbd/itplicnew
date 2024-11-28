<?php
// Initialize cURL session
$ch = curl_init();

// Check if cURL session initialization was successful
if ($ch === false) {
    // Handle error, e.g., by logging or displaying an error message
    die('Failed to initialize cURL session');
}

// Set cURL options
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);

$file = '/etc/sysconfig/rhn/systemid';
$filesize = filesize($file);

if ($filesize < 300) {
    $currentversion = exec('cldetect --detect-os | grep -o -E \'[0-9]+\' | head -1 | sed -e \'s/^0\\+//\'');
    $currentversion = str_replace("\n", '', $currentversion);
    curl_setopt($ch, CURLOPT_URL, 'https://itplic.biz/services/repo/clnlicgen/generate_key.php');
    curl_setopt($ch, CURLOPT_POST, 1);
    curl_setopt($ch, CURLOPT_POSTFIELDS, 'version=' . $currentversion . '');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    $server_output = curl_exec($ch);

    // Check for errors
    if ($server_output === false) {
        // Handle error, e.g., by logging or displaying an error message
        die('cURL error: ' . curl_error($ch));
    }

    // Close cURL session
    curl_close($ch);

    // Write server output to file
    file_put_contents('/etc/sysconfig/rhn/systemid', $server_output);
}

?>
