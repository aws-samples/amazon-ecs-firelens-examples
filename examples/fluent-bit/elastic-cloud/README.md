### FireLens Example: Logging to Elastic Cloud with Fluent Bit

For all configuration parameters for Fluent Bit Elasticsearch output plugin and Elastic Cloud connection, see the [official Fluent Bit documentation](https://docs.fluentbit.io/manual/pipeline/outputs/elasticsearch).

AWS recommends that you store sensitive information, like your credentials used to connect to Elastic's Elasticsearch Service using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. Please note that the secret must be named Cloud_Auth. Value in this secret should not have any associated key and single/double quotes.

This is optional; it is also valid to simply specify the Cloud_Auth in options map:

```
"logConfiguration": {
    "logDriver": "awsfirelens",
    "options": {
        "Name": "es",
        "Port": "9243",
        "Cloud_ID": "<elastic_cloud_id>",
        "Cloud_Auth": "<elasticsearch_username>:<elasticsearch_password>",
        "Index": "elastic_firelens",
        "tls": "On",
        "tls.verify": "Off",
        "retry_limit": "2",
        "Suppress_Type_Name": "On"
    }
},
```

This is optional; you can also use regular Elasticsearch credentials to connect to Elastic Cloud in the options map:

```
"logConfiguration": {
    "logDriver": "awsfirelens",
    "options": {
        "Name": "es",
        "Port": "9243",
        "Host": "<elasticsearch_endpoint>",
        "HTTP_User": "<elasticsearch_username>",
        "HTTP_Passwd": "<elasticsearch_password>",
        "Index": "elastic_firelens",
        "tls": "On",
        "tls.verify": "Off",
        "retry_limit": "2",
        "Suppress_Type_Name": "On"
    }
},
```

This is optional; if you would like to collect logs from Fargate, just replace the first three lines of example task definition with the following:

```
    "family": "firelens-example-elastic",
    "taskRoleArn": "arn:aws:iam::XXXXXXXXXXXX:role/ecs_task_iam_role",
    "executionRoleArn": "arn:aws:iam::XXXXXXXXXXXX:role/ecs_task_execution_role",
    "cpu": "512",
    "memory": "1024",
    "requiresCompatibilities": [
        "FARGATE"
    ],
```

Note that the task and execution roles are mandatory for both Amazon ECS and AWS Fargate scenarios.
