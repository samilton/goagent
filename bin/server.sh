#!/bin/bash

APP_HOME=$(pwd -P)

usage() {
	echo "${0} [start|stop|restart]"
}

pid() {
	PID=$(ps ax | grep peagent | grep -v grep | awk '{print $1}')
}

start() {
	pid
	if [ ${PID} ]; then
		echo "PE Agent is already running."
	else
		${APP_HOME}/peagent 2>&1 >> ${APP_HOME}/output.log &
	fi
}

stop() {
	pid
	if [ -z ${PID} ]; then
		echo "PE Agent doesn't appear to be running"
	else
		echo "Stopping PE Agent [ ${PID} ]"
		kill -9 ${PID}
	fi

}

case "$1" in
	start)
		start
		;;
	stop)
		stop
		;;
	restart)
		stop
		sleep 3
		start
		;;
	*)
		usage
		;;
esac
