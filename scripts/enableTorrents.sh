#!/bin/sh

# Start VPN
osascript -e 'tell app "Private Internet Access" to activate'
# Sleep for 10 seconds to give VPN time to connect
sleep 10
# Start Transmission
osascript -e 'tell app "Transmission" to activate'
