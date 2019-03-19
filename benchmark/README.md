# build & deploy

```
./dev-cli service deploy ./benchmark/service
./dev-cli service start benchmark-service
cd ./benchmark/app
go build
```


# run

```
./benchmark/app/main &
./benchmark/run.sh
```

# results

```
>>>> Execute Task

Summary:
  Count:	20000
  Total:	8.34 s
  Slowest:	210.29 ms
  Fastest:	2.61 ms
  Average:	81.56 ms
  Requests/sec:	2397.54

Response time histogram:
  2.608 [1]	|
  23.377 [252]	|∎∎
  44.146 [1886]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  64.914 [4162]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  85.683 [5311]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  106.452 [4201]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  127.220 [2508]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  147.989 [1198]	|∎∎∎∎∎∎∎∎∎
  168.757 [355]	|∎∎∎
  189.526 [65]	|
  210.295 [32]	|

Latency distribution:
  10% in 43.02 ms
  25% in 58.44 ms
  50% in 79.26 ms
  75% in 101.46 ms
  90% in 123.40 ms
  95% in 136.00 ms
  99% in 159.81 ms
Status code distribution:
  [OK]        19971 responses
  [Unknown]   29 responses

Error distribution:
  [29]   rpc error: code = Unknown desc = Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?


>>>> Emit Event

Summary:
  Count:	20000
  Total:	2.35 s
  Slowest:	46.46 ms
  Fastest:	7.47 ms
  Average:	23.30 ms
  Requests/sec:	8497.43

Response time histogram:
  7.470 [1]	|
  11.369 [33]	|
  15.267 [142]	|∎
  19.166 [1290]	|∎∎∎∎∎
  23.064 [9949]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  26.963 [5497]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  30.862 [2048]	|∎∎∎∎∎∎∎∎
  34.760 [801]	|∎∎∎
  38.659 [197]	|∎
  42.557 [40]	|
  46.456 [2]	|

Latency distribution:
  10% in 19.59 ms
  25% in 20.88 ms
  50% in 22.60 ms
  75% in 24.80 ms
  90% in 28.70 ms
  95% in 31.01 ms
  99% in 35.25 ms
Status code distribution:
  [OK]   20000 responses
```
