### FireLens Example: Logging to Amazon OpenSearch Serverless with Fluent Bit

Amazon OpenSearch Serverless is new offering that eliminates your need to manage OpenSearch clusters. All existing Fluent Bit OpenSearch output plugin options work with OpenSearch Serverless. The only difference with serverless from a Fluent Bit POV is that you must specify the service name as `aoss` (Amazon OpenSearch Serverless) when you enable `AWS_Auth`:

```
AWS_Auth On
AWS_Region <aws-region>
AWS_Service_Name aoss
```

### Data Access Permissions
When sending logs to OpenSearch Serverless, your task role (e.g. `arn:aws:iam::XXXXXXXXXXXX:role/ecs_task_iam_role`) needs [OpenSearch Serverless Data Access permisions](https://docs.aws.amazon.com/opensearch-service/latest/developerguide/serverless-data-access.html). Give your task role the following Data Access permissions to your serverless collection:

```
aoss:CreateIndex
aoss:UpdateIndex
aoss:WriteDocument
```
No task role IAM policies are needed to access the collection.

### Adding Permissions with AWS CLI
To add Data Access permissions to your task role via AWS CLI, use the following command along with the `aoss-data-access-policy.json` file from the `amazon-opensearch-serverless` example folder. Be sure to update the `aoss-data-access-policy.json` document with your collection name, and task role arn.
```
aws opensearchserverless create-access-policy \
    --name log-write-policy \
    --type data \
    --policy  file://./aoss-data-access-policy.json
```
 Please note that the `opensearchserverless` command was introduced in aws cli v1 in [1.27.29](https://github.com/aws/aws-cli/blob/develop/CHANGELOG.rst#L473) and in aws cli v2 in [2.9.2](https://raw.githubusercontent.com/aws/aws-cli/v2/CHANGELOG.rst). If you recieve the error: `argument command: Invalid choice`, please [update your AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html).

### Additional Information
For more information on Fluent Bit & Amazon OpenSearch Serverless, see the [official Fluent Bit documentation](https://docs.fluentbit.io/manual/pipeline/outputs/opensearch).

Learn more about Data access control for Amazon OpenSearch Serverless in the [OpenSearch Service Developer Guide](https://docs.aws.amazon.com/opensearch-service/latest/developerguide/serverless-clients.html#serverless-ingestion-permissions).