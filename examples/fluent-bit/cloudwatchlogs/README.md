### FireLens Example: Logging to CloudWatch Logs with Fluent Bit

For documentation on Fluent Bit & CloudWatch, see: [amazon-cloudwatch-logs-for-fluent-bit](https://github.com/aws/amazon-cloudwatch-logs-for-fluent-bit)

#### CloudWatch configuration options

To minimize the possibility of log loss when using the CloudWatch plugin, consider using the recommended configuration outlined in the [CloudWatch Recommendations](https://github.com/aws/aws-for-fluent-bit/issues/340) issue.

#### What if I just want the raw log line from the container to appear in CloudWatch?

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
