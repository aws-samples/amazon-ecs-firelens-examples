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

**Please note**: In order to have ECS metadata added to your logs, `enable-ecs-log-metadata` must be enabled (set to `true`) in the `firelensConfiguration` of your Task Definition. This is the default value, so as long as you have not explicitly set it to `false` it will be enabled. 

This example shows a hostname field added, which is the ENI IP of the Fargate task. The config to add hostname is optional and is commented out in the example `extra.conf` in this repo. 

In the CW Log Group for these logs, you can then create a [CloudWatch Metric Filter](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) to create metrics from the logs. Metric dimensions are customizeable and can be any value in the log JSON. 

Each Fluent Bit instance will upload these logs to a log stream named by its hostname, this should ensure unique log streams are used for each Fluent Bit instance. 

Below is a screenshot that shows how the metrics will look in CloudWatch:

![Screenshot of FB Metrics in CW](screenshot.png?raw=true "Fluent Bit Metrics in CloudWatch")

*There are two tutorials in this guide*- the first is for Fluent Bit's internal plugin metrics- error and retries and data processed by each plugin. The second shows you how to also enable Fluent Bit's storage metrics. 

## Tutorial 1: Fluent Bit Plugin Metrics

### 1. Custom Fluent Bit Configuration

This example contains the following:
1. Custom Fluent Bit configuration that:
     * Enables the [Fluent Bit monitoring endpoint](https://docs.fluentbit.io/manual/administration/monitoring)
     * Uses the [exec input](https://docs.fluentbit.io/manual/pipeline/inputs/exec) to scrape that endpoint and output the results as logs
     * *Filters out everything except for output metrics*, this can be customized/altered. Un-comment the lines in the provided configuration file and remove the filter that exludes metrics not for outputs. 
     * Parses the data out of the logs to create a JSON event
     * Sends the JSON to CloudWatch as logs
2. A custom parser file that can parse the prometheus text into a structured JSON log. 

*FAQ: Why use the exec input to scrape the Fluent Bit prometheus metrics instead of the prometheus input or the Fluent Bit metrics input?*

This is necessary because currently, the [Fluent Bit metrics](https://docs.fluentbit.io/manual/pipeline/inputs/fluentbit-metrics) and the [Prometheus Metrics](https://docs.fluentbit.io/manual/pipeline/inputs/prometheus-scrape-metrics) inputs do not emit their data as logs, and use a separate metric pipeline that most Fluent Bit plugins do not support. However, if we use the exec input to curl the metrics, then the outputted prometheus data is turned into logs which we can parse and process easily. In the future, this experience may be improved.

For a quick setup, use the [built-in plugin metrics configuration file](https://github.com/aws/aws-for-fluent-bit/blob/mainline/configs/plugin-metrics-to-cloudwatch.conf) available in AWS for Fluent Bit 2.29.1+. This built-in configuration means that you do not have to build a custom Fluent Bit image. However, please note that this built-in configuration enables all plugin metrics (not just output plugin metrics).

```
			"firelensConfiguration": {
				"type": "fluentbit",
				"options": {
					"config-file-type": "file",
					"config-file-value": "/fluent-bit/configs/plugin-metrics-to-cloudwatch.conf"
				}
			},
```

#### DIY Setup

If you have your own custom configuration already, you can import/include the built-in config:

```
@INCLUDE /fluent-bit/configs/plugin-metrics-to-cloudwatch.conf
```

Please note that the [built-in config](https://github.com/aws/aws-for-fluent-bit/blob/mainline/configs/plugin-metrics-to-cloudwatch.conf) includes a `[SERVICE]` section. This section can only be set once, so importing the built-in means that you can not have your own custom `[SERVICE]` section..

Alternatively, the [`extra.conf`](extra.conf) and [`fb_metrics_parser.conf`](fb_metrics_parser.conf) in this directory show the necessary config content. 

Please note, you *must set the environment variables described in step 2 and the metric filter in step 3*.

### 2. Deploy the FireLens task

Please follow the FireLens example for [config-file-type](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file) and use the [`Dockerfile`](Dockerfile) and [`extra.conf`](extra.conf) from this example. 

Then, customize the included task definition with your custom Fluent Bit image and set these environment variables to configure where the metric logs are sent:

```
{ "name": "FLUENT_BIT_METRICS_LOG_GROUP", "value": "fluent-bit-metrics-firelens-example-parsed" },
{ "name": "FLUENT_BIT_METRICS_LOG_REGION", "value": "us-west-2" }
```


### 3. Create Metric Filter

Create a [CloudWatch Metric Filter](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) on the log group to convert the JSON Logs to metrics. This can be customized as you desire. You can add additional dimensions. 

1. Filter pattern: `{ $.value = * }`

2. Metric Name: Choose a name

3. Metric Namespace: Choose a namespace

3. Metric Value: `$.value`

4. Unit: Count

5. Dimensions: `fbmetric_name:$.metric`

You can customize the dimensions as desired, any key in the logs can be a dimension. Here we show the Fluent Bit metric name as the only dimension.



## Tutorial 2: Fluent Bit Storage Metrics

To also send storage metrics, the technique the same. The difference is that Fluent Bit publishes storage metrics to a different HTTP path that vends the metrics in JSON instead of prometheus format. Therefore, a different input and filter is needed. 

### Quick setup: use built-in configuration

For a quick setup, use the [built-in plugin + storage metrics configuration file](https://github.com/aws/aws-for-fluent-bit/blob/mainline/configs/plugin-and-storage-metrics-to-cloudwatch.conf) available in AWS for Fluent Bit 2.29.1+. This built-in configuration means that you do not have to build a custom Fluent Bit image. However, please note that this built-in configuration enables all plugin metrics (not just output plugin metrics).

```
			"firelensConfiguration": {
				"type": "fluentbit",
				"options": {
					"config-file-type": "file",
					"config-file-value": "/fluent-bit/configs/plugin-and-storage-metrics-to-cloudwatch.conf"
				}
			},
```

Then set these environment variables on your FireLens container to configure where the metric logs are sent:

```
{ "name": "FLUENT_BIT_METRICS_LOG_GROUP", "value": "fluent-bit-metrics-firelens-example-parsed" },
{ "name": "FLUENT_BIT_METRICS_LOG_REGION", "value": "us-west-2" }
```

Please note, *you are not done yet, you must still create the metric filters on your log group described in this guide*.

### Full DIY Tutorial: Modify an existing custom config to add storage metrics

#### 1. Add an input to scrape the storage endpoint

```
[INPUT]
    Name exec
    Command curl -s http://127.0.0.1:2020/api/v1/storage && echo ""
    Interval_Sec 5
    Tag fb_metrics-storage
```

#### 2. Add a filter to parse the JSON formatted storage metrics

```
[FILTER]
    Name parser
    Match fb_metrics-storage
    Key_Name exec
    Parser json
```

The [`extra.conf`](extra.conf) file in this directory has the configuration to send storage metrics commented out on lines 42 to 63. Un-commment these lines (remove the '#' signs) to enable sending this data to CW. The data can be sent by the same output as the prometheus metrics. 

Alternatively you can import/include the built-in config:

```
@INCLUDE /fluent-bit/configs/plugin-and-storage-metrics-to-cloudwatch.conf
```

Please note that the [built-in config](https://github.com/aws/aws-for-fluent-bit/blob/mainline/configs/plugin-and-storage-metrics-to-cloudwatch.conf) includes a `[SERVICE]` section. This section can only be set once, so importing the built-in means that you can not have your own custom `[SERVICE]` section..

#### 3. Create a metric filter for the storage metrics

The storage metric JSON data in your CW log stream will look like this:

```
{
	"date": 1656218290.714451,
	"storage_layer": {
		"chunks": {
			"total_chunks": 2,
			"mem_chunks": 2,
			"fs_chunks": 0,
			"fs_chunks_up": 0,
			"fs_chunks_down": 0
		}
	},
	"input_chunks": {
		"my_input_alias": {
			"status": {
				"overlimit": false,
				"mem_size": "104b",
				"mem_limit": "0b"
			},
			"chunks": {
				"total": 1,
				"up": 1,
				"down": 0,
				"busy": 0,
				"busy_size": "0b"
			}
		}
	}
}
```

You must determine which data inside this structure is important to your use case. The `input_chunks` structure allows you to track the storage used by specific input definitions- notice that the name of the key is the name of the input definition. This can be customized with an alias, which you can see how to configure here: [Fluent Bit monitoring: configuring aliases](https://docs.fluentbit.io/manual/administration/monitoring#configuring-aliases).

In this example, we show how to expose the `mem_chunks` metric in CloudWatch. This tracks the count of chunks of data stored in memory by Fluent Bit. 

A. Filter pattern: `{ $.storage_layer.chunks.mem_chunks = * }`

B. Metric Name: Choose a name

C. Metric Namespace: Choose a namespace

D. Metric Value: `$.storage_layer.chunks.mem_chunks`

E. Unit: Count

F. Dimensions: Same as the first tutorial, you must choose this yourself, we recommend choosing uncommenting the lines in [`extra.conf`](extra.conf) to add HOSTNAME to the metric data and choosing this as a dimension. This can be accomplished by setting the dimension to `hostname:$hostname`. 
