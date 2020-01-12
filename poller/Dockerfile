FROM python:3.7

ADD requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

ADD poller.py .

CMD ["python", "poller.py"]
