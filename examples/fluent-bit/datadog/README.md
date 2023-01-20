### FireLens Example: Logging to Datadog with Fluent Bit

For documentation on sending your FireLens monitored log data to Datadog Logs, see: [Fluent Bit and Firelens](https://docs.datadoghq.com/integrations/ecs_fargate/#fluent-bit-and-firelens).

It should be noted that this example and Datadog docs show `"enable-ecs-log-metadata":"true"` (which is the default). This option tells FireLens to add ECS Task Metadata keys to logs. However, the Datadog output [intentionally](https://github.com/fluent/fluent-bit/blob/v1.9.10/plugins/out_datadog/datadog.c#L254) will not send the metadata in the log and will instead convert the metadata to DataDog tags for easy searching and integration with other DataDog observability features. 

For all configuration parameters for Fluent Bit DataDog output plugin, see [DataDog](https://docs.datadoghq.com/integrations/fluentbit/#configuration-parameters) or [Fluent Bit](https://docs.fluentbit.io/manual/output/datadog) documentation.

AWS recommends that you store sensitive information, like your Datadog API Key using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. This is optional; it is also valid to simply specify the API Key in options map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
	   "Name": "datadog",
	   "Host": "http-intake.logs.datadoghq.com",
	   "TLS": "on",
	   "apikey": "<DATADOG_API_KEY>",
	   "dd_service": "my-httpd-service",
	   "dd_source": "httpd",
	   "dd_tags": "project:example",
	   "provider": "ecs",
	   "retry_limit": "2"
   }
},
```
