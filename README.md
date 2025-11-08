# Parkside Notifier

A Telegram bot that scrapes the Portuguese Lidl website to check the availability of Parkside products. It uses the GPT-4o Mini model to analyze PDFs and extract information about Parkside items. 
You can join by accessing this [link](https://t.me/parksideNotifications).

## How to use

There's a docker-compose.yml file that creates a rqlite database and a go service that scrapes the website once a day. Every Parkside flyer is saved in the database to avoid duplicate messages being sent to the channel.

1. Create an .env file following the .env.example file
2. Replace __OPENAI_API_KEY__ with a valid open ai key that can be created in the [OpenAI dashboard](https://platform.openai.com/api-keys)

3. Start the containers
> docker compose build && docker-compose up -d

4. Create the table

```
curl --location 'http://127.0.0.1:4001/db/execute?pretty=null&timings=null' \
--header 'Content-Type: application/json' \
--data '[
    "CREATE TABLE message (url TEXT UNIQUE, notified INTEGER DEFAULT 0 NOT NULL)"
]'
```

5. Every day, the cron job will scrape the website and every new parkside sale will be sent to the channel

## Local development

1. On the docker-compose.yml, comment the go lang service and run

#### Podman
> podman compose up

#### Docker
> docker compose up


2. Create the database by call this endpoint

```
curl --location 'http://127.0.0.1:4001/db/execute?pretty=null&timings=null' \
--header 'Content-Type: application/json' \
--data '[
    "CREATE TABLE message (url TEXT UNIQUE, notified INTEGER DEFAULT 0 NOT NULL)"
]'
```

3. Create an .env file like the .env.example one

4. Run in the root folder
> go run ./src

