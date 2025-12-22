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

The script emits the following CloudWatch metrics:
- **Name**: `ThroughputMbps`
  - **Unit**: Megabytes/Second
  - **Dimensions**: LogGroup, TaskArn
- **Name**: `ThroughputKbps`
  - **Unit**: Kilobytes/Second
  - **Dimensions**: LogGroup, TaskArn

### Samples

```
{"ThroughputMbps":2.720108,"ThroughputKbps":2720.108,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1},{"Name":"ThroughputKbps","Unit":"Kilobytes/Second","StorageResolution":1}]}],"Timestamp":1762561423000},"LogGroup":"/aws/fluent-bit/analytics"}
2025-M-DT00:23:53.009000+00:00 fluent-bit/log-router/TASK {"ThroughputMbps":2.683549,"ThroughputKbps":2683.549,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1},{"Name":"ThroughputKbps","Unit":"Kilobytes/Second","StorageResolution":1}]}],"Timestamp":1762561433000},"LogGroup":"/aws/fluent-bit/analytics"}
2025-M-DT00:24:03.009000+00:00 fluent-bit/log-router/TASK {"ThroughputMbps":2.692664,"ThroughputKbps":2692.664,"TaskArn":"arn:aws:ecs:us-west-2:444455556666:task/fluentbit-workshop/TASK","_aws":{"CloudWatchMetrics":[{"Namespace":"aws-for-fluent-bit/LogThroughput","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second","StorageResolution":1},{"Name":"ThroughputKbps","Unit":"Kilobytes/Second","StorageResolution":1}]}],"Timestamp":1762561443000},"LogGroup":"/aws/fluent-bit/analytics"}
```

## How It Works

1. Tracks bytes processed from each log record
2. Every 10 seconds, calculates throughput in both MB/s and KB/s
3. Outputs EMF-formatted JSON to stdout via fluent-bit
4. Passes through original log records unchanged

## Observing Metrics in CloudWatch

Once the EMF metrics are sent to CloudWatch Logs, they are automatically extracted as CloudWatch metrics. You can view and query these metrics in several ways:

### CloudWatch Metrics Console
Navigate to CloudWatch > Metrics > Custom Namespaces > `aws-for-fluent-bit/LogThroughput` to view the `ThroughputMbps` and `ThroughputKbps` metrics.

### CloudWatch Insights Query
Use the following CloudWatch Logs Insights queries to analyze throughput data:

```sql
SELECT AVG(ThroughputMbps) FROM SCHEMA("aws-for-fluent-bit/LogThroughput", LogGroup,TaskArn) GROUP BY TaskArn
```

```sql
SELECT AVG(ThroughputKbps) FROM SCHEMA("aws-for-fluent-bit/LogThroughput", LogGroup,TaskArn) GROUP BY TaskArn
```

### CloudWatch Dashboard Widget
You can create a dashboard widget using the metrics expressions above:

```json
{
  "type": "metric",
  "properties": {
    "metrics": [
      [{ 
        "expression": "SELECT AVG(ThroughputMbps) FROM SCHEMA(\"aws-for-fluent-bit/LogThroughput\", LogGroup,TaskArn) GROUP BY TaskArn", 
        "label": "ThroughputMbps", 
        "id": "q1", 
        "period": 20 
      }],
      [{ 
        "expression": "SELECT AVG(ThroughputKbps) FROM SCHEMA(\"aws-for-fluent-bit/LogThroughput\", LogGroup,TaskArn) GROUP BY TaskArn", 
        "label": "ThroughputKbps", 
        "id": "q2", 
        "period": 20 
      }]
    ],
    "view": "timeSeries",
    "region": "${AWS::Region}",
    "title": "Log Throughput by Task"
  }
}
```

## Notes

- The TaskArn is cached from the first record containing `ecs_task_arn` field
- Throughput calculation uses an approximate byte size based on record serialization
- EMF metrics are printed to stdout and should be sent to CloudWatch Logs for automatic metric extraction
