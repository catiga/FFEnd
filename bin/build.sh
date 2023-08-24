#!/bin/sh

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o regulustar-linux main.go


echo set sytem env
set GOOS=linux
set GOARCH=arm64
set GOARM=7

HOME=$(cd `dirname $0`; pwd)
HOME=$(dirname "$HOME")

dir="$HOME/target"

echo $dir



if [ ! -d "$dir" ];then
mkdir $dir
else
echo "removing target folder..."
rm -rf $dir
mkdir $dir
fi

mkdir $dir/conf

cp $HOME/config/*.yaml $dir/conf
cp $HOME/bin/startup.sh $dir

echo "start building..."

CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build -o spw-linux $HOME/main.go

cp $HOME/spw-linux $dir