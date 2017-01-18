#!/bin/bash -ex
cd `dirname $0`/..

go install -v -race
$GOPATH/bin/prometheus-custom-metrics '{
  "Collectors": {
    "BillingGcloud": {
      "MetricsPort": 9103,
      "Bigquery": {
        "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
        "GcloudProjectId": "speech-danstutzman",
        "DatasetName": "billing_export"
      }
    }
  }
}'
