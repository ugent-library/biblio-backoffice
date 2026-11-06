#!/bin/bash
# Runs inside LocalStack when it's ready. The bucket persists across restarts,
# so only create it when it isn't there yet.
export AWS_ENDPOINT_URL=http://localhost:3071

awslocal s3api head-bucket --bucket biblio >/dev/null 2>&1 || awslocal s3 mb s3://biblio
