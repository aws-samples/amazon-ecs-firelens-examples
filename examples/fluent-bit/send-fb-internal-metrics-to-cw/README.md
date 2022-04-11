# Send Fluent Bit internal Metrics to CloudWatch

This example shows you how to monitor Fluent Bit, by having it upload its own metrics to CloudWatch. 

### How it works

Fluent Bit scrapes its own [prometheus endpoint](https://docs.fluentbit.io/manual/administration/monitoring) and then parses the logs to a JSON that looks like this after FireLens adds ECS Metadata:

```
{
    "metric": "fluentbit_output_retries_total",
    "plugin": "cloudwatch_logs.0",
    "value": "0",
    "time": "1649288960583",
    "ecs_cluster": "firelens-testing",
    "ecs_task_arn": "arn:aws:ecs:ap-south-1:144718711470:task/firelens-testing/f2ad7dba413f45ddb4d92f7853b78469",
    "ecs_task_definition": "fluentbit-metrics-to-cw-firelens-example:4",
    "hostname": "ip-10-192-21-104.ap-south-1.compute.internal"
}
```

This example shows a hostname field added, which is the ENI IP of the Fargate task. The config to add hostname is optional and is commented out in the example `extra.conf` in this repo. 

In the CW Log Group for these logs, you can then create a [CloudWatch Metric Filter](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) to create metrics from the logs. Metric dimensions are customizeable and can be any value in the log JSON. 

Each Fluent Bit instance will upload these logs to a log stream named by its hostname, this should ensure unique log streams are used for each Fluent Bit instance. 

Below is a screenshot that shows how the metrics will look in CloudWatch:

![Screenshot of FB Metrics in CW](screenshot.png?raw=true "Fluent Bit Metrics in CloudWatch")

## Tutorial

### 1. Custom Fluent Bit Configuration

This example contains the following:
1. Custom Fluent Bit configuration that:
     * Enables the [Fluent Bit monitoring endpoint](https://docs.fluentbit.io/manual/administration/monitoring)
     * Uses the [exec input](https://docs.fluentbit.io/manual/pipeline/inputs/exec) to scrape that endpoint and output the results as logs
     * Filters out everything except for output, metrics, this can be customized/altered
     * Parses the data out of the logs to create a JSON event
     * Sends the JSON to CloudWatch as logs
2. A custom parser file that can parse the prometheus text into a structured JSON log. 

*FAQ: Why use the exec input to scrape the Fluent Bit prometheus metrics instead of the prometheus input or the Fluent Bit metrics input?*

This is necessary because currently, the [Fluent Bit metrics](https://docs.fluentbit.io/manual/pipeline/inputs/fluentbit-metrics) and the [Prometheus Metrics](https://docs.fluentbit.io/manual/pipeline/inputs/prometheus-scrape-metrics) inputs do not emit their data as logs, and use a separate metric pipeline that most Fluent Bit plugins do not support. However, if we use the exec input to curl the metrics, then the outputted prometheus data is turned into logs which we can parse and process easily. In the future, this experience may be improved.

2. Deploy the FireLens task

Please follow the FireLens example for [config-file-type](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file) and use the `Dockerfile` and `extra.conf` from this example. 

Then, customize the included task definition with your custom Fluent Bit image and set these environment variables to configure where the metric logs are sent:

```
{ "name": "FLUENT_BIT_METRICS_LOG_GROUP", "value": "fluent-bit-metrics-firelens-example-parsed" },
{ "name": "FLUENT_BIT_METRICS_LOG_REGION", "value": "us-west-2" }
```


### 3. Create Metric Filter

Create a [CloudWatch Metric Filter](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) on the log group to convert the JSON Logs to metrics. This can be customized as you desire. You can add additional dimensions. 

1. Filter pattern: `[metric, value, time]`

2. Metric Name: Choose a name

3. Metric Namespace

3. Metric Value: `$value`

4. Unit: Count

5. Dimensions: `fbmetric_name:$metric`

You can customize the dimensions as desired, any key in the logs can be a dimension. Here we show the Fluent Bit metric name as the only dimension.