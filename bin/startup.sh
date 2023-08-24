#!/bin/sh

HOME="/data/dcw-regulustar"
ENV="prod"

export DUBBO_GO_CONFIG_PATH="$HOME/conf/dubbogo-$ENV.yaml"
export DALINK_GO_CONFIG_PATH="$HOME/conf/param-$ENV.yaml"

nohup ./regulustar-linux &
