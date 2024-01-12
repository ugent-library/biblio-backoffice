#!/bin/ash

# export MOCK_OIDC_PRIVATE_KEY=$(openssl genrsa)
# export MOCK_OIDC_PUBLIC_KEY=$(openssl rsa -pubout -in <(echo $MOCK_OIDC_PRIVATE_KEY))

cd /dist
mkdir -p .data
openssl genrsa -out .data/oidc.rsa
openssl rsa -in .data/oidc.rsa -pubout > .data/oidc.rsa.pub

 /dist/app server