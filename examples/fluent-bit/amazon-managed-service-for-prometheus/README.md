## FireLens Example: Prometheus forward metrics - using the Fluent Bit image with init tag

This example shows how to forward metrics with Fluent Bit's prometheus remote write output plugin to an Amazon Managed Service for Prometheus workspace.

### Step 1: Create config file locally

**fluent-bit.conf**

```
# FireLens Example: Prometheus forward metrics - using the Fluent Bit image with init tag

# Scape node metrics every 20 seconds (use any metrics input plugin)
[INPUT]
    Name                        node_exporter_metrics
    Tag                         node_metrics
    Scrape_interval             20

# Send metrics to AMP via Remote Write
[OUTPUT]
    Name prometheus_remote_write
    Match node_metrics
    Host aps-workspaces.us-west-2.amazonaws.com
    Port 443
    Uri /workspaces/< my-amp-workspace-id >/api/v1/remote_write
    AWS_Auth On
    AWS_region us-west-2
    Tls On
    Tls.verify On
    add_label app my-ecs-app
    add_label color blue
```

**Note:** you can find this config file in the `config-files` directory of this example, please modify according to the actual situation to match your needs.

### Step 2: Create an Amazon Managed Service for Prometheus Workspace

* create an Amazon Managed Service for Prometheus workspace `my-prometheus` to store metrics
* update above config file's `prometheus_remote_write` `uri` to the created workspace's remote write url endpoint

### Step 3: Upload config file to S3

* create the S3 bucket `your-bucket` to store config files
* upload above config file to this bucket

### Step 4: Create the ECS Task

* create the ECS Task using provided `task-definition.json`, which using the Fluent Bit image with init tag
* change the `taskRoleArn` and `executionRoleArn` to your own role ARN
* change the `environment` part in the task definition FireLens configuration, copy the ARN of config files and paste it as environment variable's value. The name of environment variable requires to use the prefix `aws_fluent_bit_init_s3_`

### Step 5: View the metrics with Grafana

To view the metrics sent to Prometheus, consider creating an Amazon Managed Grafana workspace and adding the Amazon Managed Service for Prometheus workspace as it's data source. Refer to following [documentation](https://docs.aws.amazon.com/grafana/latest/userguide/AMP-adding-AWS-config.html) for more information.
