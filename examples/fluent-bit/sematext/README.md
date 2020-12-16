### FireLens Example: Logging to Sematext with Fluent Bit

For documentation on sending your logs from AWS ECS running on either AWS Fargate or AWS EC2 to Sematext Logs, see: [Elastic Container Service (ECS) Logs Integration](https://sematext.com/docs/integration/ecs-logs/).

AWS recommends that you store sensitive information, like your Sematext `LOGS_TOKEN` using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. This is optional; it is also valid to simply specify the `LOGS_TOKEN` in the `logConfiguration`. Note, the `URI` field needs to have a `/` in front of the token:

```json
"logConfiguration": {
    "logDriver":"awsfirelens",
    "options": {
        "Name": "http",
        "Match": "*",
        "Header": "sourceName nginx",
        "Host": "logs-ecs-receiver.sematext.com",
        "URI": "/<LOGS_TOKEN>",
        "Port": "443",
        "TLS": "on",
        "Format": "json",
        "compress": "gzip"
    }
},
```

**Note: If you are using the EU region of Sematext, use this Host value: `"Host": "logs-ecs-receiver.eu.sematext.com"`**
