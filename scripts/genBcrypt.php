<?php

function getRandBytes($len) {
    $bytes = mcrypt_create_iv($len, MCRYPT_DEV_URANDOM);
    // Make it obvious to readers that this function
    // will return FALSE
    return $bytes === FALSE ? FALSE : $bytes;
}

$inputString = bin2hex(getRandBytes(16));
$options = ['cost' => 11, 'salt' => mcrypt_create_iv(22, MCRYPT_DEV_URANDOM)];
$bcryptedString = password_hash($inputString, PASSWORD_BCRYPT, $options);



echo "16 byte bcrypted string: $bcryptedString\n";
echo "16 byte string: $inputString";
