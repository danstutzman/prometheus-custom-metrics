#!/bin/bash -ex
cd `dirname $0`/..

go install -v
$GOPATH/bin/prometheus-custom-metrics '{
  "MemoryUsage": { "MetricsPort": 3000 },
  "PapertrailUsage": {
    "ApiTokenPath": "conf/papertrail_api_token.txt",
    "MetricsPort": 3000
  },
  "UrlToPing": {
    "MetricsPort": 3001,
    "Pop3CredsJson": "conf/pop3.creds.json",
    "EmailMaxAgeInMins": 60,
    "EmailSubject": "[FIRING:1] FakeAlertToVerifyEndToEnd",
    "SuccessUrl": "https://nosnch.in/480f8a1fa3"
  }
}'
