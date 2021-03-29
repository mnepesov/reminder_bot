## Prerequisites
* go 1.16
* docker
* Docker compose

### Create .env file in root directory and add following values:

````
TGToken=<Telegram bot token>
PostgresqlPassword=postgres
GoogleMapApiKey=<google map api key>
RabbitMQPassword=guest
````
## Libraries

* [github.com/spf13/viper](https://github.com/spf13/viper) - Yaml, Json, INI, and other formats are supported
* [github.com/jasonlvhit/gocron](https://github.com/jasonlvhit/gocron) - gocron allows you to run Go functions periodically at a predefined interval
* [github.com/streadway/amqp](https://github.com/streadway/amqp) - AMQP 0.9.1 client with RabbitMQ extensions in Go
* [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) - sqlx is a package for Go which provides a set of extensions on top of the excellent built-in database/sql package.
* [github.com/joho/godotenv](https://github.com/joho/godotenv) - loads env vars from a .env file
