## Amazon ECS FireLens Examples

Sample logging architectures for FireLens on Amazon ECS and AWS Fargate.

### Contributing

We want examples of as many use cases in this repository as possible! Submit a Pull Request if you would like to add something.

### ECS Log Collection
* [Collecting Log Files, stdout and log4j logs in ECS](examples/fluent-bit/ecs-log-collection)

### Basic FireLens examples
* [Using the 'file' 'config-file-type'](examples/fluent-bit/config-file-type-file)
* [Cross Account Log Forwarding](examples/fluent-bit/cross-account)
* [Source multiple configs from S3 or files](examples/fluent-bit/multi-config-support)
* [Using EFS to store configuration files](examples/fluent-bit/efs)
* [Specifying buffer limit size with 'Fluentd' log driver](examples/fluent-bit/log-driver-buffer-limit)
* [How to set Fluentd and Fluent Bit input parameters (including Mem_Buf_Limit) in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/)
* [How to prevent OOMKills (Out of Memory) in FireLens](examples/fluent-bit/oomkill-prevention)

### Multiline Examples
* [Concat multiline logs using regex parsers](examples/fluent-bit/filter-multiline)
* [Concat partial/split container logs](examples/fluent-bit/filter-multiline-partial-message-mode)

### Monitoring Fluent Bit
* [Send Fluent Bit internal metrics to CloudWatch](examples/fluent-bit/send-fb-internal-metrics-to-cw)
* [Fluent Bit Container Health Check Options](examples/fluent-bit/health-check)
* [CPU, Disk, and Memory Usage Monitoring with ADOT](examples/fluent-bit/adot-resource-monitoring/)

### Fluent Bit Examples

* [Send Logs to CloudWatch Logs](examples/fluent-bit/cloudwatchlogs)
* [Send EMF Metrics to CloudWatch Logs](examples/fluent-bit/cloudwatchlogs-emf)
* [Send to Kinesis Data Firehose](examples/fluent-bit/kinesis-firehose)
* [Send to Kinesis Data Stream](examples/fluent-bit/kinesis-stream)
* [Send to S3](examples/fluent-bit/s3)
* [Send to Amazon OpenSearch Service](examples/fluent-bit/amazon-opensearch)
* [Enable Debug Logging](examples/fluent-bit/enable-debug-logging)
* [Forward to a Fluentd or Fluent Bit Log Aggregator](examples/fluent-bit/forward-to-aggregator)
* [Parse Serialized JSON](examples/fluent-bit/parse-json)
* [Parse common log formats](examples/fluent-bit/parse-common-log-formats)
* [Parse Envoy Access Logs from AWS App Mesh](examples/fluent-bit/parse-envoy-app-mesh)
* [Send to multiple destinations](examples/fluent-bit/send-to-multiple-destinations)
* [Add custom metadata to logs](examples/fluent-bit/add-keys)
* [Datadog monitoring](examples/fluent-bit/datadog)
* [SignalFx monitoring](examples/fluent-bit/signalfx)
* [New Relic Logs](examples/fluent-bit/newrelic)
* [Sumo Logic](examples/fluent-bit/sumologic)
* [SolarWinds Loggly](examples/fluent-bit/solarwinds-loggly)
* [Sematext Logs](examples/fluent-bit/sematext)
* [Logstash](examples/fluent-bit/logstash)
* [Elastic Cloud](examples/fluent-bit/elastic-cloud)

### Fluentd Examples

* [Handling multiline logs](examples/fluentd/multiline-logs)

### Splitting an applications logs into multiple strings

Artifacts for the blog [Splitting an application’s logs into multiple streams: a Fluent tutorial](https://aws.amazon.com/blogs/opensource/splitting-application-logs-multiple-streams-fluent/)

* [Splitting logs artifacts/examples](examples/splitting-log-streams)

### Setup for the examples

Before you use FireLens, familiarize yourself with [Amazon ECS](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_GetStarted_EC2.html) and with the [FireLens documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html).

In order to use these examples, you will need the following IAM resources:
* A [Task IAM Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) with permissions to send logs to your log destination. Each of the examples in this repository that needs additional permissions has a sample policy.
* A [Task Execution Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html). This role is used by the ECS Agent to make calls on your behalf. If you enable logging for your FireLens container with the [`awslogs` Docker Driver](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html), you will need permissions for CloudWatch. You also need to give it S3 permissions if you are pulling an external Fluent Bit or Fluentd configuration file from S3. See the the [FireLens documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_firelens.html) for more.

Here is an example inline policy with S3 access for FireLens:

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::examplebucket/folder_name/config_file_name"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetBucketLocation"
      ],
      "Resource": [
        "arn:aws:s3:::examplebucket"
      ]
    }
  ]
}
```

### Using the Examples

You must update each Task Definition to reflect your own needs. Replace the IAM roles with your own roles. Update the log configuration with the values that you desire. And replace the app image with your own application image.

Additionally, several of these examples use a custom Fluent Bit/Fluentd configuration file in S3. You must upload it to your own bucket, and change the S3 ARN in the example Task Definition.

If you are using ECS on Fargate, then pulling a config file from S3 is not currently supported. Instead, you must create a custom Docker image with the config file.

Dockerfile to add a custom configs:
```
FROM amazon/aws-for-fluent-bit:latest
ADD extra.conf /extra.conf
```

Then update the `firelensConfiguration` `options` in the Task Definition to the following:
```
"options": {
    "config-file-type": "file",
    "config-file-value": "/extra.conf"
}
```

## License Summary

This sample code is made available under the MIT-0 license. See the LICENSE file.
