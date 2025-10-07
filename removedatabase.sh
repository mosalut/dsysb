#!/usr/bin/bash

if [ "$1" = "" ]; then
	dir="chain_0"
else
	dir="chain_$1"
fi

dir=$dir".db"

# echo $dir is removed

rm -rf $dir 
