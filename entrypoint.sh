#!/usr/bin/env bash

# Replace the logo of reporter if vendor is set
if [ -n "$VENDOR" ];then
    echo "------ copying images from mounted volume ------"
    rm -v /tmp/logo.png
    cp -v /tmp/mo-vendor/assets/images/logos/logo.png /tmp/
fi

# Starting grafana reporter
./usr/local/bin/grafana-reporter