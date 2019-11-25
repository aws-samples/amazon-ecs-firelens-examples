### FireLens Example: Logging to CloudWatch Logs with Fluent Bit

For documentation on Fluent Bit & CloudWatch, see: [amazon-cloudwatch-logs-for-fluent-bit](https://github.com/aws/amazon-cloudwatch-logs-for-fluent-bit)

#### What if I just want the raw log line from the container to appear in CloudWatch?

By default, FireLens will send a JSON event with the raw log line encapsulated in a `log` field. ECS Metadata will also be added. If you just want the raw log line, add the `log_key` option to your log configuration:

```
	"logConfiguration": {
		"logDriver":"awsfirelens",
		"options": {
			"Name": "cloudwatch",
			"region": "us-west-2",
			"log_key": "log",
			"log_group_name": "firelens-fluent-bit",
			"auto_create_group": "true",
			"log_stream_prefix": "from-fluent-bit"
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
