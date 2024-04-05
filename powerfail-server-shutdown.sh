#!/bin/bash

# Check if the script is being run as root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root"
    exit 1
fi

# Suspend the VirtualBox virtual machine
echo "Suspending the virtual machine..."
VBoxManage controlvm nebu savestate

# Gracefully shutdown the server
echo "Shutting down the server..."
shutdown -h now
