#!/bin/sh

lastchange=`chage -l jdoolin | head -n 1 | awk '{print $5, $6, $7;}'`
if [ ${lastchange} != "Jan 1, 2024"] then
    exit 0
else
    exit 1
fi