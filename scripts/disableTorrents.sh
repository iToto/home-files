#!/bin/sh

# Kill Transmission
osascript -e 'quit app "Transmission"'
# Kill VPN
ps -ef | pgrep -f pia | xargs -L1 kill
