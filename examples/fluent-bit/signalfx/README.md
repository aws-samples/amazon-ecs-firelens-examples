### FireLens Example: Logging to SignalFx with Fluent Bit

The SignalFx output plugin for Fluent Bit sends log-based metrics to SignalFx. Refer to the [SignalFx Fluent Bit integration](https://docs.signalfx.com/en/latest/integrations/integrations-reference/integrations.fluent.bit.html) documentation for details.

AWS recommends that you store sensitive information (like your SignalFx Token) using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example [task definition](task-definition.json). This is optional - it is also valid to simply specify the Token in the options map:

```
"logConfiguration": {
    "logDriver": "awsfirelens",
    "options": {
        "Name": "SignalFx",
        "MetricName": "com.firelens.example",
        "MetricType": "counter",
        "Dimensions": "ecs_cluster, container_name",
        "Token": "XXXXXXXXXXXX"
    }
}
```
