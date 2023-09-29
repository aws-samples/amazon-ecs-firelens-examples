## FireLens Example: Forward metrics to Prometheus - using the Fluent Bit image with init tag

This example shows how to forward metrics with Fluent Bit's prometheus remote write output plugin to an Amazon Managed Service for Prometheus workspace.

### Step 1: Create an Amazon Managed Service for Prometheus Workspace

* Create an Amazon Managed Service for Prometheus workspace, `my-prometheus`, to store metrics. Refer to the [AMP onboarding docs](https://docs.aws.amazon.com/prometheus/latest/userguide/AMP-onboard-create-workspace.html) for more information.
* Take note of the `prometheus_remote_write` `uri` from the created workspace's remote write url endpoint, used in Fluent Bit's configuration in step 2.

### Step 2: Create config file locally

**fluent-bit.conf**

```
# FireLens Example: Prometheus forward metrics - using the Fluent Bit image with init tag

# Scrape node metrics every 20 seconds (collect metrics using any metrics input plugin)
# See the docs for more information: https://docs.fluentbit.io/manual/pipeline/inputs/node-exporter-metrics
[INPUT]
    Name                        node_exporter_metrics
    Tag                         node_metrics
    Scrape_interval             20

# Send metrics to AMP via Remote Write
# See the docs for more information: https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-remote-write
[OUTPUT]
    Name prometheus_remote_write
    Match node_metrics
    Host aps-workspaces.< region >.amazonaws.com
    Port 443
    Uri /workspaces/< my-amp-workspace-id >/api/v1/remote_write
    AWS_Auth On
    AWS_region us-west-2
    Tls On
    Tls.verify On
    add_label app my-ecs-app
    add_label color blue
```

**Note:** You can find this config file in the `config-files` directory of this example. Please modify the `node_exporter_metrics` input configuration options according to the metrics you desire to collect, along with the `prometheus_remote_write` output configuration options for `Host` and `Uri` to match your AMP setup.

### Step 3: Upload config file to S3

* Create the S3 bucket `your-bucket` to store config files
* Upload above config file to this bucket

### Step 4: Create the ECS Task

* Create the ECS Task using provided `task-definition.json`, which uses the Fluent Bit image with init tag
* Change the `taskRoleArn` to an IAM role that has policies listed in the `permissions.json` file found in this example directory.
* Change the `executionRoleArn` to an IAM role that has necessary permissions for launching an ECS task. See [ecs documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html) for more information.
* Change the `environment` section in the task definition FireLens configuration, copy the S3 ARN of config file and paste it as environment variable's value. The name of environment variable should be `aws_fluent_bit_init_s3_1`

### Step 5: View the metrics with Grafana

To view the metrics sent to Prometheus, consider creating an Amazon Managed Grafana workspace and adding the Amazon Managed Service for Prometheus workspace as it's data source. Refer to following [documentation](https://docs.aws.amazon.com/grafana/latest/userguide/AMP-adding-AWS-config.html) for more information.
