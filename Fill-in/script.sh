#!/bin/bash
DIR=$(cd $(dirname $0) && pwd)
echo "DIR $DIR"
cd $DIR
./ver1

read -p "Press [Enter] key to start backup..."