package main

// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L. -lrwl
// #include <stdlib.h>
// #include <string.h>
// #define _C_
// #include "rwl_exp.h"
import "C"
import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	TransactionTimeSince int64 = 946684800 // 2000-01-01T00:00:00+00:00
)

type (
	Octopus struct{}

	CardReaderInfo struct {
		DeviceID                   int
		OperatorID                 int
		DeviceTime                 time.Time
		CompanyID                  int
		KeyVersion                 int
		EODVersion                 int
		BlacklistVersion           int
		FirmwareVersion            int
		CCHSVer                    int
		LocationID                 int
		InterimBlacklistVersion    int
		FunctionalBlacklistVersion int
	}

	Card struct {
		CardID                 string
		RemainingValue         int
		AddValueDetail         string
		LastAddValueDevice     string
		Class                  string
		Language               string
		AvailableAutopayAmount string

		UniqueManufactureID string

		Logs []CardLog
	}

	CardLog struct {
		ServiceProviderID string
		TransactionAmount string
		TransactionTime   time.Time
		MachineID         string
		ServiceInfo       string
	}

	XFileResult struct {
		FileName    string
		WarningCode int
	}

	DeductResult struct {
		RemainingValue int
		AdditionalInfo Hex
	}

	GetLastAddValueInfoResult struct {
		Date     string
		TypeCode string
		Type     string
		DeviceID string
	}

	InitArgs struct {
		PortNumber   int
		BaudRate     int
		ControllerID int
	}

	WriteLocationArgs struct {
		LocationID int
	}

	PollArgs struct {
		Command int
		Timeout int
	}

	Hex []byte

	DeductArgs struct {
		Value            int
		ServiceInfo      Hex
		DeferReleaseFlag int
	}

	TxnAmtArgs struct {
		Value          int
		RemainingValue int
		LED            int
		Sound          int
	}
)

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToUpper(hex.EncodeToString(h)))
}

func (h *Hex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var unquoted string
	var err error
	unquoted, err = strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	*h, err = hex.DecodeString(unquoted)
	return err
}

func (octopus *Octopus) Init(args *InitArgs, reply *CardReaderInfo) error {
	// 0: use default settings in init file
	portNumber := C.uchar(args.PortNumber)

	baudRate := C.int(args.BaudRate)

	// 0: use default settings in init file
	controllerID := C.ulong(args.ControllerID)

	initRet := int(C.InitComm(portNumber, baudRate, controllerID))
	if initRet != 0 {
		log.Println("failed to init octopus", initRet)
		return errorForCode(initRet)
	}
	log.Println("successfully inited octopus")

	return octopus.Inspect(new(int), reply)
}

func (octopus *Octopus) UpdateLocationID(args *WriteLocationArgs, reply *bool) error {
	locationID := C.uint(args.LocationID)
	locRet := int(C.WriteID(locationID))
	if locRet != 0 {
		log.Println("failed to update location", locRet)
		return errorForCode(locRet)
	}
	log.Println("successfully updated location", locationID)
	*reply = true
	return nil
}

func (octopus *Octopus) GetLastAddValueInfo(_ *int, reply *GetLastAddValueInfoResult) error {
	data := C.malloc(C.sizeof_uchar * 512)
	defer C.free(unsafe.Pointer(data))
	ret := int(C.GetExtraInfo(C.uint(0), C.uint(1), (*C.uchar)(data)))
	if ret == 0 {
		parts := strings.SplitN(C.GoString((*C.char)(unsafe.Pointer(data))), ",", 3)
		var typ string
		switch parts[1] {
		case "1":
			typ = "Cash"
		case "2":
			typ = "Online"
		case "3":
			typ = "Refund"
		case "4":
			typ = "AAVS"
		default:
			typ = "Others"
		}
		*reply = GetLastAddValueInfoResult{
			Date:     parts[0],
			TypeCode: parts[1],
			Type:     typ,
			DeviceID: parts[2],
		}
		return nil
	}
	log.Println("failed to get last add value info", ret)
	return errorForCode(ret)
}

func (octopus *Octopus) Inspect(_ *int, reply *CardReaderInfo) error {
	data := C.malloc(C.sizeof_uchar * 2046)
	defer C.free(unsafe.Pointer(data))
	tvRet := int(C.TimeVer((*C.uchar)(data)))
	if tvRet == 0 {
		info := (*C.stDevVer)(unsafe.Pointer(data))
		*reply = CardReaderInfo{
			DeviceID:                   int(info.DevID),
			OperatorID:                 int(info.OperID),
			DeviceTime:                 time.Unix(TransactionTimeSince+int64(int(info.DevTime)), 0),
			CompanyID:                  int(info.CompID),
			KeyVersion:                 int(info.KeyVer),
			EODVersion:                 int(info.EODVer),
			BlacklistVersion:           int(info.BLVer),
			FirmwareVersion:            int(info.FIRMVer),
			CCHSVer:                    int(info.CCHSVer),
			LocationID:                 int(info.CSSer),
			InterimBlacklistVersion:    int(info.IntBLVer),
			FunctionalBlacklistVersion: int(info.FuncBLVer),
		}
		return nil
	}
	log.Println("failed to inspect card reader", tvRet)
	return errorForCode(tvRet)
}

func (octopus *Octopus) Poll(args *PollArgs, reply *Card) error {
	data := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(data))

	// 0: remaining value and card ID
	// 1: remaining value, card ID and IDm (unique manufacture ID)
	// 2: remaining value, card ID, IDm and card history
	command := C.uchar(args.Command)

	// timeout in 0.1 second; max: 30
	timeout := C.uchar(args.Timeout)

	pollRet := int(C.Poll(command, timeout, (*C.char)(data)))
	if pollRet > 100000 {
		log.Println("failed to query card", pollRet)
		return errorForCode(pollRet)
	}
	b := C.GoBytes(data, 1024)
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		log.Println("failed to get data from poll")
		return errors.New("invalid data")
	}

	parts := strings.Split(string(b[:n]), ",")
	parts1 := strings.Split(parts[1], "-")
	card := Card{
		CardID:                 parts[0],
		RemainingValue:         pollRet,
		AddValueDetail:         parts1[0],
		LastAddValueDevice:     parts1[1],
		Class:                  parts1[2],
		Language:               parts1[3],
		AvailableAutopayAmount: parts1[4],
	}

	partsLen := len(parts)

	if partsLen > 2 {
		card.UniqueManufactureID = parts[2]
	}

	if partsLen > 3 {
		i := 3
		for i < partsLen {
			card.Logs = append(card.Logs, CardLog{
				ServiceProviderID: parts[i],
				TransactionAmount: parts[i+1],
				TransactionTime:   parseTime(parts[i+2]),
				MachineID:         parts[i+3],
				ServiceInfo:       parts[i+4],
			})
			i += 5
		}
	}

	*reply = card
	log.Println("successfully queried card with ID", parts[0], "with remaining value", pollRet)
	return nil
}

func (octopus *Octopus) Deduct(args *DeductArgs, reply *DeductResult) error {
	if len(args.ServiceInfo) < 5 {
		log.Println("bad deduct service info")
		return errors.New("service info must be 5 bytes")
	}

	amount := C.int(args.Value)

	data := C.malloc(C.sizeof_uchar * 128)
	defer C.free(unsafe.Pointer(data))

	ud := make([]byte, 2)
	rand.Read(ud)
	ai := append(args.ServiceInfo[0:5], ud...)
	cai := C.CString(string(ai))
	defer C.free(unsafe.Pointer(cai))
	C.memcpy(unsafe.Pointer(data), unsafe.Pointer(cai), C.size_t(len(ai)))

	// 0: Release card after deduct
	// 1: Defer card release
	deferReleaseFlag := C.int(args.DeferReleaseFlag)

	deductRet := int(C.Deduct(amount, (*C.uchar)(data), deferReleaseFlag))
	if deductRet > 100000 {
		log.Println("failed to deduct", deductRet)
		return errorForCode(deductRet)
	}
	*reply = DeductResult{
		RemainingValue: deductRet,
	}
	b := C.GoBytes(data, 128)
	n := bytes.Index(b, ud)
	if n == -1 {
		log.Println("warning: failed to get result from additional info")
		return nil
	}
	(*reply).AdditionalInfo = b[:n+len(ud)]
	log.Println("successfully deducted value, remaining", deductRet)
	return nil
}

func (octopus *Octopus) TxnAmt(args *TxnAmtArgs, reply *bool) error {
	ret := int(C.TxnAmt(C.int(args.Value), C.int(args.RemainingValue), C.uchar(args.LED), C.uchar(args.Sound)))
	if ret >= 100000 {
		log.Println("failed to execute TxtAmt", ret)
		return errorForCode(ret)
	}
	*reply = true
	return nil
}

func (octopus *Octopus) GenerateExchangeFile(_ *int, reply *XFileResult) error {
	data := C.malloc(C.sizeof_char * 128)
	defer C.free(unsafe.Pointer(data))

	ret := int(C.XFile((*C.char)(data)))
	if ret >= 100000 {
		log.Println("failed to generate exchange file", ret)
		return errorForCode(ret)
	}

	b := C.GoBytes(data, 128)
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		log.Println("failed to get data from xfile")
		return errors.New("invalid data")
	}
	filename := string(b[:n])
	*reply = XFileResult{
		FileName:    filename,
		WarningCode: ret,
	}
	log.Println("successfully generated exchange file", filename)
	return nil
}

func parseTime(input string) time.Time {
	secs, _ := strconv.ParseInt(input, 10, 0)
	return time.Unix(TransactionTimeSince+secs, 0)
}

func errorForCode(code int) error {
	return errors.New(strconv.Itoa(code))
}

func main() {
	address := flag.String("address", "127.0.0.1:12345", "address to listen to")
	flag.Parse()

	err := rpc.Register(new(Octopus))
	if err != nil {
		log.Fatal(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening", tcpAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println("opened connection from", conn.RemoteAddr().String())
		jsonrpc.ServeConn(conn)
	}
}
