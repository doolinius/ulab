#!/bin/sh

PWD=`pwd`
SHOULD_BE="/home/$USER/lab_data"
if [ "${PWD}" != "${SHOULD_BE}" ] 
then
    exit 1
else
    ls lab_data
fi
