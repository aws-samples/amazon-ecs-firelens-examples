### FireLens Example: Logging to Dynatrace with Fluent Bit

For documentation on sending your FireLens monitored log data to Dynatrace Logs, see: [Fluent Bit with Dynatrace](https://www.dynatrace.com/hub/detail/fluent-bit/).

AWS recommends that you store sensitive information, like your Dynatrace API Token using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example [Task Definition](https://github.com/dynatrace-oss-contrib/amazon-ecs-firelens-examples/blob/mainline/examples/fluent-bit/dynatrace/task-definition.json). This is optional; it is also valid to simply specify the API Token in the URI map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
	   "Name": "http",
	   "Host": "{your-environment-id}.live.dynatrace.com",
	   "TLS": "on",
	   "TLS.verify" : "off",
	   "Format": "json",
	   "Json_Date_Format": "iso8601",
	   "Json_Date_Key": "timestamp",
	   "Header: "Content-Type application/json; charset=utf-8",
	   "Port": "443",
	   "URI": "/api/v2/logs/ingest?api-token={your-API-Token-here}",
	   "Allow_Duplicated_Headers": "false",
	   "retry_limit": "2"
   }
},
```
