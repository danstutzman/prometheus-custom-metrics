#!/bin/bash -ex
cd `dirname $0`/..

go run main.go '{"PortNum": 3000,
  "MemoryUsage": true,
  "PiwikExporter": false,
  "SecurityUpdates": false,
  "UrlToPing": {
    "Pop3CredsJson": "conf/pop3.creds.json",
    "EmailMaxAgeInMins": 60,
    "EmailSubject": "[FIRING:1] FakeAlertToVerifyEndToEnd",
    "SuccessUrl": "https://nosnch.in/480f8a1fa3"
  }
}'
