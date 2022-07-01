#!/bin/bash

/go/filebeat/filebeat -c /go/filebeat.yml -e &
/go/jarvis-trading-bot server

