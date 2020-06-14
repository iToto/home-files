#!/bin/sh

# ./createVboxVM.sh vm1 Ubuntu_64 4

# VBoxManage startvm vm1 -type headless
# VBoxManage unregistervm vm1 --delete

VM_NAME=$1
OS=$2
MEMORY_GB=$3
MB_IN_GB=1024
MEMORY_SIZE=$(( $MEMORY_GB * $MB_IN_GB ))

HDD_PATH= $VM_HOST_PATH/$VM_NAME/$VM_NAME.vdi

echo $OS
echo $VM_NAME
echo $MEMORY_SIZE

exit

#Create VM
VBoxManage createvm --name $VM_NAME --ostype $OS --register --basefolder `pwd`
#Set memory and network
VBoxManage modifyvm $VM_NAME --ioapic on
VBoxManage modifyvm $VM_NAME --memory $MEMORY_SIZE --vram 128
VBoxManage modifyvm $VM_NAME --nic1 nat
#Create Disk and connect Debian Iso
VBoxManage createhd --filename `pwd`/$VM_NAME/$VM_NAME_DISK.vdi --size 250000 --format VDI
VBoxManage storagectl $VM_NAME --name "SATA Controller" --add sata --controller IntelAhci
VBoxManage storageattach $VM_NAME --storagectl "SATA Controller" --port 0 --device 0 --type hdd --medium  `pwd`/$VM_NAME/$VM_NAME_DISK.vdi
VBoxManage storagectl $VM_NAME --name "IDE Controller" --add ide --controller PIIX4
VBoxManage storageattach $VM_NAME --storagectl "IDE Controller" --port 1 --device 0 --type dvddrive --medium `pwd`/debian.iso
VBoxManage modifyvm $VM_NAME --boot1 dvd --boot2 disk --boot3 none --boot4 none

#Enable RDP
VBoxManage modifyvm $VM_NAME --vrde on
VBoxManage modifyvm $VM_NAME --vrdemulticon on --vrdeport 10001

#Start the VM
VBoxHeadless --startvm $VM_NAME
