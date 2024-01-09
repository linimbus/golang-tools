#!/bin/bash

for((i=1;i<=10;i++));
do
echo run test number $i
sleep 1
value=$?
echo $value
if [ "$value" -ne "0" ]; then
   exit 0
fi
done