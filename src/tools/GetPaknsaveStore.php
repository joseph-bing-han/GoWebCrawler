<?php

//$data = file_get_contents('https://www.paknsaveonline.co.nz/CommonApi/Store/GetStoreList');
$data = file_get_contents('https://www.ishopnewworld.co.nz/CommonApi/Store/GetStoreList');
$json = json_decode($data, true);
$stores = $json['stores'];
foreach ($stores as $store) {
    echo("\"browser_nearest_store={\\\"UserLat\\\":\\\"{$store['latitude']}\\\",\\\"UserLng\\\":\\\"{$store['longitude']}\\\",\\\"IsSuccess\\\":true}\",\n");
}

echo("\n\nBranch Name:\n");
foreach ($stores as $store) {
    echo("\"{$store['name']} ({$store['address']})\",\n");
}
