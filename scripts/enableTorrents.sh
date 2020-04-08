#!/bin/sh

# Start VPN
osascript -e 'launch app "Private Internet Access"'
# Sleep for 10 seconds to give VPN time to connect
sleep 10
# Start Transmission
osascript -e 'launch app "Transmission"'
open -a "Transmission"
# Start Sonarr
open -a Sonarr
# open NZBGet
open -a nzbget
# open Radarr
open -a Radarr
# mount volumes
osascript -e 'launch app "mount-zion-discs"'
