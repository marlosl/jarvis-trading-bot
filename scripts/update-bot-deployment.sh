#!/bin/bash

kubectl patch -n jarvis-trading-bot deployment candles-receiver-analyzer -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"
