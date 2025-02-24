#!/bin/bash

# Check if the script is being run as root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root"
    exit 1
fi

# Suspend the VirtualBox virtual machine
echo "Suspending the virtual machine..."
VBoxManage controlvm nebu savestate

# send an email to the admin notifying them of the power failure
echo "Power failure detected. The server has been shut down." | mail -s "Power Failure Detected"


# Gracefully shutdown the server
echo "Shutting down the server..."
shutdown -h now
