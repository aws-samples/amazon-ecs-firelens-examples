## FireLens Example: Logging to CloudWatch Logs with Fluent Bit

### CloudWatch Golang Plugin vs CloudWatch C Plugin

There are two Fluent Bit output plugins for sending to Amazon CloudWatch Logs:
* Plugin name `cloudwatch`: the original golang plugin with extensive templating support, including injecting ECS metadata into log stream and group name templates. See its [documentation](https://github.com/aws/amazon-cloudwatch-logs-for-fluent-bit). This directory contains an example task definition thats demonstrates this plugin's ability to inject ECS task metadata into log group and stream names.
* Plugin name `cloudwatch_logs`: the newer and higher performance cloudwatch plugin built in C in the Fluent Bit upstream code base. It has more limited log group and stream name templating support. See its [documentation](https://docs.fluentbit.io/manual/pipeline/outputs/cloudwatch). This directory contains an example task definition for the high performance plugin without templating. The log stream name will be set to be `{log_stream_prefix}{log tag}` and [FireLens sets the log tag](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#firelens-tag-and-match-pattern-and-generated-config) to be `{container name in task definition}-firelens-{task ID}`. So the log stream name for this example will be `stdout-stderr-app-firelens-{task ID}`. 


For more on the AWS Go outputs vs AWS C outputs, check out the [FAQ entry in our debugging guide](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#aws-go-plugins-vs-aws-core-c-plugins). 

#### Recommended cloudwatch_logs configuration options

To minimize the possibility of log loss when sending to CloudWatch, consider using the recommended configuration outlined in the [CloudWatch Recommendations](https://github.com/aws/aws-for-fluent-bit/issues/340) issue.

### What if I just want the raw log line from the container to appear in CloudWatch?

The example shown here is for the `cloudwatch` plugin, however, the `cloudwatch_logs` plugin has the same `log_key` option.

By default, FireLens will send a JSON event with the raw log line encapsulated in a `log` field. ECS Metadata will also be added. If you just want the raw log line, add the `log_key` option to your log configuration:

```
	"logConfiguration": {
		"logDriver":"awsfirelens",
		"options": {
			"Name": "cloudwatch",
			"region": "us-west-2",
			"log_key": "log",
			"log_group_name": "/aws/ecs/containerinsights/$(ecs_cluster)/application",
			"auto_create_group": "true",
			"log_stream_name": "$(ecs_task_id)",
			"retry_limit": "2"
		}
	},
```

This field instructs the CloudWatch plugin to only send the value of the log key. You should additionally disable ECS Metadata to prevent Fluent Bit from performing unnecessary data processing:

```
	"firelensConfiguration": {
		"type": "fluentbit",
		"options": {
			"enable-ecs-log-metadata": "false",
		}
	},
```

			"log_key": "log",
The following table shows how your logs will appear in CloudWatch with and without `log_key` and `enable-ecs-log-metadata` if "`my_sample log`" is printed to stdout:

|logConfiguration.options contains|firelensConfiguration.options contains|received by CloudWatch|description|comments|
|-|-|-|-|-|
|||{<br>    "container_id": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXXXX",<br>    "container_name": "my_container_name",<br>    "ecs_cluster": "my_cluster_name",<br>    "ecs_task_arn": "arn:aws:ecs:region:9876543210:task/my_task_arn",<br>    "ecs_task_definition": "my_task_definition:revision_number",<br>    "log": "my_sample log",<br>    "source": "stdout"<br>}|no log key is set and ecs-log-metadata is enabled by default||
||"enable-ecs-log-metadata":"false"|{<br>    "container_id": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXXXX",<br>    "container_name": "my_container_name",<br>    "log": "my_sample log",<br>    "source": "stdout"<br>}|no log key is set and ecs-log-metadata is disabled||
|"log_key":"log"||"my_sample log"|log key set to "log" and ecs-log-metadata is enabled by default|less efficient|
|"log_key":"log"|"enable-ecs-log-metadata":"false"|"my_sample log"|log key set to "log" and ecs-log-metadata is disabled|<b>more efficient</b>|
