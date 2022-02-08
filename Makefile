run:
	go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=postgres

test:
	go test -v ./...