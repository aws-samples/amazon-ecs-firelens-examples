### FireLens Example: Logging to Grafana Loki with Fluent Bit

For documentation on sending your FireLens monitored log data to Grafana Loki, see: [Fluent Bit and FireLens in the Grafana documentation](https://grafana.com/docs/loki/latest/clients/fluentbit/). Be aware there is a separate Golang output plugin provided by [Grafana](https://grafana.com/docs/loki/latest/clients/fluentbit/) with different configuration options. The Loki data source provides access to Loki, Grafana’s log aggregation system.

It should be noted that this example and Grafana docs show `"enable-ecs-log-metadata":"true"` (which is the default). This option tells FireLens to add ECS Task Metadata keys to logs.  

For all configuration parameters for Fluent Bit Grafana output plugin, see [Fluent Bit](https://docs.fluentbit.io/manual/v/1.9-pre/pipeline/outputs/loki) documentation.

AWS recommends that you store sensitive information, like your Grafana API Key using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example Task Definition. This is optional; it is also valid to simply specify the API Key in options map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
	   "Name": "loki",
	   "Host": "logs-prod-us-west2.grafana.net",
	   "port": "443",
	   "tls": "on",
	   "tls.verify": "on",
	   "http_user": "user_id",
	   "http_passwd": "password"
   }
},
```
