#!/bin/bash -ex
cd $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics

go install
go vet .

for INSTANCE in basicruby monitoring vocabincontext; do
  fwknop -s -n $INSTANCE.danstutzman.com
  ssh root@$INSTANCE.danstutzman.com <<"EOF"
    set -ex

    id -u prometheus-custom-metrics &>/dev/null || sudo useradd prometheus-custom-metrics
    sudo mkdir -p /home/prometheus-custom-metrics
    sudo chown prometheus-custom-metrics:prometheus-custom-metrics /home/prometheus-custom-metrics
    cd /home/prometheus-custom-metrics

    if [ `uname -p` == i686 ];     then ARCH=386
    elif [ `uname -p` == x86_64 ]; then ARCH=amd64; fi

    GOROOT=/home/prometheus-custom-metrics/go1.7.3.linux-$ARCH
    if [ ! -e $GOROOT ]; then
      sudo curl https://storage.googleapis.com/golang/go1.7.3.linux-$ARCH.tar.gz >go1.7.3.linux-$ARCH.tar.gz
      chown prometheus-custom-metrics:prometheus-custom-metrics go1.7.3.linux-$ARCH.tar.gz
      sudo -u prometheus-custom-metrics tar xzf go1.7.3.linux-$ARCH.tar.gz
      sudo -u prometheus-custom-metrics mv go $GOROOT
    fi
    GOPATH=/home/prometheus-custom-metrics/gopath
    sudo -u prometheus-custom-metrics mkdir -p $GOPATH
    sudo -u prometheus-custom-metrics mkdir -p $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics
EOF
  chmod 0400 conf/*
  cp $GOPATH/bin/prometheus-custom-metrics .
  time rsync -a -e "ssh -C" -r . root@$INSTANCE.danstutzman.com:/home/prometheus-custom-metrics/gopath/src/github.com/danielstutzman/prometheus-custom-metrics --include='*.go' --include='s3.creds.ini' --include='Speech-ba6281533dc8.json' --include='pop3.creds.json' --include='papertrail_api_token.txt' --include='*/' --exclude='*' --prune-empty-dirs
  rm ./prometheus-custom-metrics
  ssh root@$INSTANCE.danstutzman.com <<"EOF"
    set -ex

    if [ `uname -p` == i686 ];     then ARCH=386
    elif [ `uname -p` == x86_64 ]; then ARCH=amd64; fi

    GOROOT=/home/prometheus-custom-metrics/go1.7.3.linux-$ARCH
    GOPATH=/home/prometheus-custom-metrics/gopath
    cd $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics
    chown -R prometheus-custom-metrics:prometheus-custom-metrics .
    time sudo -u prometheus-custom-metrics GOPATH=$GOPATH GOROOT=$GOROOT $GOROOT/bin/go install

    if [ `hostname -s` == monitoring ]; then
      tee /etc/init/prometheus-custom-metrics.conf <<EOF2
        chdir /home/prometheus-custom-metrics
        start on started mysql
        setuid prometheus-custom-metrics
        setgid prometheus-custom-metrics
        respawn
        respawn limit 2 60
        script
          ./prometheus-custom-metrics '{
            "MemoryUsage": { "MetricsPort": 9102 },

            "CloudfrontLogs": {
              "MetricsPort": 9103,
              "S3CredsPath": "conf/s3.creds.ini",
              "S3Region": "us-east-1",
              "S3BucketName": "cloudfront-logs-danstutzman",
              "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
              "GcloudProjectId": "speech-danstutzman"
            },
            "PapertrailUsage": {
              "ApiTokenPath": "conf/papertrail_api_token.txt",
              "MetricsPort": 9103
            },
            "SecurityUpdates": { "MetricsPort": 9103 },
            "UrlToPing": {
              "MetricsPort": 9103,
              "Pop3CredsJson": "conf/pop3.creds.json",
              "EmailMaxAgeInMins": 60,
              "EmailSubject": "[FIRING:1] FakeAlertToVerifyEndToEnd",
              "SuccessUrl": "https://nosnch.in/480f8a1fa3"
            },

            "BillingGcloud": {
              "MetricsPort": 9104,
              "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
              "GcloudProjectId": "speech-danstutzman",
              "GcloudDatasetName": "billing_export"
            },
            "PiwikExporter": { "MetricsPort": 9104 }
          }'
          end script
EOF2
    else
      tee /etc/init/prometheus-custom-metrics.conf <<EOF2
        chdir /home/prometheus-custom-metrics
        start on started remote_syslog
        setuid prometheus-custom-metrics
        setgid prometheus-custom-metrics
        respawn
        respawn limit 2 60
        script
          ./prometheus-custom-metrics '{
            "MemoryUsage": { "MetricsPort": 9102 },
            "SecurityUpdates": { "MetricsPort": 9103 }
          }'
        end script
EOF2
    fi

    sudo service prometheus-custom-metrics stop || true
    rm -rf /home/prometheus-custom-metrics/conf
    sudo -u prometheus-custom-metrics cp -rv $GOPATH/bin/prometheus-custom-metrics ./conf \
      /home/prometheus-custom-metrics
    sudo service prometheus-custom-metrics start
    sleep 1
    curl -f http://localhost:9102/metrics >/dev/null

    sudo ufw allow from `dig +short monitoring.danstutzman.com` to any port 9102
    sudo ufw allow from `dig +short monitoring.danstutzman.com` to any port 9103
EOF
done
