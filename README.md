# Jobsity Chat Robot

This repository contains a decoupled service that parses a CSV file obtained from [stooq.com](https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv).

## RabbitMQ

This service is powered by RabbitMQ, and the configuration settings have been preconfigured in the Makefile.

## Running

You can start the service with the command:

```bash
make run
