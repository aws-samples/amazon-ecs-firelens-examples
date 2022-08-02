# FireLens Cross Account Log Push Example

This example shows you how to push logs from ECS tasks running in one account, to a CloudWatch Log Group in another account. While the example demonstrates using CloudWatch Logs, all AWS outputs in Fluent Bit support cross account. The steps in this tutorial can be followed for other AWS outputs with only slight modifications. 

This tutorial will use the following terminology:
* Task Account: The AWS Account with the running ECS task containers that produce logs. Let's assume the ID for this account is `111111111111`.
* Log Account: The AWS Account with the CloudWatch Log Group that will have the container logs. Let's assume the ID for this account is `222222222222`.


### Step 1: Configure Base Role

To push logs from the Task Account to the Log Account, Fluent Bit must use a set of base credentials which it will then use to obtain credentials for the cross account log push. 

Since the Fluent Bit FireLens sidecar runs as a container in the task, it uses the ECS Task Role for its base role/base credentials. This [role is specified in the Task Definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) with the `taskRoleArn` field. Let's assume the ARN for the ECS Task Role is `arn:aws:iam::111111111111:role/ecs-task-role`. 


The role should have the following permissions:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "*"
        }
    ]
}
```

This will allow entities that use the task role (the task containers) to call the STS AssumeRole API. This will allow Fluent Bit to assume the role in the Log Account. 

If you want to limit which roles can be assumed, that can be done with the `Resource` field in the policy above. 

### Step 2: Set up role in Log Account

Now in the other account that we will push logs to, we create a role that Fluent Bit can assume. 

The role in the Log Account will be: `arn:aws:iam::222222222222:role/cross-account-log-role`.

And it will have the following trust relationship:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "arn:aws:iam::111111111111:role/ecs-task-role"
            },
            "Action": "sts:AssumeRole",
            "Condition": {}
        }
    ]
}
```

This allows the task role in the Task Account to assume this role. Then add permissions to this role, specifically the `cloudwatch-write.json` policy in this example's directory. If you want to use a different log destination instead of CloudWatch, simply attach the permissions for your chosen destination. 

### Step 3: Configure Fluent Bit to assume the log push role

Finally, we just need to configure Fluent Bit to assume the cross account role. To do this, simply specify it with the `role_arn` parameter that all AWS outputs support:

```
role_arn     arn:aws:iam::222222222222:role/cross-account-log-role
```

The `task-definition.json` example in this directory demonstrates adding this parameter for a CloudWatch output configured via the task definition. If you have outputs in a custom configuration file, you would add the parameter like this:

```
[OUTPUT]
    Name              cloudwatch_logs
    Match *
    Log_Group_Name    cross-account-example
    Log_Stream_Prefix example-
    Auto_Create_Group On
    Role_Arn          arn:aws:iam::222222222222:role/cross-account-log-role
```

Please note that Fluent Bit will always use regional STS endpoints by default. If you have a VPC endpoint for STS that you want to use, you can specify it with the `sts_endpoint` option supported by all AWS outputs. 

