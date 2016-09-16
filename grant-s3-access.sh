#!/bin/bash -ex
BUCKET_NAME=cloudfront-logs-danstutzman

if [ ! -e s3.user.json ]; then
  echo "Creating user..."
  aws iam create-user --user-name prometheus-cloudfront-logs-exporter > s3.user.json
  chmod 0400 s3.user.json
fi

if [ ! -e s3.accesskey.json ]; then
  echo "Creating access key..."
  aws iam create-access-key --user-name prometheus-cloudfront-logs-exporter > s3.accesskey.json
  chmod 0400 s3.accesskey.json
fi

if [ ! -e s3.creds.ini ]; then
  python -c "import json; key = json.load(open('s3.accesskey.json'))['AccessKey']; print '[default]\naws_access_key_id = %s\naws_secret_access_key = %s' % (key['AccessKeyId'], key['SecretAccessKey'])" > s3.creds.ini
  chmod 0400 s3.creds.ini
fi

tee policy.json <<EOF
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
    }
  ]
}
EOF

aws iam put-user-policy --user-name prometheus-cloudfront-logs-exporter \
 --policy-name can-read-cloudfront-logs \
 --policy-document file://policy.json

rm policy.json

# How to delete the user:
#   aws iam list-user-policies --user-name prometheus-cloudfront-logs-exporter
#   aws iam delete-user-policy --user-name prometheus-cloudfront-logs-exporter --policy-name can-read-cloudfront-logs
#   aws iam delete-user --user-name prometheus-cloudfront-logs-exporter
