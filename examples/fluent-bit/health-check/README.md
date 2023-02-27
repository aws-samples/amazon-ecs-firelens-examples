# Fluent Bit FireLens Container Health Check Guidance

This guide will help you choose the right health check for FireLens for your needs. One key decision to make along with the health check is whether the FireLens container definition will have the [essential field](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ContainerDefinition.html) set to `true` or `false`. 

As explained in the [ECS Health Check documentation](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_HealthCheck.html), if a container is essential, then if it becomes unhealthy the task will become unhealthy. All essential containers in the task must be healthy for the task to be healthy. 

If you make the FireLens container non-essential, its health status will still be displayed in the ECS Console and it will be returned in the output of [DescribeTasks](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeTasks.html). 

Please note, health checks are not the only option for monitoring Fluent Bit. Fluent Bit also has [built-in metrics](https://docs.fluentbit.io/manual/administration/monitoring). We have a [tutorial on how to send these metrics to CloudWatch](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/send-fb-internal-metrics-to-cw) in this repo. 

## Health Check Options

### Output Metric Based Health Check

Fluent Bit's [http monitoring interface has a health check option](https://docs.fluentbit.io/manual/administration/monitoring#health-check-for-fluent-bit). As explained in its documentation, the healthy or unhealthy status is based on the output plugin metrics- are the outputs successfully sending logs or not? This is a "deep" health check. You can configure the error threshold and the period of evaluate errors in the Fluent Bit configuration. 

See the [task-definition-output-metrics-healthcheck.json](task-definition-output-metrics-healthcheck.json) in this directory. This health check uses the built-in `/fluent-bit/configs/output-metrics-healthcheck.conf` configuration available in AWS for Fluent Bit 2.29.1+. This means you do not need to build a custom image to use this health check. If you still want to use a custom image built with your custom configuration, the `monitoring-server-healthcheck.conf` (and DIY section below) shows the configuration required and the [config-file-type](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file) example in this guide shows how to use a custom config file. 

*Benefits:*
* If Fluent Bit can not send logs to your destination- if it can not serve its business purpose- it will become unhealthy. 

*Drawbacks:*
* If inputs are blocked/paused or are not receiving logs, the output metrics will not be affected and will still look positive. Thus, this health check can not cover all possible failure scenarios.
* If the destination is down, then all of your Fluent Bit containers across your entire fleet will become unhealthy all at the same time. If the FireLens container is essential, all tasks in your service would become unhealthy and be replaced at the same time. 

If this option is appealing but one of the drawbacks is a deal-breaker, then you may want to monitor Fluent Bit's plugin metrics instead. The health check endpoint and options are just built on top of the Fluent Bit [built-in metrics](https://docs.fluentbit.io/manual/administration/monitoring). We have a [tutorial on how to send these metrics to CloudWatch](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/send-fb-internal-metrics-to-cw) in this repo. 

#### DIY setup

If you do not use the built-in configuration for this health check, you will need to add the following to your configuration file to use this healthcheck. This is the same the `monitoring-server-healthcheck.conf` included in this directory. 

```
[SERVICE]
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_PORT    2020
    # https://docs.fluentbit.io/manual/administration/monitoring#health-check-for-fluent-bit
    Health_Check On 
    # customize error and retry thresholds and evaluation period as desired
    HC_Errors_Count 5 
    HC_Retry_Failure_Count 5 
    HC_Period 5
```


### Simple Uptime Health Check

Fluent Bit's [monitoring interface](https://docs.fluentbit.io/manual/administration/monitoring#http-server) has a `/api/v1/uptime` path which can be queried for a simple "is it responsive" health check. This is a "shallow" health check. It is built on the assumption that if the monitoring uptime interface returns a response, then Fluent Bit is still running and is not frozen or unresponsive. The Fluent Bit monitoring interface is widely used by customers who use it to get prometheus format metrics from Fluent Bit, and the AWS for Fluent Bit team has not received many bug reports for this interface over the years. You can use the buit-in `/fluent-bit/configs/output-metrics-healthcheck.conf` configuration available in AWS for Fluent Bit 2.29.1+. While the config enables the output metrics healthcheck above, it also enables the monitoring server which is used by this simple uptime healthcheck. See the DIY section below for information on setting up this healthcheck via your own custom configuration.

See the [task-definition-uptime-healthcheck.json](task-definition-uptime-healthcheck.json) in this directory. 

*Benefits:*
* Unlike the "deeper" health check option above, if your destination is down, this will not cause the uptime health check to fail. 

*Drawbacks:*
* It is a very shallow health check, Fluent Bit could be completely failing to send logs but as long as its still responsive on the monitoring interface, it will be marked as healthy.
* AWS for Fluent Bit team has only recommended/documented this health check option for users since April 2022, so if you are reading this close to that date, it has not yet been very widely tested with real production workloads.

#### DIY setup

If you do not use the built-in configuration for this health check, you will need to add the following to your configuration file to use this healthcheck:

```
[SERVICE]
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_PORT    2020
```

### [Not recommended] TCP Input Health Check

Finally, another health check option is to try to ingest logs into Fluent Bit, and have the health check validate that the logs are accepted. This can be done using the netcat utility and the [TCP Input](https://docs.fluentbit.io/manual/pipeline/inputs/tcp). Logs ingested for the health check can be discarded via the null output. This health check is an attempt at a "deep" health check that is based on the assumption that Fluent Bit must be responsive if it will accept new logs via the TCP input. If it is down or frozen or nonresponsive, the netcat command will fail.

See the [task-definition-tcp-healthcheck.json](task-definition-tcp-healthcheck.json) in this directory. This health check does not require any extra configuration because [FireLens](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/) automatically generates [the TCP input and null output for the health check](https://github.com/aws-samples/amazon-ecs-firelens-under-the-hood/blob/mainline/generated-configs/fluent-bit/generated_by_firelens.conf#L10). *It should be noted that the link in the last sentence is an example of a generated configuration for a specific task definition, [please read and understand it carefully](https://github.com/aws-samples/amazon-ecs-firelens-under-the-hood/tree/mainline/generated-configs/fluent-bit)*. If you would like to setup this health check yourself, and you do not use the auto-generated FireLens config (for example, because you followed the [How to set input parameters in FireLens blog](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/), see the DIY section below).

*Benefits:*
* Is a "deeper" health check that attempts to verify that Fluent Bit is still accepting logs. 

*Drawbacks:*
* AWS for Fluent Bit team has recommended this health check to users since the launch of FireLens at the end of 2019. We chose to build the config for the health check into FireLens to provide all users with an easy deep health check option. Since then, a minority of users have reported that the health check appears to be susceptible to false positives where Fluent Bit is marked as unhealthy despite continuing to function. Unfortunately, AWS for Fluent Bit team has never been able to independently reproduce or root-cause these failures. *However, because we have received multiple similar reports of bad behavior from this health check over time, we no longer recommend it.* 

#### DIY setup

If you do not use the built-in configuration for this health check, you will need to add the following to your configuration file to use this healthcheck:

```
# input to recieve health check logs
[INPUT]
    Name tcp
    Tag firelens-healthcheck
    Listen 127.0.0.1
    Port 8877

# null output so that health check logs do not go anywhere
[OUTPUT]
    Name null
    Match firelens-healthcheck
```

## Recommended settings for new users

If you are new to FireLens, we recommend setting the FireLens container's essential value to `true`. This way, if there is an issue with your logging config and Fluent Bit can't start, or some it fails after starting, you task will be stopped- logging failures will be easy to notice/hard to miss. 

If you choose to make the FireLens container non-essential, you must then set up some sort of separate monitoring system to check it and notify you of failures. Because if the FireLens container is non-essential then it can fail to start or crash without impacting your task. Your app would keep running but no logs would be sent, so you would be flying blind. 

For health check, we recommend new users either do not enable a health check or choose the simple uptime health check. 

