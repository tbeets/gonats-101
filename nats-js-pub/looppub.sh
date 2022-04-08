#!/bin/bash

while :
do
   ./nats-js-pub -s nats://vbox1.tinghus.net:4222 -creds "/home/todd/lab/nats-cluster1/vault/.nkeys/creds/NatsOp/AcctA/UserA1.creds" foo "Message {{Count}}"
   sleep 2 
done

