# We are choosing Go - Micro

# What do we have that we could exercise...

| Package                                                               | Built-in Plugin       | Description                                                       |
|:----------------------------------------------------------------------|:----------------------|:------------------------------------------------------------------|
| [auth](https://godoc.org/github.com/micro/go-platform/auth)           | auth-srv              | authentication and authorisation for users and services           |
| [config](https://godoc.org/github.com/micro/go-platform/config)       | config-srv            | dynamic configuration which is namespaced and versioned           |
| [db](https://godoc.org/github.com/micro/go-platform/db)               | db-srv                | distributed database abstraction                                  |
| [discovery](https://godoc.org/github.com/micro/go-platform/discovery) | discovery-srv         | extends the go-micro registry to add heartbeating, etc            |
| [event](https://godoc.org/github.com/micro/go-platform/event)         | event-srv             | platform event publication, subscription and aggregation          |
| [kv](https://godoc.org/github.com/micro/go-platform/kv)               | distributed in-memory | simply key value layered on memcached, etcd, consul               |
| [log](https://godoc.org/github.com/micro/go-platform/log)             | file                  | structured logging to stdout, logstash, fluentd, pubsub           |
| [monitor](https://godoc.org/github.com/micro/go-platform/monitor)     | monitor-srv           | add custom healthchecks measured with distributed systems in mind |
| [metrics](https://godoc.org/github.com/micro/go-platform/metrics)     | telegraf              | instrumentation and collation of counters                         |
| [router](https://godoc.org/github.com/micro/go-platform/router)       | router-srv            | global circuit breaking, load balancing, A/B testing              |
| [sync](https://godoc.org/github.com/micro/go-platform/sync)           | consul                | distributed locking, leadership election, etc                     |
| [trace](https://godoc.org/github.com/micro/go-platform/trace)         | trace-srv             | distributed tracing of request/response                           |

# Preferred modules...

- Discovery
- Config
- Log
- Router
- Auth
- Db
- Metrics
- Event

# Ideas for building?

We could build:

- Translator service
- Aphorisms service
- Fortunes/Quotes (or character's phrases etc.) service
- Email sender (which consumes fortures service)
- File sharing (dropbox etc.)
- Image sharing (Imgur etc.)
- Manga Reader/Tracker
- Metacritic Review site
- Captcha service
- Weather service
- ~~Open Air Replacement~~
- Jason F's Chocolate Addiction Product (choco-otaku @chocotweets)
    - Twitter Consumer (producer?) for creating geoloaction / review
    - Chocolate geolocation / pictures / review / rating

# Simple Spec of Jason F's Chocoloate Addiction Product

- [ ] you'll be able to use a web based (mobile friendly) user interface to write information about a yummy chocolate
- [ ] you'll be able to view a big board of yummy chocolate
- [ ] you'll be able to search for any yummy chocolates (you wrote about)
- [ ] you'll be able to visibly sort the big board of yummy chocolate
- [ ] Only you can change your treasure trove of yummy chocolates
- [ ] you'll be able to post a tweet for a yummy chocolate (you wrote about)
- [ ] Anyone can view my big board of yummy chocolate

### How does this decompose into services?


### Build plan!

- Green field, direct to micro-services YUP/YAY!
