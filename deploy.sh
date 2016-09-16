#!/bin/bash -ex

go vet .

rm -f ./prometheus-piwik-exporter
GOOS=linux GOARCH=amd64 go build

INSTANCE_IP=`tugboat droplets | grep 'monitoring ' | egrep -oh "[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+" || true`
scp -C -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -P 2222 ./Speech-ba6281533dc8.json root@$INSTANCE_IP:/root/Speech-ba6281533dc8.json
scp -C -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -P 2222 ./prometheus-cloudfront-logs-exporter root@$INSTANCE_IP:/root/prometheus-cloudfront-logs-exporter.new

tugboat ssh -n monitoring <<EOF
set -ex

mv /root/prometheus-cloudfront-logs-exporter.new \
  /root/prometheus-cloudfront-logs-exporter
/root/prometheus-cloudfront-logs-exporter --pem_path Speech-ba6281533dc8.json

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
