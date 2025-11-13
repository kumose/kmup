#!/bin/sh

set -e

if [ ! -f ./build/test-env-prepare.sh ]; then
  echo "${0} can only be executed in kmup source root directory"
  exit 1
fi

echo "change the owner of files to kmup ..."
chown -R kmup:kmup .
