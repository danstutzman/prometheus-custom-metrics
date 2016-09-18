#!/bin/bash -ex
go run main.go bigquery.go s3.go collector.go \
  --s3_creds_path ./s3.creds.ini \
  --s3_region us-east-1 \
  --s3_bucket_name cloudfront-logs-danstutzman \
  --gcloud_pem_path Speech-ba6281533dc8.json \
  --gcloud_project_id speech-danstutzman \
  --port_num 3000
