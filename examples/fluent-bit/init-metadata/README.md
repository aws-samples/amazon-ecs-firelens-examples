## FireLens Example: Using ECS Metadata Provided by Init tag

This example shows you how to use ECS Metadata in configuration. For more information on how to use the init tag feature, please see our [use case guide](https://github.com/aws/aws-for-fluent-bit/blob/mainline/use_cases/init-process-for-fluent-bit/README.md).

### Environment Variables set by init tag

Init sets the following useful env vars:

```
AWS_REGION / ECS_LAUNCH_TYPE / ECS_CLUSTER / ECS_FAMILY
ECS_TASK_ARN / ECS_TASK_ID / ECS_REVISION / ECS_TASK_DEFINITION
```

You can use these to inject metadata values in your config:

```
[FILTER]
    Name record_modifier
    Match *
    Record ecs_task_id ${ECS_TASK_ID}
```

Remember, if you use FireLens and you did not disable `enable-ecs-log-metadata`, then your logs will [already include cluster name, task ARN, and task definition](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#what-will-the-logs-collected-by-fluent-bit-look-like). 

### Templating Log Group and Stream Name using init metadata

Here's an example using logConfiguration options to create an output:

```
"logConfiguration": {
                "logDriver": "awsfirelens",
                "options": {
                    "Name": "cloudwatch_logs",
                    "region": "${AWS_REGION}",
                    "log_group_name": "${ECS_CLUSTER}/application",
                    "auto_create_group": "true",
                    "log_stream_name": "${ECS_TASK_ID}",
                    "retry_limit": "2"
                }
            },
```

Here's the same output created directly using a config file:

```
[OUTPUT]
    Name cloudwatch_logs
    region ${AWS_REGION}
    log_group_name ${ECS_CLUSTER}/application
    log_stream_name ${ECS_TASK_ID}
    auto_create_group true
    retry_limit 2
```

As you can see, the init tag makes it easy to template your configuration with ECS metadata.





