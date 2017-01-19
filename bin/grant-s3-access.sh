#!/bin/bash -ex
cd `dirname $0`/..

if [ ! -e conf/s3.user.json ]; then
  echo "Creating conf/s3.user.json ..."
  aws iam create-user --user-name prometheus-custom-metrics \
    > conf/s3.user.json
  chmod 0400 conf/s3.user.json
fi

if [ ! -e conf/s3.accesskey.json ]; then
  echo "Creating conf/s3.accesskey.json ..."
  aws iam create-access-key --user-name prometheus-custom-metrics \
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
      "Resource": "arn:aws:s3:::cloudfront-logs-danstutzman/*"
    }, {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": "arn:aws:s3:::cloudfront-logs-danstutzman"
    }, {
      "Effect": "Allow",
      "Action": ["s3:DeleteObject"],
      "Resource": "arn:aws:s3:::cloudfront-logs-danstutzman/*"
    }
  ]
}
EOF

aws iam put-user-policy --user-name prometheus-custom-metrics \
 --policy-name can-read-cloudfront-logs \
 --policy-document file://conf/policy.json

tee conf/policy.json <<EOF
{
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject"],
      "Resource": "arn:aws:s3:::billing-danstutzman/*"
    }, {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": "arn:aws:s3:::billing-danstutzman"
    }, {
      "Effect": "Allow",
      "Action": ["s3:DeleteObject"],
      "Resource": "arn:aws:s3:::billing-danstutzman/*"
    }
  ]
}
EOF

aws iam put-user-policy --user-name prometheus-custom-metrics \
 --policy-name can-read-billing \
 --policy-document file://conf/policy.json

rm conf/policy.json

# How to delete the user:
#   aws iam list-user-policies --user-name prometheus-custom-metrics
#   aws iam delete-user-policy --user-name prometheus-custom-metrics --policy-name can-read-cloudfront-logs
#   aws iam list-access-keys --user-name prometheus-custom-metrics
#   aws iam delete-access-key --user-name prometheus-custom-metrics --access-key-id AKIAJ2OHSRS75K464FUA
#   aws iam delete-user --user-name prometheus-custom-metrics
