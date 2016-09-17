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

echo "- targets: [ 'localhost:9102' ]" \
    >/root/prometheus_configs/prometheus-cloudfront-logs-exporter.yml
curl -X POST http://localhost:9090/-/reload

tee /etc/init/prometheus-cloudfront-logs-exporter.conf <<EOF2
chdir /root/prometheus-cloudfront-logs-exporter
start on startup
script
  ./prometheus-cloudfront-logs-exporter \
    --s3_creds_path ./s3.creds.ini \
    --s3_region us-east-1 \
    --s3_bucket_name cloudfront-logs-danstutzman \
    --gcloud_pem_path ./Speech-ba6281533dc8.json \
    --gcloud_project_id speech-danstutzman \
    --port_num 9102
end script
EOF2

sudo service prometheus-cloudfront-logs-exporter stop || true
mv /root/prometheus-cloudfront-logs-exporter/prometheus-cloudfront-logs-exporter.new \
  /root/prometheus-cloudfront-logs-exporter/prometheus-cloudfront-logs-exporter
sudo service prometheus-cloudfront-logs-exporter start
EOF
