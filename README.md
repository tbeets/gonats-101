# gonats-101

# The Basics

| app | description |
|-----|-----|
| nats-bench | handy core NATS benchmark tool |
| nats-echo | echo back request on REPLY subject |
| nats-pub | ye olde fire-and-forget PUB |
| nats-qsub | SUB as part of a subscriber group emulating a queue |
| nats-req | PUB an api request message w/REPLY subject |
| nats-rply | SUB as a service api and PUB a reply message |
| nats-sub | ye olde SUB interest |

# JetStream

| app                     | description                                                                              |
|-------------------------|------------------------------------------------------------------------------------------|
| nats-js-addconsumer     | Create a durable JS Consumer programmatically                                            |
| nats-js-addsourcestream | Create a JetStream sourcing another JetStream                                            |
| nats-js-addstream       | Create a JetStream programmatically                                                      |
| nats-js-pub             | PUB a message to JetStream, client waiting for an ingest acknowledgement                 |
| nats-js-pubasync        | PUB a message to JetStream asynchronously (getting an ingest acknowledgement callback)   |
| nats-js-subdds          | SUB to a Pull JS Consumer using Dynamic Delivery Subject (DDS) mode                      |
| nats-js-subsds          | SUB to a Push JS Consumer using Static Delivery Subject (SDS) mode                       |         

## Using ngs-xacct-demo for example JetStreams

See [ngs-xacct-demo](https://github.com/ConnectEverything/ngs-xacct-demo)

###### nats-js-subsdds example
```bash
# FULFILLEVENTS-C1 is a Pull JS Consumer
nats-js-subdds -s "tls://connect.ngs.synadia-test.com" -creds "/home/todd/.nkeys/creds/test-syn/todd-test-a/test-ash.creds" FULFILLEVENTS FULFILLEVENTS-C1

# FULFILLEVENTS stream subscribes to retail.v1.fulfill.>
./nats-js-pubasync -s "tls://connect.ngs.synadia-test.com" -creds "/home/todd/.nkeys/creds/test-syn/todd-test-a/test-ash.creds" "retail.v1.fulfill.completed" "Fulfill for order 1234 completed!"
```

###### nats-js-subsds example
```bash
# ORDEREVENTS-C1 is a Push JS Consumer with static delivery subject of deliver.retail.v1.order.events and deliver group of order-processor
./nats-js-subsds -s "tls://connect.ngs.synadia-test.com" -creds "/home/todd/.nkeys/creds/test-syn/todd-test-a/test-ash.creds" ORDEREVENTS ORDEREVENTS-C1

# ORDEREVENTS stream subscribes to retail.v1.order.>
./nats-js-pubasync -s "tls://connect.ngs.synadia-test.com" -creds "/home/todd/.nkeys/creds/test-syn/todd-test-a/test-ash.creds" "retail.v1.order.captured" "Captured order 1234!"
```
