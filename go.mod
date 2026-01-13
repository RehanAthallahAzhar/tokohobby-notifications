module github.com/RehanAthallahAzhar/tokohobby-notifications

go 1.24.5

require (
	github.com/RehanAthallahAzhar/tokohobby-messaging-go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.5.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/RehanAthallahAzhar/tokohobby-messaging-go => ../messaging
