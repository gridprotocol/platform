#!/bin/bash

ps -ef | grep platform | grep -v 'color' | awk '{print $2}' | xargs kill -9

nohup ./platform daemon run -c dev > log 2>&1 &
