#!/bin/bash

SUFFIX=$1
EXECFILE="./main"
BUCKETNAME="mybucket"
#loop for set
for var in `seq 1000`
do
 $EXECFILE key set --key key_${SUFFIX}_$var --value value$var --bucket $BUCKETNAME 
  sleep 1
done






