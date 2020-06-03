#!/bin/bash

# turn on bash's job control
set -m

# start daemon
mesg-daemon start &
daemon=$!

# function that stop child processes
function stop {
  echo "stopping daemon"
  kill $daemon
  wait $daemon
  exit 0
}

# trap both sigint (ctrl+c) and sigterm (OS ask process to be stopped)
trap stop SIGINT
trap stop SIGTERM

# wait 5 sec for the daemon to start rpc server
sleep 5 &
wait $!

# start lcd
echo "starting lcd"
mesg-cli rest-server --laddr tcp://0.0.0.0:1317 &
lcd=$!

# start orchestrator
echo "starting orchestrator"
mesg-cli orchestrator start &
orchestrator=$!

# this variable is used to control the monitoring of the child process
monitoring=true

# function that stop child processes
function stop {
  monitoring=false
  echo "stopping all child processes"
  kill $daemon $lcd $orchestrator
  wait $daemon $lcd $orchestrator
}

# trap both sigint (ctrl+c) and sigterm (OS ask process to be stopped)
trap stop SIGINT
trap stop SIGTERM

# start the monitoring loop
while $monitoring
do
  # check if all sub processes are running correctly
  if [[ -n "$(ps -p $daemon -o pid=)" ]] && [[ -n "$(ps -p $lcd -o pid=)" ]] && [[ -n "$(ps -p $orchestrator -o pid=)" ]]; then
    sleep 2 &
    wait $!
  else
    # if one child process is not running, stopping all of them and exit with error code 1
    echo "one child process is not running"
    stop
    exit 1
  fi
done
