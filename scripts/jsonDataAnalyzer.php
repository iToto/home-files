<?php

if (! count($argv == 2)) {
    die("Usage: jsonDataAnalyzer.php file_name.json");
}
$file = $argv[1];

if (!file_exists($file)) {
    die("Error: File Not Found");
}

$json = file_get_contents($file);

if (!$data = json_decode($json, true)) {
    die("Invalid Json Found");
}

$common = [];
$totalDocuments = count($data);

foreach ($data as $index => $row) {
    if ($index == 0) {
        $common = array_keys($row);
        continue;
    }

    $common = array_intersect($common, array_keys($row));
}

echo "Total Lines: $totalDocuments \n";
echo "Common Keys: \n";
var_dump($common);
