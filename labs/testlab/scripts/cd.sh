#!/bin/sh

PWD=`pwd`
SHOULD_BE="$HOME/lab_data"
if [ "${PWD}" != "${SHOULD_BE}" ] 
then
    exit 1
else
    exit 0
fi
