<?php

generateHashAndMD5();

function generateHashAndMD5()
{
    $options = array('cost' => 11);
    $currentTime = time();
    $md5 = md5($currentTime);
    $hash = password_hash($md5, PASSWORD_BCRYPT, $options)."\n";

    echo "MD5: $md5";
    echo "\n";
    echo "HASH: $hash";
}
