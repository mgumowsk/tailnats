# tailnats
Tail a file and publish to NATS

Usage:

    $ tailnats /var/log/test.log

Environment variables:

	CLUSTER_NAME = "test-cluster"
	NATS_SERVER" = "nats://localhost:4222"
	NATS_CLIENT_NAME" = "tailnats"
