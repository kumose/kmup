#!/bin/sh

set -e

if [ ! -f ./build/test-env-check.sh ]; then
  echo "${0} can only be executed in kmup source root directory"
  exit 1
fi


echo "check uid ..."

# the uid of kmup defined in "https://kmup.com/kmup/test-env" is 1000
kmup_uid=$(id -u kmup)
if [ "$kmup_uid" != "1000" ]; then
  echo "The uid of linux user 'kmup' is expected to be 1000, but it is $kmup_uid"
  exit 1
fi

cur_uid=$(id -u)
if [ "$cur_uid" != "0" -a "$cur_uid" != "$kmup_uid" ]; then
  echo "The uid of current linux user is expected to be 0 or $kmup_uid, but it is $cur_uid"
  exit 1
fi
