# New Relic Logs Firelens Example

For more detailed documentation about configuring Firelens for New Relic logs please refer to: [Enable New Relic Logs for AWS Firelens](https://docs.newrelic.com/docs/logs/new-relic-logs/enable-logs/enable-new-relic-logs-aws-firelens)

## Images for Log Router Container
New Relic uses a custom Fluent Bit plugin and has provided custom images for US/EU regions

| AWS Region   | Full Image Name                                                                     |
|--------------|-------------------------------------------------------------------------------------|
| us-east-1    | 533243300146.dkr.ecr.us-east-1.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| us-east-2    | 533243300146.dkr.ecr.us-east-2.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| us-west-1    | 533243300146.dkr.ecr.us-east-1.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| us-west-2    | 533243300146.dkr.ecr.us-west-2.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| ca-cental-1  | 533243300146.dkr.ecr.ca-central-1.amazonaws.com/newrelic/logging-firelens-fluentbit |
| eu-central-1 | 533243300146.dkr.ecr.eu-central-1.amazonaws.com/newrelic/logging-firelens-fluentbit |
| eu-west-1    | 533243300146.dkr.ecr.eu-west-1.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| eu-west-2    | 533243300146.dkr.ecr.eu-west-2.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| eu-west-3    | 533243300146.dkr.ecr.eu-west-3.amazonaws.com/newrelic/logging-firelens-fluentbit    |
| eu-north-1   | 533243300146.dkr.ecr.eu-north-1.amazonaws.com/newrelic/logging-firelens-fluentbit   |

## Application Container logConfiguration
AWS recommends that you store your New Relic Insights Insert Key with the [AWS Secrets Manager](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data.html) as shown below:

```
"logConfiguration": {
    "logDriver":"awsfirelens",
    "options": {
        "Name": "newrelic" 
    },         
    "secretOptions": [{
        "name": "apiKey",
        "valueFrom": "arn:aws:secretsmanager:region:aws_account_id:secret:secret_name-AbCdEf"
    }]     
}

```

You can also alternatively add your New Relic Insights Insert Key directly into the logConfiguration block as shown below:
```
"logConfiguration": {
    "logDriver":"awsfirelens",
    "options": {
       	"Name": "newrelic"
		“apiKey”: “[YOUR_INSIGHTS_INSERT_KEY]”
    }  
 }

```