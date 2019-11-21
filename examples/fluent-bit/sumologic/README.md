### FireLens Example: Logging to Sumologic with Fluent Bit

For documentation on sending your FireLens monitored log data to Sumologic, see: [Collect AWS ECS Fargate and EC2 Container Logs](https://help.sumologic.com/03Send-Data/Collect-from-Other-Data-Sources/AWS_Fargate_log_collection).

AWS recommends that you store sensitive information (like your Sumologic URI) using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html) as shown in the example [task definition](task-definition.json). This is optional - it is also valid to simply specify the URI in options map:

```
"logConfiguration": {
	"logDriver":"awsfirelens",
	"options": {
		"Name":"http",
		"Host": "<endpoint>",
		"URI": "/receiver/v1/http/<token>",
		"Port": "443",
		"tls": "on",
		"tls.verify": "off",
		"Format": "json_lines"
	}
}
```
