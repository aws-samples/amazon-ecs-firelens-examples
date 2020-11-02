## ECS FireLens Log Collection Tutorial

This tutorial shows multiple methods to collect logs from containers on Amazon ECS. All of the examples here use CloudWatch Logs as the log destination, however, you can change the destination to any. The goal of this tutorial is to show how to collect logs.

By default, FireLens collects stdout logs. If your logs are sent to stdout, collecting and sending them is easy, and any of the other tutorials in this repo or in the [public FireLens documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html) will suffice.

This guide is for more advanced cases where some of your logs do not go to standard out.

### Background

Make sure you are familiar with the [Fluent Bit documentation](https://docs.fluentbit.io/manual/) and with [AWS for Fluent Bit](https://github.com/aws/aws-for-fluent-bit) and its container image distribution.

### Tutorial 1: Tail log files and capture stdout

This tutorial will walk you through a scenario where there are two separate log streams emitted from a single application container:
- A Log file named `/var/log/service.log` which may or may not have its name rotated (service1.log, service2.log, etc).
- Standard out logs

The file `task-definition-tail.json` contains the full task definition.

Note one key thing: the app container is named `app`, which means that FireLens will assign the tag `app-firelens-{task ID}` to its standard out logs. Routing in Fluent Bit based on tags, and FireLens always assigns a tag of the form `{container name}-firelens-{task ID}` to standard out logs.

First, we create our Fluent Bit configuration file:

```
[INPUT]
    Name              tail
    Tag               service-log
    Path              /var/log/service*.log
    DB                /var/log/flb_service.db
    DB.locking        true
    Skip_Long_Lines   On
    Refresh_Interval  10
    Rotate_Wait       30

# Optional: append EC2 Instance Metadata to all logs
[FILTER]
    Name aws
    Match *

# Output for stdout logs
# Change the match if your container is not named 'app'
[OUTPUT]
    Name cloudwatch
    Match   app-firelens*
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/app/$(ec2_instance_id)-$(ecs_task_id)
    auto_create_group true

# Output for the file logs
[OUTPUT]
    Name cloudwatch
    Match   service-log
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/service/$(ec2_instance_id)-$(ecs_task_id)
    auto_create_group true
```

With this output, you will get all of your logs in a single log group with log stream names that tell you whether they came from stdout or the service log file. All log streams will have the EC2 instance ID and ECS Task ID. You can of course modify these patterns to suit your own needs.

Please check out the [AWS Metadata filter](https://docs.fluentbit.io/manual/pipeline/filters/aws-metadata). If you do not include it, remove `$(ec2_instance_id)` from the log stream names- that field must be injected into the logs from the filter. The other special values- `$(ecs_cluster)` and `$(ecs_task_id)` are provided by the CloudWatch Logs output when you run in ECS. For more information, [see its readme](https://github.com/aws/amazon-cloudwatch-logs-for-fluent-bit). This feature is included in AWS for Fluent Bit 2.9.0 and later.

Next, for your task to use this configuration, you should upload it to S3 and reference it in the FireLens configuration:
```
"firelensConfiguration": {
	"type": "fluentbit",
	"options": {
		"config-file-type": "s3",
		"config-file-value": "arn:aws:s3:::yourbucket/yourdirectory/tail.conf"
	}
},
```

Finally, you must mount the log file directory in to the log router container so that Fluent Bit can read it. This is show in the task definition with the volume configuration. An ephemeral volume is created which is mounted at `/var/log` in both containers.


### Tutorial 2: Same as 1 but run on Fargate

If you wish to run on Fargate, then use the Fluent Bit configuration below which removes the AWS metadata filter since EC2 metadata will not be availabe on Fargate.

```
[INPUT]
    Name              tail
    Tag               service-log
    Path              /var/log/service*.log
    DB                /var/log/flb_service.db
    DB.locking        true
    Skip_Long_Lines   On
    Refresh_Interval  10
    Rotate_Wait       30

# Output for stdout logs
# Change the match if your container is not named 'app'
[OUTPUT]
    Name cloudwatch
    Match   app-firelens*
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/app/$(ecs_task_id)
    auto_create_group true

# Output for the file logs
[OUTPUT]
    Name cloudwatch
    Match   service-log
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/service/$(ecs_task_id)
    auto_create_group true
```

Finally, with Fargate you can not store your configuration file in S3. You must bake it into a custom Fluent Bit image. Instructions can be found in the [tutorial in this repo](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file).

### Tutorial 3: Using Log4j with TCP

If your application uses log4j, you can send logs to Fluent Bit using TCP. Log4j can be set up as [explained in its documentation](https://logging.apache.org/log4j/2.x/manual/cloud.html). To capture the logs, use the Fluent Bit [TCP input](https://docs.fluentbit.io/manual/pipeline/inputs/tcp).

Note that when you use FireLens, ECS will inject the environment variable `FLUENT_HOST` at runtime. The value of this environment variable will be the IP address at which Fluent Bit is reachable within your task. You can use this env var inside your Java code to route Log4j TCP connections to Fluent Bit.

The full task definition can be found in this directory as `task-definition-tcp.json`.

Your Fluent Bit configuration file should look like the following:

```
[INPUT]
    Name        tcp
    Listen      0.0.0.0
    Port        5170
    Chunk_Size  32
    Buffer_Size 64
    Format      json
	Tag         tcp-logs

# Optional: append EC2 Instance Metadata to all logs
[FILTER]
    Name aws
    Match *

# Output for all logs
# if you app emits anything to stdout it will go here as well
# see tutorial 1 for how to send stdout to a different output
[OUTPUT]
    Name cloudwatch
    Match   *
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/$(ec2_instance_id)-$(ecs_task_id)
    auto_create_group true
```

Upload this file to S3 and reference it in your FireLens configuration:
```
"firelensConfiguration": {
	"type": "fluentbit",
	"options": {
		"config-file-type": "s3",
		"config-file-value": "arn:aws:s3:::yourbucket/yourdirectory/tcp.conf"
	}
},
```

See the notes in Tutorial 1 an 2 about EC2 metadata and how to run on Fargate.
