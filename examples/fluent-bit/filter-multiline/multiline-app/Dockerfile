FROM public.ecr.aws/amazonlinux/amazonlinux:latest
ADD test.log /test.log

RUN yum upgrade -y && yum install -y python3 pip3

RUN pip3 install boto3

WORKDIR /usr/local/bin

COPY main.py .

CMD ["python3", "main.py"]