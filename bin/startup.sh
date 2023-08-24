#!/bin/sh

HOME="/data/aspw/server"
ENV="prod"

export DUBBO_GO_CONFIG_PATH="$HOME/conf/dubbogo-$ENV.yaml"
export DALINK_GO_CONFIG_PATH="$HOME/conf/param-$ENV.yaml"

nohup ./spw-linux &
