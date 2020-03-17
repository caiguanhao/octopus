octopus
=======

Interact Octopus Card Reader/Writer with JSON-RPC (Golang). (Not all octopus commands are implemented.)

## Commands

Examples of requests and responses:

```
# this connects to server
# nc 127.0.0.1 12345

req: {"id": 1, "method": "Octopus.Init", "params": [{"PortNumber": 0, "BaudRate": 115200, "ControllerID": 0}]}
res: (see "Octopus.Inspect" below)

req: {"id": 3, "method": "Octopus.Inspect", "params": []}
res: {
       "id": 3,
       "result": {
         "DeviceID": 9999999,
         "OperatorID": 11111111,
         "DeviceTime": "2020-01-01T00:00:00+08:00",
         "CompanyID": 111,
         "KeyVersion": 0,
         "EODVersion": 1111,
         "BlacklistVersion": 0,
         "FirmwareVersion": 8389448,
         "CCHSVer": 0,
         "LocationID": 0,
         "InterimBlacklistVersion": 0,
         "FunctionalBlacklistVersion": 0
       },
       "error": null
     }

req: {"id": 2, "method": "Octopus.UpdateLocationID", "params": [{"LocationID": 234567}]}
res: {"id": 2, "result": true, "error": null}

req: {"id": 4, "method": "Octopus.Poll", "params": [{"Command": 0, "Timeout": 5}]}
req: {"id": 4, "method": "Octopus.Poll", "params": [{"Command": 1, "Timeout": 5}]}
req: {"id": 4, "method": "Octopus.Poll", "params": [{"Command": 2, "Timeout": 5}]}
res: {
       "id": 4,
       "result": {
         "CardID": "11111111",
         "RemainingValue": 2222,
         "AddValueDetail": "0",
         "LastAddValueDevice": "3333333",
         "Class": "4",
         "Language": "0",
         "AvailableAutopayAmount": "5",
         "UniqueManufactureID": "0101010101010101",
         "Logs": [
           {
             "ServiceProviderID": "180",
             "TransactionAmount": "-1",
             "TransactionTime": "2020-01-01T00:00:00+08:00",
             "MachineID": "11111",
             "ServiceInfo": "0000000000"
           },
           ...
         ]
       },
       "error": null
     }

req: {"id": 3, "method": "Octopus.GetLastAddValueInfo", "params": []}
res: {"id": 3, "result": {"Date": "2019-12-26", "TypeCode": "4", "Type": "AAVS", "DeviceID": "56FFC0"}, "error": null}

req: {"id": 5, "method": "Octopus.Deduct", "params": [{"Value": 1, "ServiceInfo": "0F0C62D1F7", "DeferReleaseFlag": 1}]}
res: {"id": 5, "result": {"RemainingValue": 2000, "AdditionalInfo": "000000000000000000000000000000000000000000000000000000000000"}, "error": null}

req: {"id": 6, "method": "Octopus.GenerateExchangeFile", "params": []}
res: {"id": 6, "result": {"FileName": "MPS.FF30.20200101000000", "WarningCode": 0}, "error": null}

req: {"id": 7, "method": "Octopus.TxnAmt", "params": [{"Value": 0, "RemainingValue": 0, "LED": 0, "Sound": 1}]}
res: {"id": 7, "result": true, "error": null}
```

To compile, run `make` on a 32-bit Linux system.

The `octopus` program can also be run on 64-bit system.

It is recommended to build this project using [go1.8.7.linux-386.tar.gz](https://dl.google.com/go/go1.8.7.linux-386.tar.gz).

The header file `rwl_exp.h` and library file `librwl.a` are from Octopus company.

You might try to run the command as **root** if you're having trouble initializing Octopus.

The config file `RWL.INI` (case sensitive) should be in the same directory of `octopus`.

Set `PORTTYPE=USB` in `RWL.INI` if you are using USB to serial adapter.

If you change the content of `RWL.INI`, remember to restart the `octopus` program.

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
