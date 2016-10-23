#!/bin/bash -ex
cd $GOPATH/src/github.com/danielstutzman/prometheus-cloudfront-logs-exporter

go build -i
rm prometheus-cloudfront-logs-exporter
go vet .

ssh -p 2222 root@monitoring.danstutzman.com <<"EOF"
set -ex
GOROOT=/root/go1.7.3.linux-amd64
if [ ! -e $GOROOT ]; then
  cd /root
  sudo curl https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz >go1.7.3.linux-amd64.tar.gz
  tar xzf go1.7.3.linux-amd64.tar.gz
  mv go $GOROOT
fi
GOPATH=/root/gopath
mkdir -p $GOPATH
mkdir -p $GOPATH/src/github.com/danielstutzman/prometheus-cloudfront-logs-exporter
EOF
time rsync -a -e "ssh -C -p 2222 -o StrictHostKeyChecking=no" -r . root@monitoring.danstutzman.com:/root/gopath/src/github.com/danielstutzman/prometheus-cloudfront-logs-exporter --include='*.go' --include='s3.creds.ini' --include='Speech-ba6281533dc8.json' --include='*/' --exclude='*' --prune-empty-dirs
ssh -p 2222 root@monitoring.danstutzman.com <<"EOF"
set -ex
GOROOT=/root/go1.7.3.linux-amd64
GOPATH=/root/gopath
cd $GOPATH/src/github.com/danielstutzman/prometheus-cloudfront-logs-exporter
time GOPATH=$GOPATH GOROOT=$GOROOT $GOROOT/bin/go build -i

tee /etc/init/prometheus-cloudfront-logs-exporter.conf <<EOF2
chdir /root/prometheus-cloudfront-logs-exporter
start on startup
respawn
respawn limit 2 60
script
  ./prometheus-cloudfront-logs-exporter '{"PortNum": 9102,
    "CloudfrontLogs": {
      "S3CredsPath": "conf/s3.creds.ini",
      "S3Region": "us-east-1",
      "S3BucketName": "cloudfront-logs-danstutzman",
      "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
      "GcloudProjectId": "speech-danstutzman"
    }
  }'
end script
EOF2

sudo service prometheus-cloudfront-logs-exporter stop || true
mkdir -p /root/prometheus-cloudfront-logs-exporter
cp -rv ./prometheus-cloudfront-logs-exporter ./conf /root/prometheus-cloudfront-logs-exporter
sudo service prometheus-cloudfront-logs-exporter start
curl -f http://localhost:9102/metrics >/dev/null
EOF
