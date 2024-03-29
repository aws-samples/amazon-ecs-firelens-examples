### FireLens Example: Sending to multiple destinations

Sending to multiple destinations is easy with Fluent Bit. Log messages will be sent to any matching output. In this example, we have 2 firehose outputs, 1 cloudwatch_logs, and 1 datadog each of which match the tag `*`. 

For information on the tag patterns used by FireLens see here: https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#firelens-tag-and-match-pattern-and-generated-config

Since the destinations are configured in the custom config file, the app log configuration does not have any options. To use FireLens for logs, and container only needs to specify the `awsfirelens` log driver- the options are optional.
