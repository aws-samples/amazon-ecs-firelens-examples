FROM public.ecr.aws/amazonlinux/amazonlinux:latest
ADD houndofbaskerville.txt /houndofbaskerville.txt

RUN yum upgrade -y && yum install -y python3

WORKDIR /usr/local/bin

COPY main.py .

CMD ["python3", "main.py"]
