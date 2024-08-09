run:
	go run cmd/main.go

mutex:
	go run cmd/mutex/mutex.go

semaphore:
	go run cmd/semaphore/semaphore.go

message:
	go run cmd/message-exchange/message.go

gen:
	go run cmd/generator/generator.go

clean:
	rm -rf log.log
