# Ioio

The goal is to make a simulator of networked infrastructure, so that people can visually see different failure modes.

## Nodes

Nodes are groupings of cores that execute requests.

## Cores

Cores can execute requests in parallel.

## Network connections

Edges joining different nodes, it's only possible to make requests along those.

## Requests

Requests are entities that encode programs, with the following two components:

- nodes to be visited
- amount of instructions required

Fore example a request for a basic website render could be encoded in the following yaml. *s* stands for sacond, an imaginary unit of time, or tick. 1000*ms* is 1*s* etc.

```yml
origin:
  ctl: [ o:origin ]
  prg: 
  - p/10:
    - get:
      ctl: [ h:lb, t:10s ]
      prg: @lb

lb:
  - lb/fair:
    - @web
  
web:
  - get
    ctl: [ h:web, t:10s, r:backoff/3 ]
    prg:
    - r/10: [ c/100ms, get/db ]

db:
  ctl: [ h:db, t:1t ]
  prg: [ c/10ms ]
```

This reads as: send 10 parallel requests to web box (or tier), with 10s timeout and 3 retries (with backoff), repeating 10 times a 100ms compute followed by a call to db box, which in turn does 10ms of compute with 1s timeout and no retries.