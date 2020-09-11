#!/bin/bash

SCRIPT_PATH=`realpath $0`
PROJECT_PATH=`dirname $SCRIPT_PATH`

export DIST_NAME=aasaam
if [ -n "$BRAND_ICON" ]; then
  export DIST_NAME=$BRAND_ICON
fi

echo "Building nginx protection for $DIST_NAME..."

cd $PROJECT_PATH/web
export ASSETS_PATH=challenge/assets/
nodejs build/build.js

cd $PROJECT_PATH
rice embed-go
rm -rf nginx-protection
rm -rf nginx-protection*.tgz
go build .
tar -czf nginx-protection.$DIST_NAME.tar.gz nginx-protection
