package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/runner-mei/goipmi"
	ipmi_proto "github.com/runner-mei/goipmi/protocol"
)

func main() {
	cli, err := ipmi_proto.NewClient(&ipmi_proto.ConnectionOption{
		Hostname:  "192.168.1.15",
		Port:      623,
		Username:  "Administrator",
		Password:  "123456abc",
		Interface: "lanplus",
	})

	if nil != err {
		fmt.Println(err)
		return
	}
	client := &goipmi.Client{cli}

	if err := client.Open(); nil != err {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := client.Close(); nil != err {
			fmt.Println(err)
		}
	}()

	resp, err := client.GetDeviceID()
	if nil != err {
		log.Fatalln(err)
	}

	fmt.Println("DeviceID:", resp.DeviceID)
	fmt.Println("DeviceRevision:", resp.DeviceRevision)
	fmt.Println("FirmwareRevision1:", resp.FirmwareRevision1)
	fmt.Println("FirmwareRevision2:", resp.FirmwareRevision2)

	getACPIPowerStateResponse, err := client.GetACPIPowerState()
	if err != nil {
		log.Fatalln("get Sdr info,", err)
	}
	fmt.Println("================ ACPI State ================")
	fmt.Printf("SystemPower = %v\r\n", getACPIPowerStateResponse.GetACPISystemPowerState())
	fmt.Printf("DevicePower = %v\r\n", getACPIPowerStateResponse.GetACPIDevicePowerState())

	getChassisCapabilitiesResponse, err := client.GetChassisCapabilities()
	if err != nil {
		log.Fatalln("get Chassis Capabilities info,", err)
	}
	fmt.Println("================ Chassis Capabilities ================")
	fmt.Printf("PowerInterlock = %v\r\n", getChassisCapabilitiesResponse.PowerInterlock())
	fmt.Printf("DiagnosticInterrupt = %v\r\n", getChassisCapabilitiesResponse.DiagnosticInterrupt())
	fmt.Printf("FrontPanelLockout = %v\r\n", getChassisCapabilitiesResponse.FrontPanelLockout())
	fmt.Printf("HasInstrusionSensor = %v\r\n", getChassisCapabilitiesResponse.HasInstrusionSensor())
	fmt.Printf("ChassisFruInfoDeviceAddress = %v\r\n", getChassisCapabilitiesResponse.ChassisFruInfoDeviceAddress)
	fmt.Printf("ChassisSDRDeviceAddress = %v\r\n", getChassisCapabilitiesResponse.ChassisSDRDeviceAddress)
	fmt.Printf("ChassisSELDeviceAddress = %v\r\n", getChassisCapabilitiesResponse.ChassisSELDeviceAddress)
	fmt.Printf("ChassisSystemManagementDeviceAddress = %v\r\n", getChassisCapabilitiesResponse.ChassisSystemManagementDeviceAddress)
	fmt.Printf("ChassisBridgeDeviceAddress = %v\r\n", getChassisCapabilitiesResponse.ChassisBridgeDeviceAddress)

	getChassisStatusResponse, err := client.GetChassisStatus()
	if err != nil {
		log.Fatalln("get Chassis Status,", err)
	}
	fmt.Println("================ Chassis Status ================")
	fmt.Printf("PowerRestorePolicy = %v\r\n", getChassisStatusResponse.PowerRestorePolicyString())
	fmt.Printf("PowerControlFault = %v\r\n", getChassisStatusResponse.PowerControlFault())
	fmt.Printf("PowerFault = %v\r\n", getChassisStatusResponse.PowerFault())
	fmt.Printf("Interlock = %v\r\n", getChassisStatusResponse.Interlock())
	fmt.Printf("PowerOverload = %v\r\n", getChassisStatusResponse.PowerOverload())
	fmt.Printf("PowerOn = %v\r\n", getChassisStatusResponse.PowerOn())
	fmt.Printf("LastPowerEventPowerOn = %v\r\n", getChassisStatusResponse.LastPowerEventPowerOn())
	fmt.Printf("LastPowerEventPowerDownByFault = %v\r\n", getChassisStatusResponse.LastPowerEventPowerDownByFault())
	fmt.Printf("LastPowerEventPowerDownByInterlock = %v\r\n", getChassisStatusResponse.LastPowerEventPowerDownByInterlock())
	fmt.Printf("LastPowerEventPowerDownByOverload = %v\r\n", getChassisStatusResponse.LastPowerEventPowerDownByOverload())
	fmt.Printf("LastPowerEventACFailed = %v\r\n", getChassisStatusResponse.LastPowerEventACFailed())
	fmt.Printf("IdentityCommandSupported = %v\r\n", getChassisStatusResponse.IdentityCommandSupported())
	fmt.Printf("FanFault = %v\r\n", getChassisStatusResponse.FanFault())
	fmt.Printf("DriverFault = %v\r\n", getChassisStatusResponse.DriverFault())
	fmt.Printf("FrontPanelLockoutActived = %v\r\n", getChassisStatusResponse.FrontPanelLockoutActived())
	fmt.Printf("ChassisInstrusionActived = %v\r\n", getChassisStatusResponse.ChassisInstrusionActived())
	fmt.Printf("FrontPanelStandbyButtonDisableAllowed = %v\r\n", getChassisStatusResponse.FrontPanelStandbyButtonDisableAllowed())
	fmt.Printf("FrontPanelDiagnosticInterruptButtonDisableAllowed = %v\r\n", getChassisStatusResponse.FrontPanelDiagnosticInterruptButtonDisableAllowed())
	fmt.Printf("FrontPanelResetButtonDisableAllowed = %v\r\n", getChassisStatusResponse.FrontPanelResetButtonDisableAllowed())
	fmt.Printf("FrontPanelPowerOffButtonDisableAllowed = %v\r\n", getChassisStatusResponse.FrontPanelPowerOffButtonDisableAllowed())
	fmt.Printf("FrontPanelStandbyButtonDisabled = %v\r\n", getChassisStatusResponse.FrontPanelStandbyButtonDisabled())
	fmt.Printf("FrontPanelDiagnosticInterruptButtonDisabled = %v\r\n", getChassisStatusResponse.FrontPanelDiagnosticInterruptButtonDisabled())
	fmt.Printf("FrontPanelResetButtonDisabled = %v\r\n", getChassisStatusResponse.FrontPanelResetButtonDisabled())
	fmt.Printf("FrontPanelPowerOffButtonDisabled = %v\r\n", getChassisStatusResponse.FrontPanelPowerOffButtonDisabled())

	getSystemRestartCauseResponse, err := client.GetSystemRestartCause()
	if err != nil {
		log.Fatalln("get System restart cause,", err)
	}
	fmt.Println("================ System restart cause ================")
	fmt.Printf("ResetCauseString = %v\r\n", getSystemRestartCauseResponse.GetResetCauseString())

	getSDRInfoResponse, err := client.GetSDRRepositoryInfo()
	if err != nil {
		log.Fatalln("get Sdr repo info,", err)
	}
	fmt.Println("================ SDR Info ================")
	fmt.Printf("Version = %v\r\n", getSDRInfoResponse.GetVersion())
	fmt.Printf("RecordCount = %v\r\n", getSDRInfoResponse.RecordCount)
	fmt.Printf("RecentAddTimestamp = %v\r\n", time.Unix(int64(getSDRInfoResponse.RecentAddTimestamp), 0))
	fmt.Printf("RecentDelTimestamp = %v\r\n", time.Unix(int64(getSDRInfoResponse.RecentDelTimestamp), 0))
	fmt.Printf("OperationSupport = %v\r\n", getSDRInfoResponse.OperationSupport)

	reserveSDRResponse, err := client.GetReserveSDRRepository()
	if err != nil {
		log.Fatalln("get Reserve Sdr info,", err)
	}
	fmt.Println("================ Reserve SDR Info ================")
	fmt.Printf("Id = %v\r\n", reserveSDRResponse.Id)

	fmt.Println("================ List SDR ================")
	records, err := client.ListSDR(reserveSDRResponse.Id)
	if err != nil {
		log.Fatalln("get List Sdr,", err)
	}

	fullRecords, readings, err := client.ListFullSDRReading(records)
	if err != nil {
		log.Fatalln("get Sdr reading,", err)
	}
	for idx, full := range fullRecords {
		if readings[idx].Error != nil {
			fmt.Println(full.RecordId, full.EntityId, full.IdString, readings[idx].Error)
			continue
		}
		value, err := full.Calc(int32(readings[idx].Response.Reading), 8)
		fmt.Println(full.RecordId, full.EntityId, full.IdString, value, err)

		fmt.Println(readings[idx].Response.ToEventString(full.EventOrReadingTypeCode))

		response, err := client.GetSensorThresholds(full.SensorNumber)
		if nil != err {
			fmt.Println(err)
			continue
		}

		fmt.Println("HasUpperNonrecoverableThreshold:", response.HasUpperNonrecoverableThreshold())
		fmt.Println("HasUpperCriticalThreshold:", response.HasUpperCriticalThreshold())
		fmt.Println("HasUpperNonCriticalThreshold:", response.HasUpperNonCriticalThreshold())
		fmt.Println("HasLowerNonrecoverableThreshold:", response.HasLowerNonrecoverableThreshold())
		fmt.Println("HasLowerCriticalThreshold:", response.HasLowerCriticalThreshold())
		fmt.Println("HasLowerNonCriticalThreshold:", response.HasLowerNonCriticalThreshold())

		fmt.Println("UpperNonrecoverableThreshold:", response.UpperNonrecoverableThreshold)
		fmt.Println("UpperCriticalThreshold:", response.UpperCriticalThreshold)
		fmt.Println("UpperNonCriticalThreshold:", response.UpperNonCriticalThreshold)
		fmt.Println("LowerNonrecoverableThreshold:", response.LowerNonrecoverableThreshold)
		fmt.Println("LowerCriticalThreshold:", response.LowerCriticalThreshold)
		fmt.Println("LowerNonCriticalThreshold:", response.LowerNonCriticalThreshold)
	}

	for _, rec := range records {
		if _, ok := rec.(*goipmi.FullSensorRecord); !ok {

			fmt.Printf("%T\r\n", rec)
			bs, _ := json.Marshal(rec)
			fmt.Println(string(bs))
		}
	}
}
