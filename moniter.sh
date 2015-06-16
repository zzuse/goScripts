#!/bin/sh
while true
do
sleep 10
date
if [ `ps -ef |grep media_format|grep -v grep | wc -l` -le 12 ]
then
    echo "Error"
else
	echo "Right"
	sleep 60
fi
done
