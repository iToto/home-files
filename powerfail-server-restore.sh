#!/bin/bash

echo "restoring nebu state..."

# Restore state from VirtualBox
VBoxManage startvm nebu --type headless

echo "nebu state restored."

echo "updating packages..."

# Update Packages
sudo apt-get update

echo "packages updated."
