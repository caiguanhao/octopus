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

{"id": 7, "method": "Octopus.TxnAmt", "params": [{"Value": 0, "RemainingValue": 0, "LED": 0, "Sound": 1}]}
```

To compile, run `make` on a 32-bit Linux system.

The `octopus` program can also be run on 64-bit system.

The header file `rwl_exp.h` and library file `librwl.a` are from Octopus company.

You might try to run the command as root user if you're having trouble initializing Octopus.

The config file `RWL.INI` (case sensitive) should be in the same directory of `octopus`.

## Example

```
node example.js

// tap card

{ id: 406988,
  result:
   { CardID: '12345678',
     RemainingValue: 49898,
     AddValueDetail: '0',
     LastAddValueDevice: '1234567',
     Class: '4',
     Language: '0',
     AvailableAutopayAmount: '20',
     UniqueManufactureID: '0101010101010101',
     Logs: null },
  error: null }
```

---

LICENSE: MIT

Author: caiguanhao
