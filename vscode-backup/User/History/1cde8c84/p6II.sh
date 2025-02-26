#!/bin/sh

start=7445766284
end=7445766460
url="https://pictime1eus1public.azureedge.net/pictures/26/641/26641158/bdy8cue/lowres/"

for (( counter=start; counter<end; counter++ ))
do
download_url="${url}${counter}.jpg"
wget ${download_url}
done



# https://pictime1eus1public.azureedge.net/pictures/26/641/26641158/bdy8cue/lowres/7445766284.jpg
# https://pictime1eus1public.azureedge.net/pictures/26/641/26641158/bdy8cue/lowres/7445766460.jpg

# 7445766460 - 7445766284
