#!/bin/bash -ex
cd `dirname $0`/..

go install -v -race
$GOPATH/bin/prometheus-custom-metrics '{
  "BillingGcloud": {
    "MetricsPort": 9103,
    "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
    "GcloudProjectId": "speech-danstutzman",
    "GcloudDatasetName": "billing_export"
  }
}'
