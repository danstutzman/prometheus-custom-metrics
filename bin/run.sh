#!/bin/bash -ex
cd `dirname $0`/..

go run main.go '{"PortNum": 3000,
  "CloudfrontLogs": {
    "S3CredsPath": "conf/s3.creds.ini",
    "S3Region": "us-east-1",
    "S3BucketName": "cloudfront-logs-danstutzman",
    "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
    "GcloudProjectId": "speech-danstutzman"
  },
  "MemoryUsage": true,
  "PiwikExporter": false,
  "UrlToPing": ""
}'
