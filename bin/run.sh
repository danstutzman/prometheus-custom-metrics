#!/bin/bash -ex
cd `dirname $0`/..

go install -v -race
$GOPATH/bin/prometheus-custom-metrics '{
  "BillingAws": {
    "MetricsPort": 9103,
    "Bigquery": {
      "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
      "GcloudProjectId": "speech-danstutzman",
      "DatasetName": "billing_export"
    }, "S3": {
      "CredsPath": "conf/s3.creds.ini",
      "Region": "us-east-1",
      "BucketName": "billing-danstutzman"
    }
  }
}'
