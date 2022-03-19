FROM python:3.7

RUN pip install flask

ADD ./main.py .
ADD ./certs/server.key .
ADD ./certs/server.crt .
ADD ./certs/ca.crt .


CMD ["python", "-u", "main.py"]