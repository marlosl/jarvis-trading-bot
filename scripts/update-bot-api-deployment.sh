#!/bin/bash

kubectl patch -n jarvis-trading-bot deployment trade-bot-api -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"