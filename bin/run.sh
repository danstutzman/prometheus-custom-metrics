#!/bin/bash -ex
cd `dirname $0`/..

go run main.go bigquery.go s3.go collector.go \
  --s3_creds_path conf/s3.creds.ini \
  --s3_region us-east-1 \
  --s3_bucket_name cloudfront-logs-danstutzman \
  --gcloud_pem_path conf/Speech-ba6281533dc8.json \
  --gcloud_project_id speech-danstutzman \
  --port_num 3000
