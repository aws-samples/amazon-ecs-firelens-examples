# Fluent Bit Throughput Measurement

This example demonstrates how to measure log throughput in Fluent Bit using a Lua filter that outputs metrics in AWS CloudWatch EMF (Embedded Metric Format).

## Overview

The Lua script calculates throughput every 10 seconds by tracking the total bytes processed and emits CloudWatch metrics with the following dimensions:
- `LogGroup`: The CloudWatch log group name
- `TaskArn`: The ECS task ARN (if available)

## Configuration

### Environment Variables

- `LOG_GROUP` - CloudWatch log group name (default: `/aws/fluent-bit/analytics`)
- `METRIC_NAMESPACE` - CloudWatch metric namespace (default: `aws-for-fluent-bit/LogThroughput`)

### Usage Options

**Option 1: External Lua file**
```ini
[FILTER]
    Name lua
    Match *
    script measure-throughput.lua
    call measure_throughput
```

**Option 2: Inline Lua code**
See `fluent-bit.conf` for the inline version.

## Metrics

The script emits the following CloudWatch metric:
- **Name**: `ThroughputMbps`
- **Unit**: Megabytes/Second
- **Dimensions**: LogGroup, TaskArn

### Samples

```
{"ThroughputMbps":2.720108,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1}]}],"Timestamp":1762561423000},"LogGroup":"/aws/fluent-bit/analytics"}
2025-M-DT00:23:53.009000+00:00 fluent-bit/log-router/TASK {"ThroughputMbps":2.683549,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1}]}],"Timestamp":1762561433000},"LogGroup":"/aws/fluent-bit/analytics"}
2025-M-DT00:24:03.009000+00:00 fluent-bit/log-router/TASK {"ThroughputMbps":2.692664,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1}]}],"Timestamp":1762561443000},"LogGroup":"/aws/fluent-bit/analytics"}
```

## How It Works

1. Tracks bytes processed from each log record
2. Every 10 seconds, calculates throughput in MB/s
3. Outputs EMF-formatted JSON to stdout
4. Passes through original log records unchanged

## Notes

- The TaskArn is cached from the first record containing `ecs_task_arn` field
- Throughput calculation uses an approximate byte size based on record serialization
- EMF metrics are printed to stdout and should be sent to CloudWatch Logs for automatic metric extraction
