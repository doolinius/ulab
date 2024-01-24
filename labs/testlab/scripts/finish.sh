#!/bin/sh

PWD=`pwd`
SHOULD_BE="$HOME"
if [ "${PWD}" != "${SHOULD_BE}" ] 
then
    exit 1
else
    if ls lab_data/
    then
        exit 1
    else
        exit 0
    fi
fi
