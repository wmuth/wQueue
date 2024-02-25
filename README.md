# wQueue
A queueing system written in Go

Using mostly the go standard library except GORM for database interaction and gorilla for some web toolkit functionality.

PostgreSQL acts as the backend.

# To run
- Make sure you have a .env file, check `example.env` for formatting
- ```go run main.go``` or ```go build && ./wQueue```

# Features
## Guests can
- List available queues
- See queuers in the queues
- Look at some help/about information

## Users can
- Do everything guests can
- Join queues
- Leave queues

## Admins can
- Do everything users can
- Manage queues they are in charge of
- Open/Close queues
- Change the displayed queue message
- Send an alert message to the queuers
- Remove queuers from the queue
