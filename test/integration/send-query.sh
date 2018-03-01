#!/bin/sh
# darwest@ebay.com <darryl.west>
# 2018-03-01 07:24:59
#

curl -v -H 'Content-type: application/json' -H 'Accept: application/json' -X PUT -d '{"req":"status"}' http://127.0.0.1:3300/status
