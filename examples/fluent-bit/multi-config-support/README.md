## FireLens Example: Multiple Config support - using the Fluent Bit image with init tag

This example shows you how to set multiple config files for Fluent Bit on ECS. For more information on how to use the Multi-config support feature, please see our [use case guide](https://github.com/aws/aws-for-fluent-bit/blob/mainline/use_cases/init-process-for-fluent-bit/README.md).

This example simulates a situation: A ECS Task include two containers, the app container using your own image, and the log router container using Fluent Bit image, which will forward the logs from the App container to CloudWatch. 

**Based on above situation, suppose your app container also generate a log file which path is `/logs/app.log`. we want to use multiple config feature to forward the content of this log file, parse it first, and forward the parsed logs to S3.**



### Step 1: Create config files locally

**tail-input.conf**

```
[INPUT]
    Name            tail
    Tag             app
    Path            /logs/app.log
    Read_from_Head  True
```

**your-filter.conf**

 ```
[FILTER]
    Name        parser
    Match       app
    Key_Name    data
    Parser      app_test
 ```

**your-parser.conf**

```
[PARSER]
    Name        app_test
    Format      regex
    Regex       ^(?<INT>[^ ]+) (?<FLOAT>[^ ]+) (?<BOOL>[^ ]+) (?<STRING>.+)$
```

**dummy-s3-output.conf**

```
[OUTPUT]
    Name                         s3
    Match                        app
    bucket                       your-bucket
    region                       ${AWS_REGION}
    total_file_size              1M
    upload_timeout               1m
    use_put_object               On
```



**Note:** you can find these config files in the `config-files` directory of this example, please modify them according to the actual situation to match your needs.

### Step 2: Upload config files to S3

* create the S3 bucket `your-bucket` to store config files
* upload above config files to this bucket
* create the S3 bucket `your-result` to receive the forwarded logs



### Step 3: Create the ECS Task

* create the ECS Task using provided `task-definition.json`, which using the Fluent Bit image with init tag
* change the `taskRoleArn` and `executionRoleArn` to your own role ARN
* change the `environment` part in the task definition FireLens configuration, copy the ARN of config files and paste it as environment variable's value. The name of environment variable requires to use the prefix `aws_fluent_bit_init_s3_`

**Note:** you need mount the log file directory in to the log router container so that Fluent Bit can read it. This is shown in the task definition with the volume configuration. An ephemeral volume is created which is mounted at `/logs` in both containers.
For ECS FireLens Log Collection Tutorial, [please see here.](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/ecs-log-collection)



### Step 4: Run the Task and check the result

Run the task then go to S3 to check the result in `your-result` bucket. Content of the log file of your app container should have been forwarded to S3.
