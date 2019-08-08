#!/bin/bash

SUFFIX=$1
EXECFILE="./main"
BUCKETNAME="bucket007"


#loop for set
for var in `seq 1000`
do
 $EXECFILE key set --key key_${SUFFIX}_$var --value value$var --bucket $BUCKETNAME --addr localhost:8888
#  echo $cmd
 #`$cmd`
  sleep 1
done






