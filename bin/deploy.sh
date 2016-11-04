#!/bin/bash -ex
cd $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics

go build -i
rm prometheus-custom-metrics
go vet .

for INSTANCE in basicruby monitoring vocabincontext; do
  ssh -p 2222 -o StrictHostKeyChecking=no root@$INSTANCE.danstutzman.com <<"EOF"
    set -ex

    if [ `uname -p` == i686 ];     then ARCH=386
    elif [ `uname -p` == x86_64 ]; then ARCH=amd64; fi

    GOROOT=/root/go1.7.3.linux-$ARCH
    if [ ! -e $GOROOT ]; then
      cd /root
      sudo curl https://storage.googleapis.com/golang/go1.7.3.linux-$ARCH.tar.gz >go1.7.3.linux-$ARCH.tar.gz
      tar xzf go1.7.3.linux-$ARCH.tar.gz
      mv go $GOROOT
    fi
    GOPATH=/root/gopath
    mkdir -p $GOPATH
    mkdir -p $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics
EOF
  time rsync -a -e "ssh -C -p 2222 -o StrictHostKeyChecking=no" -r . root@$INSTANCE.danstutzman.com:/root/gopath/src/github.com/danielstutzman/prometheus-custom-metrics --include='*.go' --include='s3.creds.ini' --include='Speech-ba6281533dc8.json' --include='*/' --exclude='*' --prune-empty-dirs
  ssh -p 2222 root@$INSTANCE.danstutzman.com <<"EOF"
    set -ex

    if [ -e /root/go1.7.3.linux-amd64 ]; then GOROOT=/root/go1.7.3.linux-amd64
    else                                      GOROOT=/root/go1.7.3.linux-386; fi

    GOPATH=/root/gopath
    cd $GOPATH/src/github.com/danielstutzman/prometheus-custom-metrics
    time GOPATH=$GOPATH GOROOT=$GOROOT $GOROOT/bin/go build -i

    tee /etc/init/prometheus-custom-metrics.conf <<EOF2
    chdir /root/prometheus-custom-metrics
    start on startup
    respawn
    respawn limit 2 60
    script
EOF2
    if [ `hostname -s` == monitoring ]; then
      tee -a /etc/init/prometheus-custom-metrics.conf <<EOF2
        ./prometheus-custom-metrics '{"PortNum": 9102,
          "CloudfrontLogs": {
            "S3CredsPath": "conf/s3.creds.ini",
            "S3Region": "us-east-1",
            "S3BucketName": "cloudfront-logs-danstutzman",
            "GcloudPemPath": "conf/Speech-ba6281533dc8.json",
            "GcloudProjectId": "speech-danstutzman"
          },
          "MemoryUsage": true,
          "PiwikExporter": true,
          "SecurityUpdates": true,
          "UrlToPing": "https://nosnch.in/480f8a1fa3"
        }'
        end script
EOF2
    else
      tee -a /etc/init/prometheus-custom-metrics.conf <<EOF2
        ./prometheus-custom-metrics '{"PortNum": 9102,
          "MemoryUsage": true,
          "SecurityUpdates": true}'
        end script
EOF2
    fi

    sudo service prometheus-custom-metrics stop || true
    mkdir -p /root/prometheus-custom-metrics
    cp -rv ./prometheus-custom-metrics ./conf /root/prometheus-custom-metrics
    sudo service prometheus-custom-metrics start
    curl -f http://localhost:9102/metrics >/dev/null

    sudo ufw allow from `dig +short monitoring.danstutzman.com` to any port 9102
EOF
done
