### FireLens Example: Sending Logs to hosted Logstash with Fluent Bit

You can use Fluent Bit's [HTTP output plugin](https://docs.fluentbit.io/manual/pipeline/outputs/http) to send your container's log to external or aws hosted logstash. 

To know more about Logstash, please refer [here](https://www.elastic.co/logstash) and [here](https://aws.amazon.com/elasticsearch-service/the-elk-stack/logstash/).

Logstash works with different plugins, for logs to be processed and transformed, you'll have to enable http input plugins.

An input plugin enables a specific source of events to be read by Logstash. To read more about input plugin, please refer [here](https://www.elastic.co/guide/en/logstash/current/input-plugins.html).

Logstash configuration example to enable source of events to be read by Logstash, here we have used port 8090 - 

```
input {
    beats {
        port => 5044
        client_inactivity_timeout => 500
    }
    http {
        port => 8080
        type => "elb_healthcheck"
    }
    http {
        type => 8090
        type => "ecs_fluent_bit_logs"
    }
}
```

AWS recommands that you store sensitive information (like the URI containing some sensative token) using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html), as shown in this example [task definition](task-definition.json). However using URI is optional and it is also valid to use URI as part of `options` map (URI is optional as part of http output plugin of Fluent Bit):

```
"logConfiguration": {
        "logDriver": "awsfirelens",
        "options": {
                "Name": "http",
                "Host": "api.logstash.fake.domain",
                "URI": "/some/<token>/tag/<tag>/",
                "Port": "8090",
                "Format": "json",
        }
}
```