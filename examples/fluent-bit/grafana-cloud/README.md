### FireLens Example: Logging to Grafana Loki with Fluent Bit

For documentation on sending your FireLens monitored log data to Grafana Loki, see: [Fluent Bit loki Plugin Documentation](https://docs.fluentbit.io/manual/v/1.9-pre/pipeline/outputs/loki). Be aware there is a separate Golang output plugin provided by [Grafana](https://grafana.com/docs/loki/latest/clients/fluentbit/) with different configuration options. Here is a [blog](https://calyptia.com/blog/how-to-send-logs-to-loki-using-fluent-bit) on sending logs to Loki using the loki Fluent Bit plugin. The Loki data source provides access to Loki, Grafana’s log aggregation system.

It should be noted that this example and Grafana docs show `"enable-ecs-log-metadata":"true"` (which is the default). This option tells FireLens to add ECS Task Metadata keys to logs.  

For all configuration parameters for Fluent Bit Loki output plugin, see [Fluent Bit loki Plugin Documentation](https://docs.fluentbit.io/manual/v/1.9-pre/pipeline/outputs/loki) documentation.

If you are looking for `"bearer_token"` support, please use the docker hub upstream [image](https://hub.docker.com/r/fluent/fluent-bit) and follow the [Fluent Bit loki latest Documentation](https://docs.fluentbit.io/manual/pipeline/outputs/loki).

AWS recommends that you store sensitive information, like your Datadog API Key using secretOptions as shown in the example Task Definition. This is optional; it is also valid to simply specify the Https password in options map:

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
    "http_passwd": "<HTTP Password>"
   }
},
```
