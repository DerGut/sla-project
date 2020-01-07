import json
import time

import pika
from pika import URLParameters

RABBIT_HOST = "queue"
RABBIT_QUEUE = "tasks"


def callback(ch, method, properties, body):
    task = json.loads(body)
    print("Received {}".format(task), flush=True)

    ch.basic_ack(delivery_tag=method.delivery_tag)


print("Sleeping for 30sec to wait for rabbit")
time.sleep(30)

connection = pika.BlockingConnection(
    parameters=URLParameters("amqp://guest:guest@{}:5672".format(RABBIT_HOST)))
channel = connection.channel()

print("Connected to rabbit")

channel.queue_declare(queue=RABBIT_QUEUE, durable=True)
channel.basic_qos(prefetch_count=1)
channel.basic_consume(queue=RABBIT_QUEUE, on_message_callback=callback)

try:
    channel.start_consuming()
except KeyboardInterrupt:
    channel.stop_consuming()
connection.close()
