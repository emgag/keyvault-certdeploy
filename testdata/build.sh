#!/usr/bin/env bash

CERTS="c1.example.org c2.example.org"

for cn in $CERTS; do
    for k in ecdsa rsa; do
        echo "$cn:$k"
        cat $cn.cert.$k.pem chain.pem > $cn.fullchain.$k.pem
        cat $cn.fullchain.$k.pem $cn.privkey.$k.pem > $cn.fullchain.key.$k.pem
    done
done