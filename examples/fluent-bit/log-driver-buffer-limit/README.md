### FireLens Example: Specifying buffer limit size with 'Fluentd' log driver

`log-driver-buffer-limit` option is now supported on ECS EC2 and ECS Fargate with PV1.4. This setting helps FireLens to configure the Fluentd Docker Log Driver field [fluentd-buffer-limit](https://docs.docker.com/config/containers/logging/fluentd/#fluentd-buffer-limit). 

As we know, FireLens is a container log router for Amazon ECS and AWS Fargate that gives customers extensibility to use the breadth of services at AWS or partner solutions for log analytics and storage. FireLens works with Fluentd and Fluent Bit and makes it easy to use these two popular open source logging projects.

![FireLens](https://d2908q01vomqb2.cloudfront.net/fe2ef495a1152561572949784c16bf23abb28057/2019/11/16/Screen-Shot-2019-09-26-at-5.21.35-PM-1024x572.png)

The diagram above shows how FireLens works. Container standard out logs are sent to the FireLens container over a Unix socket via the [Fluentd Docker Log Driver](https://docs.docker.com/config/containers/logging/fluentd/). 

To set Fluentd or Fluent Bit output plugin configurations in FireLens, we can configure them in the log configuration section of the Task Definition. We specify `logDriver` as `awsfirelens` to use FireLens and provide config details for the Fluentd or Fluent Bit output plugin in `options` field. The option `log-driver-buffer-limit` will specify the limit size for number of events buffered on the memory. It can help to resolve potential log loss issue because high throughput could result in running out of memory for buffer inside of Docker and the log driver must discard messages from the buffer to add new ones if the buffer is full. This will produce tons of error messages in Docker and the lost logs will make it impossible to troubleshoot problems in the application. Increasing the buffer limit size by configuring the option `log-driver-buffer-limit` could be a good approach to avoid this problem.

By default, the buffer limit will be set to `1MiB`. In order to increase or decrease the size, you can customize the `log-driver-buffer-limit` option in your log configuration:

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
			"log-driver-buffer-limit": "2097152",
			"retry_limit": "2"
		}
	},
```

The unit for this field is `byte` so above task definition will set the buffer limit size to `2MiB`.

*Note*: 
- Fargate PV1.3 is on deprecation. The feature is only supported after Fargate PV1.4.
- The value for `log-driver-buffer-limit` should be an integer between 0 and 536870912 (`512MiB`).
- The total amount of memory allocated at the task level must be greater than the amount of memory allocated for all containers in addition to the memory buffer limit for the FireLens log driver.
- The total amount of memory buffer specified must be lesser than 536870912 (`512MiB`) when the container `memory` and `memoryReservertion` values aren't specified. More specifically, you can have an app container with `awsfirelens` log driver and option `log-driver-buffer-limit` set to `300MiB`. However, you won't be allowed to run tasks if you have more than two containers with `log-driver-buffer-limit` as `300MiB`(`300MiB` * 2 > `512MiB`).