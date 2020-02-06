## Amazon ECS FireLens Examples

Sample logging architectures for FireLens on Amazon ECS and AWS Fargate.

### Contributing

We want examples of as many use cases in this repository as possible! Submit a Pull Request if you would like to add something.

### Fluent Bit Examples

* [Send to CloudWatch Logs](examples/fluent-bit/cloudwatchlogs)
* [Send to Kinesis Data Firehose](examples/fluent-bit/kinesis-firehose)
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

### Fluentd Examples

* [Handling multiline logs](examples/fluentd/multiline-logs)

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
