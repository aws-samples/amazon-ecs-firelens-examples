# FireLens Example: Sending Custom EMF Metrics to CloudWatch

CloudWatch supports the Embedded Metrics Format which allows you to create json events sent with PutLogEvents that are processed and become CloudWatch Metrics. The PutLogEvents call must include an HTTP header noting that all the log events are EMF; Fluent Bit can be configured to send this header. This tutorial will show you how to collect and send custom EMF metrics with Fluent Bit. Fluent Bit can be used a replacement for the CloudWatch Agent to upload EMF. 

You can instrument your application to produce EMF metrics by manually constructing the JSON or by using one of the AWS EMF libraries (search GitHub for the most up to date list of available libraries):
* [aws-embedded-metrics-java](https://github.com/awslabs/aws-embedded-metrics-java)
* [aws-embedded-metrics-dotnet](https://github.com/awslabs/aws-embedded-metrics-dotnet)
* [aws-embedded-metrics-python](https://github.com/awslabs/aws-embedded-metrics-python)
* [aws-embedded-metrics-node](https://github.com/awslabs/aws-embedded-metrics-node)

Once your application is instrumented to produce EMF metrics, you then need to ingest them into Fluent Bit using an [input plugin](https://docs.fluentbit.io/manual/pipeline/inputs). There are two commonly options for ingesting EMF, which are outlined in the tutorials below. 

Once the EMF records are ingested into Fluent Bit, you simply need to send them to a cloudwatch output that has `log_format json/emf` set its options. This tells Fluent Bit to set the EMF HTTP header when it calls PutLogEvents. Please note that since all records will be interpreted as EMF, the logs sent by the output with the `log_format json/emf` option configured should probably be all EMF events. This is not a hard requirement, however, the best practice is to split EMF events and other events to separate output instances matching different tags and sending to different log streams. 

Both the [older golang CloudWatch plugin](https://github.com/aws/amazon-cloudwatch-logs-for-fluent-bit#new-higher-performance-core-fluent-bit-plugin) and the [newer high performance plugin](https://docs.fluentbit.io/manual/pipeline/outputs/cloudwatch) support this option. 

## Tutorial 1: Send EMF to Fluent Bit over local TCP connection

This is the recommended option if you are using one of the AWS EMF libraries. 

See the [emf-over-tcp](emf-over-tcp) folder for example Task Definitions and Fluent Bit configuration.  

First, we create a Fluent Bit TCP input and cloudwatch output for the EMF:

```
# TCP input used for EMF payloads
[INPUT]
    Name        tcp
    Listen      0.0.0.0
    Port        25888
    Chunk_Size  32
    Buffer_Size 64
    Format      none
    Tag         emf-${HOSTNAME}
    # This tag is used by the output plugin to determine the LogStream
    # including the HOSTNAME is a way to increase the number of LogStreams.
    # The maximum throughput on a
    # single LogStream is 5 MB/s (max 1 MB at max 5 TPS).
    # In AWSVPC mode, the HOSTNAME is the ENI private IP
    # in bridge mode, the HOSTNAME is the Docker container ID

# Output for EMF over TCP -> CloudWatch
[OUTPUT]
    Name                cloudwatch
    Match               emf-*
    region              us-west-2
    log_key             log
    log_group_name      aws-emf-ecs-firelens-example-metrics
    log_stream_prefix   from-fluent-bit-
    auto_create_group   true
    log_format          json/emf
```

Fluent Bit will accept EMF over TCP port `25888`. Build this config file into a custom Fluent Bit image as shown in the base [config-file-type](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file) example. 

Now we configure the EMF libraries in your app to send EMF to TCP port `25888`/ The setup steps are similar to using the [CW Agent to collect EMF in ECS](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Embedded_Metric_Format_Generation_CloudWatch_Agent.html). 

#### Bridge Mode

If your network mode is `bridge`, you can add the following env var to your app container definition to tell the library where to emit EMF. This env var uses a Docker hostname inserted by the link to the firelens container:

```
"name": "app",
"links": [ "fluentbit" ], // the link is based on the container name
"environment": [{
                "name": "AWS_EMF_AGENT_ENDPOINT",
                "value": "tcp://fluentbit:25888"
              }]
```

And then setup a port mapping on the FireLens container:
```
"name": "fluentbit",
"portMappings": [{
                  "protocol": "tcp",
                  "containerPort": 25888
              }],
```

See the full task definition in the [emf-over-tcp](emf-over-tcp) directory.

### AWSVPC or Host Mode

If you are using host or AWSVPC network mode, the FireLens container can be reached on localhost. Note that for host mode, the host ENI is used, and so Fluent Bit must be the only process on the host using this port. 

In this case, you simply set an env var on the app container for the EMF libraries:

```
"environment": [{
                "name": "AWS_EMF_AGENT_ENDPOINT",
                "value": "tcp://127.0.0.1:25888"
              }],
```

## Tutorial 2: Send EMF to Fluent Bit via stdout

If your app's EMF library is configured to send EMF json events to stdout, then you can simply use the FireLens log driver options to send it:

```
	"logConfiguration": {
		"logDriver":"awsfirelens",
		"options": {
			"Name": "cloudwatch_logs",
			"region": "us-west-2",
			"log_key": "log",
			"log_group_name": "emf-log-group/application",
			"auto_create_group": "true",
			"log_stream_prefix": "emf-",
			"log_format": "json/emf",
			"retry_limit": "2"
		}
	},
```

However, note that if you choose this option, all events emitted to stdout will be sent with the `json/emf` header to a single log stream. 

## Other Options for Ingesting EMF

In theory, you could use almost any [Fluent Bit input plugin](https://docs.fluentbit.io/manual/pipeline/inputs) to ingest EMF. In practice, the only other commonly used option besides those mentioned above is the [tail input](https://docs.fluentbit.io/manual/pipeline/inputs/tail)- your application could write its EMF json events to a file. 