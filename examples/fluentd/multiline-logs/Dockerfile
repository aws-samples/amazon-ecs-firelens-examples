# Be sure to update the versions to the ones you require
FROM fluent/fluentd:v1.7.4-1.0

# Required if you're not storing this config on S3.
# In this case, you need to change the reference in the task definition as well.
# ADD fluentd.conf /extra.conf

# Use root account to install apk
USER root

# Below RUN installs the concat and cloudwatch-logs plugin.
# You may want to adjust these to the plugins you require
RUN apk add --no-cache --update --virtual .build-deps \
        sudo build-base ruby-dev \
 && sudo gem install fluent-plugin-concat \
 && sudo gem install fluent-plugin-cloudwatch-logs \
 && sudo gem sources --clear-all \
 && apk del .build-deps \
 && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

USER fluent
