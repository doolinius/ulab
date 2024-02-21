#!/bin/sh

part=$1
file=$2
want=$3
#data='-rw-r--r-- 1 alice developers 263 Feb 20 09:07 fib.c'

data=`ls -l ${2}`

case $part in
	perm) got=`echo ${data} | cut -d ' ' -f 1`;;
	user) got=`echo ${data} | cut -d ' ' -f 3`;;
	group) got=`echo ${data} | cut -d ' ' -f 4`;;
	pg) got=`echo ${data} | awk '{print($1 $4);}'`
	*) echo "WTF?";;
esac

if [ "${got}" = "${want}" ]
then
	exit 0
else
	exit 1
fi

