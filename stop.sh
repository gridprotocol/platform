#!/bin/bash

ps -ef | grep platform | grep -v 'color' | awk '{print $2}' | xargs kill -9
