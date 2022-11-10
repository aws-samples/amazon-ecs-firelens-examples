# FireLens Example: Quickly Enable Resource Monitoring with ADOT

The [AWS Distro for OpenTelemetry](https://aws-otel.github.io/) has many features for collecting telemetry from your applications. In this tutorial, you will learn how to quickly enable resource utilization monitoring for Fluent Bit using ADOT. OpenTelemetry can do a lot more than is shown in this tutorial. 

This guide is broken down into two tutorials:
- [Tutorial 1: Quickly Enable Task Resource Monitoring](#tutorial-1-quickly-enable-task-resource-monitoring)
- [Tutorial 2: Enable Container Level Metric Collection](#tutorial-2-enable-container-level-metric-collection)

### Tutorial 1: Quickly Enable Task Resource Monitoring

To quickly enable CPU, memory, network and disk usage metrics for your FireLens task, simply add the following container definition to your current task definition:

```
        {
            "essential": true,
            "image": "public.ecr.aws/aws-observability/aws-otel-collector:latest",
            "name": "aws-otel-collector",
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "aws-otel-collector",
                    "awslogs-region": "us-west-1",
                    "awslogs-create-group": "true",
                    "awslogs-stream-prefix": "ecs-"
                }
            },
            "memoryReservation": 100,
            "command": [
                "--config=/etc/ecs/container-insights/otel-task-metrics-config.yaml"
            ]
        }
```

This uses the [built-in ADOT configuration](https://github.com/aws-observability/aws-otel-collector/blob/main/config/ecs/container-insights/otel-task-metrics-config.yaml) for ECS Task Level Metrics. Please read the full [documentation](https://aws-otel.github.io/docs/components/ecs-metrics-receiver).

This configuration will send metrics to CloudWatch with the following sets of dimensions:
- `ClusterName`
- `ClusterName` and `TaskDefinitionFamily`

This means that you will get CPU, Memory, Disk, and Network metrics for the entire task aggregated across your cluster, and the union of cluster and TaskDefinition Family.

### Tutorial 2: Enable Container Level Metric Collection

For this tutorial, you must do the following:

1. Copy the "Full configuration for task- and container-level metrics" example from the [ADOT ECS Documentation](https://aws-otel.github.io/docs/components/ecs-metrics-receiver). 
2. Optionally modify the dimensions and metrics sent to CloudWatch with the `awsemf` exporter. For example, you could add `TaskId` and `ContainerName` as dimensions to monitoring specific Fluent Bit containers. However, please understand that this will lead to a very large number of individual metrics in CloudWatch which will induce cost. 
3. Follow the [ADOT ECS Custom Config tutorial](https://aws-otel.github.io/docs/setup/ecs/config-through-ssm) to run ADOT with your custom configuration.