# Ochoku

Contraction-combination-wordcarwreck of Chocolate and Otaku, with a
delicious pun, _O_ prefix means Great, _choku_ is a contraction of
chocolate in Japanese.

# What does it do?!?

(see check list as a reference) Maximize your chocolate obsession by
focusing on your favorite chocolates!

Curate your treasure trove of chocolate knowledge, and make your
(twitter) friends green with salivating envy over your eurdite
and beautifully photographed chocolate collection.

# Setting up...

You'll need a server...

- Consul for service discovery
- mariadb for data storage
- Server space for images (or possibly an AWS/S3 account. If we get that far!)
- A Twitter Account

build with:

    make

Setup and boot services:

    ochoku setup

Start services:

    ochoku start

But: Until that all works use something like this:

    ./bin/ochoku start image-service
    ./bin/ochoku start auth-service
    ./bin/ochoku start user-service
    ./bin/ochoku start db-service
    ./bin/ochoku start chocolate-service
    ./bin/ochoku start ui-service
    ./bin/ochoku start twitter-service

# Usage

Point your browser at http://localhost:9001 start adding mother flippin' chocolates, awww yiss.
