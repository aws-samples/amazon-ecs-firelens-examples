FROM public.ecr.aws/aws-observability/aws-for-fluent-bit:stable

ADD extra.conf /fluent-bit/alt/fluent-bit.conf
CMD echo -n "AWS for Fluent Bit Container Image Version " && \
    cat /AWS_FOR_FLUENT_BIT_VERSION && echo "" && \
    exec /fluent-bit/bin/fluent-bit -e /fluent-bit/firehose.so -e /fluent-bit/cloudwatch.so -e /fluent-bit/kinesis.so -c /fluent-bit/alt/fluent-bit.conf