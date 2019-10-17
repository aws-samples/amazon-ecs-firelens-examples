### FireLens Example: Parsing container stdout logs that are serialized JSON

As of [AWS for Fluent Bit](https://hub.docker.com/r/amazon/aws-for-fluent-bit) version 1.3, an external configuration file is not needed to parse JSON. The config shown in [extra.conf](extra.conf) is included in the image. Simply reference it in your FireLens configuration:

```
"firelensConfiguration": {
    "type": "fluentbit",
    "options": {
        "config-file-type": "file",
        "config-file-value": "/fluent-bit/configs/parse-json.conf"
    }
},
```

This Fluent Bit config file will parse any logs that are JSON.
For example, if the logs at your destination looked like this without JSON parsing:

```
{
    "source": "stdout",
    "log": "{\"requestID\": \"b5d716fca19a4252ad90e7b8ec7cc8d2\", \"requestInfo\": {\"ipAddress\": \"204.16.5.19\", \"path\": \"/activate\", \"user\": \"TheDoctor\"}}",
    "container_id": "e54cccfac2b87417f71877907f67879068420042828067ae0867e60a63529d35",
    "container_name": "/ecs-demo-6-container2-a4eafbb3d4c7f1e16e00"
    "ecs_cluster": "mycluster",
    "ecs_task_arn": "arn:aws:ecs:us-east-2:01234567891011:task/mycluster/3de392df-6bfa-470b-97ed-aa6f482cd7a6",
    "ecs_task_definition": "demo:7"
    "ec2_instance_id": "i-06bc83dbc2ac2fdf8"
}
```

Then with JSON parsing they'll look like this:

```
{
    "source": "stdout",
    "container_id": "e54cccfac2b87417f71877907f67879068420042828067ae0867e60a63529d35",
    "container_name": "/ecs-demo-6-container2-a4eafbb3d4c7f1e16e00"
    "ecs_cluster": "mycluster",
    "ecs_task_arn": "arn:aws:ecs:us-east-2:01234567891011:task/mycluster/3de392df-6bfa-470b-97ed-aa6f482cd7a6",
    "ecs_task_definition": "demo:7"
    "ec2_instance_id": "i-06bc83dbc2ac2fdf8"
    "requestID": "b5d716fca19a4252ad90e7b8ec7cc8d2",
    "requestInfo": {
        "ipAddress": "204.16.5.19",
        "path": "/activate",
        "user": "TheDoctor"
    }
}
```

As you can see, the serialized JSON is expanded into top level fields in the final JSON. For more information on JSON parsing see, the [Fluent Bit documentation](https://docs.fluentbit.io/manual/filter/parser).
