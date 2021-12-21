## Challenges and why Kafka
1. Systems: Source and target
2. 4 source systems + 6 target systems = 24 integrations
3. Each integration - difficulties of its own:
    1. Protocol - TCP, HTTP, REST, GRPC…
    2. Data format - Binary, csv, Json, protobuf…
    3. Data schema and evolution - how data is shaped and how it may change
    4. Increased load on source systems after an integration
4. Solution - Kafka
    1. Decouple source systems and target systems.
    2. Source systems end up publishing data to Kafka
    3. Target systems source data from kafka
    4. Any source streams can be integrated with Kafka:
        1. Website events
        2. Pricing data
        3. Financial transactions
        4. User interactions, others
    5. Any target system can source data from Kafka:
        1. Database
        2. Analytics
        3. Email system
        4. Audit
5. By LinkedIn, now open source - maintained by Confluent
6. Distributed, resilient and fault tolerant
7. Horizontal scalability
    1. 100s of brokers
    2. Millions of messages/sec
8. High performant - latency less than 10 ms
9. Use cases:
    1. Messaging system
    2. Activity tracking
    3. Stream processing
    4. Gathering metrics
    5. App log gathering
    6. De-coupling system dependencies
    7. Integration with Spark, link, Storm, Hadoop, …
10. Real uses of Kafka -
    1. Netflix - recommendation engine
    2. Uber - gather data user, taxi and trip data - forecast demand, compute surge pricing
    3. LinkedIn - prevent spam, collect user interactions..
11. Kafka is best at transportation mechanism

## Topics, partitions and offsets
1. Topic - A particular stream of data
2. Topic - Similar to table in DB
3. As many topics can be created
4. Topic is identified by its name
5. Topics are split by partitions
6. Each partition will be ordered
7. Each message within a partition gets an incremental id and called offset, from 0 to …
8. Creation of topic requires to mention partitions - but it can be changed later
9. Message Address -> Topic.Partition.Offset
10. Offset only have a meaning for a specific partition
11. Order is guaranteed only within a partition
12. Data is kept in Kafka only for 1 week - default 1 week
13. Offset is always incremental - never go back to 0
14. Once data written to partition - can’t be changed - immutable
15. Data is sent randomly to partition - unless key is provided

## Broker and Topics
1. Kafka cluster - composed of multiple brokers [servers]
2. Each brokers has an ID [a number]
3. Each broker - certain partition - because kafka is distributed
4. After connecting to broker, you will be connected to entire cluster
5. Good no. of brokers - 3 to 100
6. If you create a topic with 3 partitions, each partition will reside in a different broker

## Replication factor
1. Usually 2/3 - 3 is old standard
2. If a broker is down another broker can serve the data
3. Example topic-A has 2 partitions and replication factor of 2 [given 3 brokers]
    1. Broker-1: Partition-0.Topic-A
    2. Broker-2: Partition-1.Topic-A, Partition-0.Topic-A
    3. Broker-3: Partition-1.Topic-A
4. Replications are linked
5. Leader of partition
    1. One broker can be leader of a partition at a time
    2. That leader receive and serve data for a partition
    3. Other brokers are just passive replicas and they just synchronise data
    4. 1 leader and multiple ISR [in-sync replica]
    5. If leader broker goes down, other take over

## Producers
1. Write data to topics
2. Producer know which broker/partition to write to
3. In case of broker failures, producers will automatically recover
4. Producers can choose to receive acknowledgment of data writes:
    1. acks=0, Producer won’t wait for confirmation - possible data loss
        1. Useful where potentially losing data is okay - metrics collection or logs collection
        2. Performance is a priority over data loss
    2. acks=1, Producer will wait for leader to acknowledge [limited data loss]
        1. Leader responds with ack and no guarantee from replications
        2. Potential data loss if leader goes down before replication has taken place
        3. default
    3. acks=all, Leader and replicas acknowledge [no data loss]
        1. Downside: Adds a bit of latency
        2. Upside: No data loss given enough replicas, basically safe
        3. min.insync.replica
            1. broker level and topic level
            2. topic level can override
            3. At least that many no.s of brokers including the leader must ack
            4. If replication.factor = 3, min.insync=2 and acks=all -> only 1 broker going down can be tolerated, else exception [NOT_ENOUGH_REPLICA] will be received from producer
            5. Producer retries at configured no. of times
5. Retries
    1. Transient failures:
        1. NotEnoughReplicasException
    2. Default:
        1. For <2.0 of kafka - retries = 0
        2. For >=2.1 of kafka - retries = 214783647
    3. retry.backoff.ms - retrial interval - default 100ms
    4. delivery.timeout.ms = 120000ms or 2 min - producer will retry for 2 min
    5. Messages will be sent out of order
    6. Key based ordering will be an issue in case of retries
        1. max.in.flight.requests.per.connection - no of parallel requests - default = 5
        2. Set to 1 to ensure ordering - will impact throughtput
6. Idempotent producer
    1. Duplicate messages produced by producers because of n/w issues
    2. Flow
        1. Good request -> produce - commit - ack
        2. Duplicate request -> produce - commit - ack [failed to reach producer] - retry produce [since retry > 1] - duplicate commit - ack
        3. Solution -> produce - commit - ack [failed to reach producer] - retry produce - doesn’t commit [identifies the duplicate commit by an id] - ack
    3. Default settings
        1. Retries = 2^31 - 1
        2. max.in.flights.requests = 1 [0.11] and 5 [>= 1.0 or higher]
        3. acks = all
        4. enable.idempotence = true
7. Safe producer
    1. For kafka version < 0.11
        1. acks = all -> ensures data replication and acks received
        2. min.insync.replicas = 2 -> broker/topic level -> ensures 2 brokers in ISR at least have the data in ack
        3. retries = MAX_INT -> producer level -> ensures transient errors are retried indefinitely
        4. max.in.flights.requests.per.second = 1 -> producer level -> ensures only 1 request is tried at any time, preventing message reordering in case of retries
    2. For kafka version >= 0.11
        1. enable.idempotence=true (producer level)
        2. min.insync.replicas=2 (broker/topic level)
    3. Impacts latency and throughput, test as per use case to be on safer side
8. Message keys
    1. Producer can choose to send keys
    2. If key not sent, data sent in round robin to partitions
    3. If key is sent, all messages with that key will go to a particular partition
    4. Key can be anything [a string, a number…]
    5. Key is sent if message ordering is required for a specific field
9. Compression
   1. Producer level
   2. Works great on batch
   3. Upto 4x - 10MB to 2.5MB - less latency and better throughput - better disk utilisation in kafka - Disadvantages: CPU cycles at kafka to compress data. CPU cycles at consumer to decompress data
   4. compression.type
       1. none - default
       2. gzip
       3. lz4
       4. snappy

## Consumers
1. Read data from topics
2. They know which broker to read from
3. In case of broker failures, they recover
4. Data is read in order within partition

## Consumer groups
1. Each consumer in a group read data from exclusive partitions
2. More consumers than partitions - some consumers needs to be inactive
3. Example:
    1. Consumer-Grp-1
        1. Consumer-1: Topic-A-Partition-1
        2. Consumer-2: Topic-A-Partition-2, Topic-A-Partition-3
    2. Consumer-Grp-2
        1. Consumer-1: Topic-A-Partition-1
        2. Consumer-2: Topic-A-Partition-2
        3. Consumer-3: Topic-A-Partition-3
    3. Consumer-Grp-3
        1. Consumer-1: Topic-A-Partition-1, Topic-A-Partition-2, Topic-A-Partition-3
4. Consumers will automatically use a group-coordinator and consumer-coordinator to assign consumers to a partition

## Consumer offsets
1. Kafka stores the offsets at which a consumer group has been reading - like bookmarking
2. These offsets are committed live in a kafka topic __consumer_offsets
3. When a consumer group has processed data received from kafka, it would be committing the offsets
4. If consumer dies, and comes live up, so that it can read from the offset it stopped reading
5. Delivery semantics:
    1. Consumers choose when to commit offsets
        1. At most once: offsets are committed as soon as message is received, if processing goes wrong the data is lost
        2. At least once: offsets committed after message is processed, Make sure message processing is idempotent, because duplication can happen
        3. Exactly once: Only from kafka to kafka

## Kafka broker discovery:
1. Every kafka broker - bootstrap server
2. Consumer and Producer know which broker-partition to send/read data to/from
3. Each broker knows about all brokers, topics and partitions [metadat]
4. Client -> connection + metadata request. Broker -> list of brokers. Client -> connect to needed brokers
5. This implementation is already done

## Zookeeper
1. Manages brokers - keeps the list of them
2. Helps performing in leader election
3. They send notifications to kafka about changes like new topic, broker dies, broker comes up, delete topics…
4. Zookeeper by design operates in odd numbers, like 3,5,7.. servers
5. Zookeeper has a leader [handle writes], the rest of servers are followers [handler reads]
6. Zookeeper does not store consumer offsets since kafka v0.10
7. Broker connected to diff Zookeepers

## Kafka guarantees:
1. Messages are appended to topic-partition in the order they are sent
2. Consumers read messages in order stored in topic-partition
3. RepFactor of N - producer and consumer tolerate upto N-1 brokers being down
4. As long as no. of partitions are constant, the same key will go to same partition

## Start Kafka
1. Download kafka binaries online - https://www.apache.org/dyn/closer.cgi?path=/kafka/3.0.0/kafka_2.13-3.0.0.tgz - https://dlcdn.apache.org/kafka/3.0.0/kafka_2.13-3.0.0.tgz
3. Unzip the binaries, cd into kafka-v… directory, command in terminal - bin/topics.sh - to get some output and check if java is running in system
4. Steps to install java:
    1. brew tap caskroom/versions - got this: caskroom/versions was moved. Tap homebrew/cask-versions instead
    2. brew tap homebrew/cask-versions
    3. brew cask install java8 - for caskroom and brew install --cask java8 for homebrew
5. Update PATH, export Path=$PATH:<path-to-kafka-bin>, in .bash_profile/.zshrc
6. Brew install kafka - another way to install kafka, to run Kafka commands w/o .sh
7. Start zookeeper
    1. Create a directory data/zookeeper
    2. Update config/zookeeper.properties -> data:<path-to-above-data/zookeeeper-dir>
    3. Zookeeper-server-start.sh config/zookeeper.properties
    4. Keep it running in a separate terminal window
8. Create a directory data/kafka
9. Update config/server.properties -> log.dirs: <path-to-above-data/kafka-dir>
10. kafka-server-start.sh config/server.properties

## CLI [using —bootstrap-server instead of —zookeeper and port 9092 is used instead of 2181 in new kafka version]
1. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --create - will throw error because we missed to provide partitions while creations
2. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --create --partitions 3 - will throw error because we missed to provide replication factor
3. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --create --partitions 3 --replication-factor 2 - will throw error because we provided replication factor [2] greater than no. of broker [1]
4. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --create --partitions 3 --replication-factor 1 - successfully creates topic
    1. WARNING: Due to limitations in metric names, topics with a period ('.') or underscore ('_') could collide. To avoid issues it is best to use either, but not both. Created topic first_topic.
5. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --list - shows topics
6. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --describe - details about topic
7. kafka-topics.sh --bootstrap-server 127.0.0.1:9092  --topic first_topic --delete - will remove delete.topic.enable = true else will have no impact

## Kafka console Producer
1. Required args: --broker-list, —topic
2. kafka-console-producer.sh --broker-list 127.0.0.1:9092 —topic first_topic
3. kafka-console-producer.sh --broker-list 127.0.0.1:9092 —topic first_topic --producer-property acks=all
4. kafka-console-producer.sh --broker-list 127.0.0.1:9092 —topic new_topic -> if send messages to a topic that does not exist:
    1. First message sent -> WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 3 : {new_topic=LEADER_NOT_AVAILABLE} (org.apache.kafka.clients.NetworkClient) -> the topic gets created but leader is not elected, hence the warning
    2. Subsequent messages -> it produces successfully
    3. The new topic get created with default properties
    4. Defaults can be set at server.properties -> num.partitions
5. kafka-console-producer.sh -broker-list 127.0.0.1:9092 -topic new_topic —property parse.key=true —property key.separator=,

## Kafka console Consumer
1. Required args: —bootstrap-server, —topic
2. kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic first_topic > will not broadcast messages that are already present in the kafka, will broadcast the ones it receives after it starts
3. kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic first_topic --from-beginning -> gives all messages [read/unread]
4. kafka-console-consumer.sh —bootstrap-server 127.0.0.1:9092 —topic new_topic --from-beginning —property parse.key=true —property key.separator=,

## Kafka consumers in a group
1. kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic first_topic --group my-first-app
2. kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic first_topic --group my-second-app —from-beginning
    1. 1st time -> reads all message
    2. 2nd time -> since group was specified, the offsets are committed and the previous messages won’t show up
3. kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic first_topic --group my-second-app -> if messages were already there in kafka while no consumer, the are read

## Kafka consumers group in cli
1. kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list
2. kafka-consumer-groups.sh --bootstrap-server localhost:9092 —describe —group my-first-app -> CURRENT-OFFSET, LOG_END-OFFSET, LAG
3. Set watch on above command to monitor which consumer gets the data

## Resetting offsets
1. For a defined consumer group from where should it start reading messages
2. —to-datetime, —by-period, —to-earliest, —to-latest, —shift-by, —from-file, —to-current
3. kafka-consumer-group —bootstrap-server localhost:9092 —group my-first-app —reset-offsets —to-earliest —execute —topic first_topic -> NEW_OFFSET -> 0
4. kafka-consumer-group —bootstrap-server localhost:9092 —group my-first-app —reset-offsets —shift-by -2 —execute —topic first_topic -> NEW_OFFSET -> x-2 -> all consumers are shifted by 2 [LAG becomes 2]

KafkaCat/Conduktor can be replacement of KafkaCLI