#!/bin/sh

while true
do
  echo -ne "HTTP/1.1 200 OK\r\n\r\nok\r\n" | nc -l -p 80 -w 0
done
