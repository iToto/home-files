#!/bin/sh

# ./createVboxVM.sh nebu Ubuntu_64 4 /media/z500gb1/virtual-machines/hosts /media/z500gb1/virtual-machines/downloads/ubuntu64-20.04-server.iso

# VBoxManage startvm vm1 -type headless
# VBoxManage unregistervm vm1 --delete

VM_NAME=$1
OS=$2
MEMORY_GB=$3
VM_HOST_PATH=$4
ISO_PATH=$5
MB_IN_GB=1024
MEMORY_SIZE=$(( $MEMORY_GB * $MB_IN_GB ))
HDD_PATH=$VM_HOST_PATH/$VM_NAME/$VM_NAME.vdi

# echo $VM_NAME
# echo $OS
# echo $MEMORY_GB
# echo $VM_HOST_PATH
# echo $ISO_PATH
# echo $MEMORY_SIZE
# echo $HDD_PATH
# 
# exit

#Create VM
VBoxManage createvm --name $VM_NAME --ostype $OS --register --basefolder `pwd`
#Set memory and network
VBoxManage modifyvm $VM_NAME --ioapic on
VBoxManage modifyvm $VM_NAME --memory $MEMORY_SIZE --vram 128
VBoxManage modifyvm $VM_NAME --nic1 nat
#Create Disk and connect Debian Iso
VBoxManage createhd --filename $HDD_PATH --size 250000 --format VDI
VBoxManage storagectl $VM_NAME --name "SATA Controller" --add sata --controller IntelAhci
VBoxManage storageattach $VM_NAME --storagectl "SATA Controller" --port 0 --device 0 --type hdd --medium  $HDD_PATH
VBoxManage storagectl $VM_NAME --name "IDE Controller" --add ide --controller PIIX4
VBoxManage storageattach $VM_NAME --storagectl "IDE Controller" --port 1 --device 0 --type dvddrive --medium $4
VBoxManage modifyvm $VM_NAME --boot1 dvd --boot2 disk --boot3 none --boot4 none

#Enable RDP
VBoxManage modifyvm $VM_NAME --vrde on
VBoxManage modifyvm $VM_NAME --vrdemulticon on --vrdeport 10001

#Start the VM
VBoxHeadless --startvm $VM_NAME
