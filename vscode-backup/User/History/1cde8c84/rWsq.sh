#!/bin/sh

start=1
end=10
url="https://MY-URL.com/path/to/files/"

for (( counter=start; counter<end; counter++ ))
do
download_url="${url}${counter}.jpg"
wget ${download_url}
done
