#!/bin/bash
set -e

# path to watch
WATCH_PATH=${WATCH_PATH:-"/tmp/gmi.storage.sh"}
# path to log
LOG_PATH=${LOG_PATH:-"/tmp/gmi.storage.log"}

# clear log file
echo "" > $LOG_PATH

current_version=""

# set exit flag
EXIT_FLAG=0

# log function: log [INFO|WARN|ERROR] "message"
log() {
  local level="${1:-INFO}"
  local message="$2"
  local timestamp=$(date "+%Y-%m-%d %H:%M:%S")
  
  # set color based on log level
  local color=""
  local reset="\033[0m"
  
  case "$level" in
    INFO)  color="\033[0;32m" ;; # 绿色
    WARN)  color="\033[0;33m" ;; # 黄色
    ERROR) color="\033[0;31m" ;; # 红色
    *)     color="\033[0m"    ;; # 默认
  esac
  
  # output to console (with color) and log file (without color)
  echo -e "${color}[${timestamp}] [${level}] ${message}${reset}" >> $LOG_PATH
}

# cleanup function
cleanup() {
  local signal=$1
  log INFO "received signal $signal, prepare to exit..."
  EXIT_FLAG=1
  
  # umount all mount points
  nsenter -t 1 -m -u -i -n -p -- "$WATCH_PATH" "umount" | tee -a $LOG_PATH || true
  
  # wait for background processes to finish
  wait 2>/dev/null || true
  
  log INFO "script terminated"
  exit 0
}

# set signal handler
trap 'cleanup SIGTERM' TERM
trap 'cleanup SIGINT' INT
trap 'cleanup SIGHUP' HUP

# handle file change
handle_file_change() {
  # check $WATCH_PATH version
  local version=$(md5sum "$WATCH_PATH" | awk '{print $1}')
  if [ "$version" = "$current_version" ]; then
    log INFO "file version not changed: $WATCH_PATH"
    return
  fi
  current_version=$version

  log INFO "check storage script changed: $WATCH_PATH"
  
  # check file has execute permission
  if [ ! -x "$WATCH_PATH" ]; then
    log INFO "add execute permission: $WATCH_PATH"
    chmod +x "$WATCH_PATH"
  fi
  # execute file
  log INFO "execute storage script: $WATCH_PATH"
  nsenter -t 1 -m -u -i -n -p -- chmod +x "$WATCH_PATH"
  nsenter -t 1 -m -u -i -n -p -- "$WATCH_PATH" "mount" | tee -a $LOG_PATH || true
  log INFO "storage script executed: $WATCH_PATH"
}

# start file watcher
start_file_watcher() {
  log INFO "start watch: $WATCH_PATH"
  # start continuous monitoring process
  inotifywait -m -r -e modify --format '%w%f' "$WATCH_PATH" | while read file; do
    if [ -f "$WATCH_PATH" ]; then
      handle_file_change "$WATCH_PATH"
    fi
  done &
}

# main function
main() {
  while [ $EXIT_FLAG -eq 0 ]; do
    if [ ! -f "$WATCH_PATH" ]; then
      log WARN "file not found: $WATCH_PATH"
      sleep 1
    else
      break
    fi
  done
  # check if exit signal received
  if [ $EXIT_FLAG -eq 1 ]; then
    log INFO "received exit signal, stop initialization"
    exit 0
  fi

  handle_file_change $WATCH_PATH
  start_file_watcher

  while [ $EXIT_FLAG -eq 0 ]; do
    log INFO "waiting for exit signal..."
    sleep 10
  done
}

main "$@"