# executor

A little program executing command from shell, and prepending the command name to stdout &amp; stderr (with color support)

Tested in Linux only (as I'm not going to use it in Windows at least for now)

## Why?

When I am doing experiments that involves changing several parameters and comparing their results, mostly I choose to write a simple shell script to do that. Here is an example:

```sh
#!/bin/bash

# This is only an example

CNT=0
for X in `seq 1 10 100`; do
  for Y in 0.1 0.2 0.3; do
    PYTHONUNBUFFERED=TRUE ANOTHERENV=$(((CNT)%4)) python exp.py --X "$X" --Y "$Y" > log/exp-"$X"-"$Y".log &
    
    CNT=$((CNT+1))
  done;
done;
```

It is naive, but it works at most time.

However, sometimes there may be something wrong with your environment, or your experiment script, and these processes will complain everything directly to your `stderr`. And it is just so hard to find out whom says what.

And `executor` can help a little:

```sh
#!/bin/bash

# This is only an example

CNT=0
for X in `seq 1 10 100`; do
  for Y in 0.1 0.2 0.3; do
    PYTHONUNBUFFERED=TRUE ANOTHERENV=$(((CNT)%4)) ./executor python exp.py --X "$X" --Y "$Y" > log/exp-"$X"-"$Y".log &
    
    CNT=$((CNT+1))
  done;
done;
```

When things go wrong, you can figure out the cmdline and pid of the complaining process just from the output.

## How to build

```shell
$ go build
```

## Examples

```shell
$ ./executor ping
[ping]: Executing
[ping]: PID=3290
[ping] 3290 stderr: Usage: ping [-aAbBdDfhLnOqrRUvV] [-c count] [-i interval] [-I interface]
[ping] 3290 stderr:             [-m mark] [-M pmtudisc_option] [-l preload] [-p pattern] [-Q tos]
[ping] 3290 stderr:             [-s packetsize] [-S sndbuf] [-t ttl] [-T timestamp_option]
[ping] 3290 stderr:             [-w deadline] [-W timeout] [hop1 ...] destination
2020/12/01 16:31:34 [ping] 3290 exited with status 2
$ ./executor uname -a
[uname -a]: Executing
[uname -a]: PID=28629
[uname -a] 28629 stdout: Linux <redacted> 4.19.0-8-cloud-amd64 #1 SMP Debian 4.19.98-1 (2020-01-26) x86_64 GNU/Linux
$ ./executor -shell echo '$PWD'
[echo $PWD]: Executing
[echo $PWD]: PID=35979
[echo $PWD] 35979 stdout: /tmp
$ ./executor -shell 'echo $PWD'
[echo $PWD]: Executing
[echo $PWD]: PID=36854
[echo $PWD] 36854 stdout: /tmp
$ ./executor non-exist
[non-exist]: Executing
2020/12/01 16:42:54 [non-exist] cmd.Start(): exec: "non-exist": executable file not found in $PATH
$ ./executor -shell non-exist
[non-exist]: Executing
[non-exist]: PID=45538
[non-exist] 45538 stderr: sh: 1: non-exist: not found
2020/12/01 16:43:23 [non-exist] 45538 exited with status 127
```

---

`executor` is a practice of golang when I'm learning this language, and golang looks decent during this experience. At the beginning I would like to implement this in Python, but how can I get stdout/stderr real-time? Threading/Multiprocessing in Python just sounds so creepy (they have given me many "unforgettable experiences" before), and I'm not familiar with `asyncio`.

And goroutines just works here. It's nice.
