#!/bin/bash -ex
cd `dirname $0`/..

go run main.go '{"PortNum": 3000,
  "MemoryUsage": true,
  "PiwikExporter": false,
  "SecurityUpdates": false,
  "UrlToPing": ""
}'
