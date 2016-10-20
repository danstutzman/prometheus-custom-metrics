#!/bin/bash -ex
cd `dirname $0`/..
BUCKET_NAME=cloudfront-logs-danstutzman

if [ ! -e conf/s3.user.json ]; then
  echo "Creating conf/s3.user.json ..."
  aws iam create-user --user-name prometheus-cloudfront-logs-exporter \
    > conf/s3.user.json
  chmod 0400 conf/s3.user.json
fi

if [ ! -e conf/s3.accesskey.json ]; then
  echo "Creating conf/s3.accesskey.json ..."
  aws iam create-access-key --user-name prometheus-cloudfront-logs-exporter \
    > conf/s3.accesskey.json
  chmod 0400 conf/s3.accesskey.json
fi

if [ ! -e conf/s3.creds.ini ]; then
  echo "Creating conf/s3.creds.ini ..."
  python -c "import json; key = json.load(open('conf/s3.accesskey.json'))['AccessKey']; print '[default]\naws_access_key_id = %s\naws_secret_access_key = %s' % (key['AccessKeyId'], key['SecretAccessKey'])" > conf/s3.creds.ini
  chmod 0400 conf/s3.creds.ini
fi

tee conf/policy.json <<EOF
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": "arn:aws:s3:::$BUCKET_NAME/*"
    }, {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": "arn:aws:s3:::$BUCKET_NAME"
    }, {
      "Effect": "Allow",
      "Action": ["s3:DeleteObject"],
      "Resource": "arn:aws:s3:::$BUCKET_NAME/*"
    }
  ]
}
EOF

aws iam put-user-policy --user-name prometheus-cloudfront-logs-exporter \
 --policy-name can-read-cloudfront-logs \
 --policy-document file://conf/policy.json

rm conf/policy.json

# How to delete the user:
#   aws iam list-user-policies --user-name prometheus-cloudfront-logs-exporter
#   aws iam delete-user-policy --user-name prometheus-cloudfront-logs-exporter --policy-name can-read-cloudfront-logs
#   aws iam delete-user --user-name prometheus-cloudfront-logs-exporter
