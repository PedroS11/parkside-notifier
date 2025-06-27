# Parkside Notifier

Telegram bot that crawls portuguese Lidl website for Parkside products being available. You can join be accessing this [link](https://t.me/parksideNotifications).

## How to use

There's a docker-compose.yml file that creates a rqlite database and a go service that scrapes the website once a day. Every Parkside event is saved in the database to avoid duplicate messages being sent to the channel.

1. Start the containers
> docker compose build && docker-compose up -d

2. Create the table

```
curl --location 'http://127.0.0.1:4001/db/execute?pretty=null&timings=null' \
--header 'Content-Type: application/json' \
--data '[
    "CREATE TABLE message (url TEXT UNIQUE, notified INTEGER DEFAULT 0 NOT NULL)"
]'
```

3. Every day, the cron job will scrape the website and every new parkside sale will be sent to the channel

