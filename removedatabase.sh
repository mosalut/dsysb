#!/usr/bin/bash

if [ "$1" = "" ]; then
	dir="chain"
else
	dir="chain_$1"
fi

dir=$dir".db"

echo $dir

rm -rf $dir 
