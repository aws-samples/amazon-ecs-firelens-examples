## FireLens Example: Multiple Config support - using the Fluent Bit image with init tag

This example shows you how to set multiple config files for Fluent Bit on ECS. 

This example simulates a situation: A ECS Task include two containers, the app container using httpd image, and the log router container using Fluent Bit image, which will forward the logs from the App container to CloudWatch. 

Based on above situation, we want to use multiple config feature to parse dummy data logs first, and forward the parsed dummy data logs to S3.



### Step 1: Create config files locally

**dummy-input.conf**

```
[INPUT]
    Name 					dummy
    Tag 	  			dummy.data
    Dummy   			{"data":"100 0.5 true This is Demo"}
```

**dummy-filter.conf**

 ```
[FILTER]
    Name					 parser
    Match					 dummy.*
    Key_Name 			 data
    Parser 				 dummy_test
 ```

**dummy-parser.conf**

```
[PARSER]
    Name 						dummy_test
    Format 					regex
    Regex 					^(?<INT>[^ ]+) (?<FLOAT>[^ ]+) (?<BOOL>[^ ]+) (?<STRING>.+)$
```

**dummy-s3-output.conf**

```
[OUTPUT]
    Name                   s3
    Match                  dummy.*
    bucket                 firelens-example-result
    region                 ${AWS_REGION}
    total_file_size        1M
    upload_timeout         1m
    use_put_object         On
```



### Step 2: Upload config files to S3

* create the S3 bucket `firelens-example-resources` (or any bucket name you want) to store config files
* upload above four config files to this bucket
* create the S3 bucket `firelens-example-result` (or any bucket name you want, should be same with the bucket name in `output.conf`) to receive the forwarded logs



### Step 3: Create the ECS Task

* create the ECS Task using provided `task-definition.json`, which using the Fluent Bit image with init tag

* change the `taskRoleArn` and `executionRoleArn` to your own role ARN
* change the `environment` part in the task definition FireLens configuration, copy the ARN of config files and paste it as environment variable's value. The name of environment variable requires to use the prefix `aws_fluent_bit_init_s3_`



### Step 4: Run the Task and check the result

Run the task then go to S3 to check the result in the `firelens-example-result` bucket

**original dummy log:**

 `{"data":"100 0.5 true This is Demo"}`

**dummy log forwarded to s3:** 

`{"date":"2022-08-03T03:28:36.519299Z","INT":"100","FLOAT":"0.5","BOOL":"true","STRING":"This is Demo"}`