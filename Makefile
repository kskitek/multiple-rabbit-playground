docker:
	docker run -d --rm -p 5672:5672 rabbitmq:alpine

run-forwarder:
	go run cmd/forwarder/forwarder.go -exchanges=2 -queues=2

run-consumer1:
	go run cmd/consumer/consumer.go -queue=1

run-consumer2:
	go run cmd/consumer/consumer.go -queue=2

