#!/bin/sh

PWD=`pwd`
SHOULD_BE="$HOME/datafiles"
if [ "${PWD}" != "${SHOULD_BE}" ] 
then
    exit 0
else
    exit 1
fi
