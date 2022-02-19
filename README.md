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

| app                     | description                                                           |
|-------------------------|-----------------------------------------------------------------------|
| nats-js-addconsumer     | Create a durable JS Consumer programmatically                         |
| nats-js-addsourcestream | Create a JetStream sourcing another JetStream                         |
| nats-js-addstream       | Create a JetStream programmatically                                   |
| nats-js-pub             | PUB a message to JetStream, waiting for an ingest acknowledgement     |
| nats-js-pubasync        | PUB a message to JetStream asynchronously (getting a status callback) |
| nats-js-subdds          | SUB to a JS Consumer using Dynamic Delivery Subject (DDS) mode        |


