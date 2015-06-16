#!/bin/ksh
#grep keys in file1, if also found in file2, sumrize
cat $1 |while read line
do
echo $line
a=`grep $line $2`
if [[ -z $a ]]
then
	echo $line >>wrong.txt;
else
	((sum=1+sum));echo $sum;
fi
done

