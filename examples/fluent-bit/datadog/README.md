### FireLens Example: Logging to Datadog with Fluent Bit

For documentation on sending your FireLens monitored log data to Datadog Logs, see: [Fluent Bit and Firelens](https://docs.datadoghq.com/integrations/ecs_fargate/#fluent-bit-and-firelens).

AWS recommends that you store sensitive information, like your Datadog API Key using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. This is optional; it is also valid to simply specify the API Key in options map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
	   "Name": "datadog",
	   "apiKey": "<DATADOG_API_KEY>",
	   "dd_service": "my-httpd-service",
	   "dd_source": "httpd",
	   "dd_tags": "project:example",
	   "provider": "ecs"
   }
},
```
