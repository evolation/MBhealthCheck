//go:build !js || !wasm

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

var (
	DEVELOPMENT = "true"
	rs485_port  = "COM9"
	version     = "2.0.1"
)

type Config struct {
	SlaveID byte `json:"slave_id"`
}

type Report struct {
	Active       []Config `json:"active"`
	TimeoutError []Config `json:"timeoutError"`
	CrcError     []Config `json:"crcError"`
}

func LoadConfig(file string) ([]Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var configs []Config
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func CreateRTUClient(config Config, p_port string) (modbus.Client, *modbus.RTUClientHandler, error) {
	handler := modbus.NewRTUClientHandler(p_port)
	handler.BaudRate = 38400
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 2
	handler.SlaveId = config.SlaveID
	handler.Timeout = 2 * time.Second

	err := handler.Connect()
	if err != nil {
		return nil, nil, err
	}

	rtuClient := modbus.NewClient(handler)
	return rtuClient, handler, nil
}

func ReadRegister(rtuClient modbus.Client, config Config) error {
	_, err := rtuClient.ReadHoldingRegisters(0x00, 1)
	return err
}

func runModbusHealthcheckTarget(p_filename string, p_port string) {

	configs, err := LoadConfig(p_filename)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}
	report := Report{}
	for _, config := range configs {

		rtuClient, handler, err := CreateRTUClient(config, p_port)
		if err != nil {
			log.Printf("Failed to create client for config: %v", err)
			continue
		}

		err = ReadRegister(rtuClient, config)

		if err != nil {
			mbErr, ok := err.(*modbus.ModbusError)
			if ok {
				switch mbErr.ExceptionCode {
				case modbus.ExceptionCodeIllegalFunction:
					log.Printf("Illegal function")
				case modbus.ExceptionCodeIllegalDataAddress:
					log.Printf("Illegal data address")
				case modbus.ExceptionCodeIllegalDataValue:
					log.Printf("Illegal data value")
				// Add more cases as needed
				default:
					log.Printf("Unknown Modbus error: %v", mbErr)
				}
				report.CrcError = append(report.CrcError, config)

			} else {
				// log.Printf("Failed to read register: %v", err)
				report.TimeoutError = append(report.TimeoutError, config)
			}
			fmt.Println("Slave ID: ", config.SlaveID, " ERROR")
		} else {
			report.Active = append(report.Active, config)
			fmt.Println("Slave ID: ", config.SlaveID, " OK")

		}
		handler.Close()
	}
	data, _ := json.MarshalIndent(report, "", "  ")

	ioutil.WriteFile("./report.json", data, 0644)
}
