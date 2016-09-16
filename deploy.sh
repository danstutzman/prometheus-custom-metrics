#!/bin/bash -ex

go vet .

rm -f ./prometheus-piwik-exporter.new
GOOS=linux GOARCH=amd64 go build
mv ./prometheus-cloudfront-logs-exporter ./prometheus-cloudfront-logs-exporter.new

echo "mkdir -p /root/prometheus-cloudfront-logs-exporter" | tugboat ssh -n monitoring

INSTANCE_IP=`tugboat droplets | grep 'monitoring ' | egrep -oh "[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+" || true`
scp -C -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -P 2222 \
  ./s3.creds.ini \
  ./Speech-ba6281533dc8.json \
  ./prometheus-cloudfront-logs-exporter.new \
  root@$INSTANCE_IP:/root/prometheus-cloudfront-logs-exporter

tugboat ssh -n monitoring <<EOF
set -ex

cd /root/prometheus-cloudfront-logs-exporter
mv ./prometheus-cloudfront-logs-exporter.new ./prometheus-cloudfront-logs-exporter
./prometheus-cloudfront-logs-exporter \
  --s3_creds_path ./s3.creds.ini \
  --s3_region us-east-1 \
  --s3_bucket_name cloudfront-logs-danstutzman \
  --gcloud_pem_path ./Speech-ba6281533dc8.json \
  --gcloud_project_id speech-danstutzman

#tee /etc/init/prometheus-piwik-exporter.conf <<EOF2
#start on startup
#script
#  /root/prometheus-piwik-exporter -port :9101
#end script
#EOF2
#
#sudo service prometheus-piwik-exporter stop || true
#mv /root/prometheus-piwik-exporter.new /root/prometheus-piwik-exporter
#sudo service prometheus-piwik-exporter start
EOF
