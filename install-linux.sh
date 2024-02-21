#!/bin/bash

filename=passgen-linux-amd64.tar.xz

echo "Installing for linux..."

cd $(dirname $0)
dir=/tmp/$(mktemp -u passgen-XXXX)
mkdir $dir
cp $filename $dir/
cd $dir
tar -xf $filename
result=$(make user-install 2>&1)
exitCode=$?
if [ $exitCode == 0 ]; then
	echo "Successfully Installed."
else
	echo -e "Installation failed:\n\n$result"
fi
rm -r $dir
