octopus
=======

Interact Octopus Card Reader/Writer with JSON-RPC (Golang).

## Commands

(Not all commands are implemented)

```
# this connects to server
# nc 127.0.0.1 1234

{"id": 1, "method": "Octopus.Init", "params": [{"PortNumber": 0, "BaudRate": 115200, "ControllerID": 0}]}

{"id": 2, "method": "Octopus.UpdateLocationID", "params": [{"LocationID": 234567}]}

{"id": 3, "method": "Octopus.Inspect", "params": []}

{"id": 4, "method": "Octopus.Poll", "params": [{"Command": 0, "Timeout": 5}]}
{"id": 4, "method": "Octopus.Poll", "params": [{"Command": 1, "Timeout": 5}]}
{"id": 4, "method": "Octopus.Poll", "params": [{"Command": 2, "Timeout": 5}]}

{"id": 5, "method": "Octopus.Deduct", "params": [{"Value": 1, "ServiceInfo": "0F0C62D1F7", "DeferReleaseFlag": 1}]}

{"id": 6, "method": "Octopus.GenerateExchangeFile", "params": []}
```

To compile, run `make` on a 32-bit Linux system.

The `octopus` program can also be run on 64-bit system.

The header file `rwl_exp.h` and library file `librwl.a` are from Octopus company.

---

LICENSE: MIT

Author: caiguanhao
