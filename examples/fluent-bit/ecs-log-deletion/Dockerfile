# Use the standard Amazon Linux base
FROM amazonlinux

# Install Cron
RUN yum install -y cronie which findutils perl-core compat-libcap1 lsof tar shadow-utils setuptool procps && yum clean all

# Copy crontab that deletes logs
COPY cron/delete-log-files /etc/cron.d/delete-log-files

# Give execution rights on the cron job
RUN chmod 0644 /etc/cron.d/delete-log-files

