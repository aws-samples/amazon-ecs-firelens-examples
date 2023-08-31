### FireLens Example: Logging to Grafana Loki with Fluent Bit

For documentation on sending your FireLens monitored log data to Grafana Loki, see: [Fluent Bit and Firelens in the Grafana documentation] (https://grafana.com/docs/loki/latest/clients/fluentbit/). The Loki data source provides access to Loki, Grafana’s log aggregation system.

It should be noted that this example and Grafana docs show `"enable-ecs-log-metadata":"true"` (which is the default). This option tells FireLens to add ECS Task Metadata keys to logs.  

For all configuration parameters for Fluent Bit Grafana output plugin, see [Grafana](https://grafana.com/docs/loki/latest/clients/fluentbit/#configuration-options) or [Fluent Bit](https://docs.fluentbit.io/manual/v/1.9-pre/pipeline/outputs/loki) documentation.

AWS recommends that you store sensitive information, like your Grafana API Key using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. This is optional; it is also valid to simply specify the API Key in options map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
	   "Name": "grafana-loki",
	   "Url": "https://<userid>:<apikey>@logs-prod-us-west2.grafana.net/loki/api/v1/push",
	   "Labels": "{job=\"firelens\"}",
	   "RemoveKeys": "container_id,ecs_task_arn",
	   "LabelKeys": "container_name,ecs_task_definition,source,ecs_cluster",
	   "LineFormat": "key_value"
   }
},
```
