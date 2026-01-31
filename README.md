
# Wait2Calm

This is a simple utility that waits until the computer calms down. 

* You can use this before starting a CPU-hungry application.
* You can add random delays, it is useful when multiple `wait2calm` instances are running in parallel.

This is how it works:

1. wait a random amount of time, maximum wait is given by `--random-wait-before`
2. then start a loop. In each iteration:
   1. measure the load average (1, 5 or 15 minute, specified by `--load-type`)
   2. if the value is below `--immediate-start-below` then break the loop
   3. if the value is below `--delayed-start-below` in two consecutive iterations, then break the loop
   4. if the loop stated more than `--do-not-wait-after` ago, then break the cycle
   5. wait `--measurement-interval` before starting the next cycle
3. wait a random amount of time, maximum wait is given by `--random-wait-after`

Note: `--do-not-wait-after` also affects step 1 and step 3.

## Command line options

```
Usage:
  wait2calm [OPTIONS]

Application Options:
  -v, --verbose                Show verbose information
  -d, --debug                  Show debug information
      --version                Show version information and exit
  -b, --random-wait-before=    Max. random wait before first measurement, default is 15s (default: 10s)
  -a, --random-wait-after=     Max. random wait after first measurement, default is 15s (default: 10s)
  -t, --load-type=             Load average type, can be 1=load1,5=load5,15=load15 (default: 1)
  -m, --measurement-interval=  Interval between measurements, default is 10s (default: 10s)
  -i, --immediate-start-below= Immediately start below this load value
  -l, --delayed-start-below=   Start if two subsequent loads are below this value
      --do-not-wait-after=     Do not wait inside the loop after this amount of time, regardless of the load average. A non-positive value disables this function.
                               Resolution is 100msec.
      --success-on-timeout     Return with zero exit code, even if timed out on --do-not-wait-after

Help Options:
  -h, --help                   Show this help message

Usage:
  wait2calm [OPTIONS]

Application Options:
  -v, --verbose                Show verbose information
  -d, --debug                  Show debug information
      --version                Show version information and exit
  -b, --random-wait-before=    Max. random wait before first measurement, default is 15s (default: 10s)
  -a, --random-wait-after=     Max. random wait after first measurement, default is 15s (default: 10s)
  -t, --load-type=             Load average type, can be 1=load1,5=load5,15=load15 (default: 1)
  -m, --measurement-interval=  Interval between measurements, default is 10s (default: 10s)
  -i, --immediate-start-below= Immediately start below this load value
  -l, --delayed-start-below=   Start if two subsequent loads are below this value
      --do-not-wait-after=     Do not wait inside the loop after this amount of time, regardless of the load average. A non-positive value disables this function.
                               Resolution is 100msec.
      --success-on-timeout     Return with zero exit code, even if timed out on --do-not-wait-after

Help Options:
  -h, --help                   Show this help message

```

## Exit codes

* 0 - calmed down, or timed out with `--success-on-timeout`
* 1 - timed out on `--do-not-wait-after`
* 2 - other error

## Examples

```bash
wait2calm -v -b 10s -i 0.1 -l 1 --do-not-wait-after 5s-v -b 5s -i 0.1 -l 1 --do-not-wait-after 20s


Jan 31 16:01:45.574 INF Wait before first measurement max=5s actual=2.611179969s
Jan 31 16:01:48.194 INF Delayed start (first measurement) load=0.7 delayed-start-below=1
Jan 31 16:01:48.194 INF Wait before next measurement measurement-interval=10s elapsed=2.620105823s
Jan 31 16:01:58.233 INF Delayed start (second measurement) load=0.6 delayed-start-below=1
Jan 31 16:01:58.234 INF Wait after calm down max=10s actual=1.231541683s
Jan 31 16:01:59.472 INF Calmed down!
```

```bash
-v -b 5s -i 0.1 -l 1 --do-not-wait-after 10s

Jan 31 16:03:14.745 INF Wait before first measurement max=5s actual=2.408857066s
Jan 31 16:03:17.160 INF Do not start load=1.24 delayed-start-below=1
Jan 31 16:03:17.160 INF Wait before next measurement measurement-interval=10s elapsed=2.414654734s
Jan 31 16:03:24.784 WRN Do not wait anymore --do-not-wait-after=10s elapsed=10.039400089s
```
## Building

Tested on linux only. Does not work on linux (because github.com/shirou/gopsutil cannot get load average on Windows).

If you really want to use it under Windows, please contact me, cpu load could possibly be used instead of load average.

```bash
scripts/build.sh # to build for your platform
scripts/build.sh # to build for multiple OS and architectures
```
