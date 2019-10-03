### FireLens Example: Adding Keys to the Log Events

With the custom configuration file in this example, you can add a key to each log message. This is similar to how FireLens adds ECS Metadata to your logs.

In this example, we add a field called `app-version`- this will allow us to correlate log messages with the version of our app that generated them. The App Version is set via an environment variable, which can be referenced in the Fluent Bit configuration file.

Assuming ECS Log Metadata is enabled, the final log events in Firehose will look something like the following:
```
{
    "source": "stdout",
    "app-version": "v1.1.14",
    "log": "172.17.0.1 - - [03/Oct/2019:00:06:20 +0000] \"GET / HTTP/1.1\" 200 612 \"-\" \"curl/7.54.0\" \"-\"",
    "container_id": "e54cccfac2b87417f71877907f67879068420042828067ae0867e60a63529d35",
    "container_name": "/ecs-demo-6-container2-a4eafbb3d4c7f1e16e00"
    "ecs_cluster": "mycluster",
    "ecs_task_arn": "arn:aws:ecs:us-east-2:01234567891011:task/mycluster/3de392df-6bfa-470b-97ed-aa6f482cd7a6",
    "ecs_task_definition": "demo:7",
    "ec2_instance_id": "i-06bc83dbc2ac2fdf8"
}
```

Keys can be added and removed via the record_modifier filter- for more information see the [Fluent Bit documentation](https://fluentbit.io/documentation/0.12/filter/record_modifier.html).
