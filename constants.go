package goipmi

type EntityId uint8

const (
	ENTITY_ID_UNSPECIFIED                         = 0
	ENTITY_ID_OTHER                               = 1
	ENTITY_ID_UNKNOWN                             = 2
	ENTITY_ID_PROCESSOR                           = 3
	ENTITY_ID_DISK                                = 4
	ENTITY_ID_PERIPHERALBAY                       = 5
	ENTITY_ID_SYSTEMMANAGEMENTMODULE              = 6
	ENTITY_ID_SYSTEMBOARD                         = 7
	ENTITY_ID_MEMORYMODULE                        = 8
	ENTITY_ID_PROCESORMODULE                      = 9
	ENTITY_ID_POWERSUPPLY                         = 10
	ENTITY_ID_ADDINCARD                           = 11
	ENTITY_ID_FRONTPANELBOARD                     = 12
	ENTITY_ID_BACKPANELBOARD                      = 13
	ENTITY_ID_POWERSYSTEMBOARD                    = 14
	ENTITY_ID_DRIVEBACKPLANE                      = 15
	ENTITY_ID_SYSTEMINTERNALEXPANSIONBOARD        = 16
	ENTITY_ID_OTHERSYSTEMBOARD                    = 17
	ENTITY_ID_PROCESSORBOARD                      = 18
	ENTITY_ID_POWERUNIT                           = 19
	ENTITY_ID_POWERMODULE                         = 20
	ENTITY_ID_POWERMANAGEMENT                     = 21
	ENTITY_ID_CHASSISBACKPANELBOARD               = 22
	ENTITY_ID_SYSTEMCHASSIS                       = 23
	ENTITY_ID_SUBCHASSIS                          = 24
	ENTITY_ID_OTHERCHASSIS                        = 25
	ENTITY_ID_DISKDRIVEBAY                        = 26
	ENTITY_ID_PERIPHERALBAY2                      = 27
	ENTITY_ID_DEVICEBAY                           = 28
	ENTITY_ID_FAN                                 = 29
	ENTITY_ID_COOLINGUNIT                         = 30
	ENTITY_ID_CABLEINTERCONNECT                   = 31
	ENTITY_ID_MEMORYDEVICE                        = 32
	ENTITY_ID_SYSTEMMANAGEMENTSOFTWARE            = 33
	ENTITY_ID_SYSTEMFIRMWARE                      = 34
	ENTITY_ID_OPERATINGSYSTEM                     = 35
	ENTITY_ID_SYSTEMBUS                           = 36
	ENTITY_ID_GROUP                               = 37
	ENTITY_ID_REMOTEMANAGEMENTCOMMUNICATIONDEVICE = 38
	ENTITY_ID_EXTERNALENVIRONMENT                 = 39
	ENTITY_ID_BATTERY                             = 40
	ENTITY_ID_PROCESSINGBLADE                     = 41
	ENTITY_ID_CONNECTIVITYSWITCH                  = 42
	ENTITY_ID_PROCESSORMEMORYMODULE               = 43
	ENTITY_ID_IOMODULE                            = 44
	ENTITY_ID_PROCESSORIOMODULE                   = 45
	ENTITY_ID_MANAGEMENTCONTROLLERFIRMWARE        = 46
	ENTITY_ID_IPMICHANNEL                         = 47
	ENTITY_ID_PCIBUS                              = 48
	ENTITY_ID_PCIEXPRESSBUS                       = 49
	ENTITY_ID_SCSIBUS                             = 50
	ENTITY_ID_SATABUS                             = 51
	ENTITY_ID_FRONTSIDEBUS                        = 52
	ENTITY_ID_REALTIMECLOCK                       = 53
	ENTITY_ID_AIRINLET                            = 55
	ENTITY_ID_AIRINLET2                           = 64
	ENTITY_ID_PROCESSOR2                          = 65
	ENTITY_ID_BASEBOARD                           = 66
)

/**
 * Specifies available units for sensors' measurements.
 */
type SensorUnit uint8

const (
	SENSOR_UNIT_OTHER              = 0
	SENSOR_UNIT_UNSPECIFIED        = 0
	SENSOR_UNIT_DEGREESC           = 1
	SENSOR_UNIT_DEGREESF           = 2
	SENSOR_UNIT_DEGREESK           = 3
	SENSOR_UNIT_VOLTS              = 4
	SENSOR_UNIT_AMPS               = 5
	SENSOR_UNIT_WATTS              = 6
	SENSOR_UNIT_JOULES             = 7
	SENSOR_UNIT_COULOMBS           = 8
	SENSOR_UNIT_VA                 = 9
	SENSOR_UNIT_NITS               = 10
	SENSOR_UNIT_LUMEN              = 11
	SENSOR_UNIT_LUX                = 12
	SENSOR_UNIT_CANDELA            = 13
	SENSOR_UNIT_KPA                = 14
	SENSOR_UNIT_PSI                = 15
	SENSOR_UNIT_NEWTON             = 16
	SENSOR_UNIT_CFM                = 17
	SENSOR_UNIT_RPM                = 18
	SENSOR_UNIT_HZ                 = 19
	SENSOR_UNIT_MICROSECOND        = 20
	SENSOR_UNIT_MILLISECOND        = 21
	SENSOR_UNIT_SECOND             = 22
	SENSOR_UNIT_MINUTE             = 23
	SENSOR_UNIT_HOUR               = 24
	SENSOR_UNIT_DAY                = 25
	SENSOR_UNIT_WEEK               = 26
	SENSOR_UNIT_MIL                = 27
	SENSOR_UNIT_INCHES             = 28
	SENSOR_UNIT_FEET               = 29
	SENSOR_UNIT_CUIN               = 30
	SENSOR_UNIT_CUFEET             = 31
	SENSOR_UNIT_MM                 = 32
	SENSOR_UNIT_CM                 = 33
	SENSOR_UNIT_M                  = 34
	SENSOR_UNIT_CUCM               = 35
	SENSOR_UNIT_CUM                = 36
	SENSOR_UNIT_LITERS             = 37
	SENSOR_UNIT_FLUIDOUNCE         = 38
	SENSOR_UNIT_RADIANS            = 39
	SENSOR_UNIT_STERADIANS         = 40
	SENSOR_UNIT_REVOLUTIONS        = 41
	SENSOR_UNIT_CYCLES             = 42
	SENSOR_UNIT_GRAVITIES          = 43
	SENSOR_UNIT_OUNCE              = 44
	SENSOR_UNIT_POUND              = 45
	SENSOR_UNIT_FTLB               = 46
	SENSOR_UNIT_OZIN               = 47
	SENSOR_UNIT_GAUSS              = 48
	SENSOR_UNIT_GILBERTS           = 49
	SENSOR_UNIT_HENRY              = 50
	SENSOR_UNIT_MILLIHENRY         = 51
	SENSOR_UNIT_FARAD              = 52
	SENSOR_UNIT_MICROFARAD         = 53
	SENSOR_UNIT_OHMS               = 54
	SENSOR_UNIT_SIEMENS            = 55
	SENSOR_UNIT_MOLE               = 56
	SENSOR_UNIT_BECQUEREL          = 57
	SENSOR_UNIT_PARTSPERMILION     = 58
	SENSOR_UNIT_DECIBELS           = 60
	SENSOR_UNIT_DBA                = 61
	SENSOR_UNIT_DBC                = 62
	SENSOR_UNIT_GRAY               = 63
	SENSOR_UNIT_SIEVERT            = 64
	SENSOR_UNIT_COLORTEMPDEGK      = 65
	SENSOR_UNIT_BIT                = 66
	SENSOR_UNIT_KILOBIT            = 67
	SENSOR_UNIT_MEGABIT            = 68
	SENSOR_UNIT_GIGABIT            = 69
	SENSOR_UNIT_BYTE               = 70
	SENSOR_UNIT_KILOBYTE           = 71
	SENSOR_UNIT_MEGABYTE           = 72
	SENSOR_UNIT_GIGABYTE           = 73
	SENSOR_UNIT_WORD               = 74
	SENSOR_UNIT_DWORD              = 75
	SENSOR_UNIT_QWORD              = 76
	SENSOR_UNIT_LINE               = 77
	SENSOR_UNIT_HIT                = 78
	SENSOR_UNIT_MISS               = 79
	SENSOR_UNIT_RETRY              = 80
	SENSOR_UNIT_RESET              = 81
	SENSOR_UNIT_OVERRUNOVERFLOW    = 82
	SENSOR_UNIT_UNDERRUN           = 83
	SENSOR_UNIT_COLLISION          = 84
	SENSOR_UNIT_PACKETS            = 85
	SENSOR_UNIT_MESSAGES           = 86
	SENSOR_UNIT_CHARACTERS         = 87
	SENSOR_UNIT_ERROR              = 88
	SENSOR_UNIT_CORRECTABLEERROR   = 89
	SENSOR_UNIT_UNCORRECTABLEERROR = 90
	SENSOR_UNIT_FATALERROR         = 91
	SENSOR_UNIT_GRAMS              = 92
)

const (
	SENSOR_TEMPERATURE                      = 1
	SENSOR_VOLTAGE                          = 2
	SENSOR_CURRENT                          = 3
	SENSOR_FAN                              = 4
	SENSOR_PHYSICALSECURITY                 = 5
	SENSOR_PLATFORMSECURITYVIOLATIONATTEMPT = 6
	SENSOR_PROCESSOR                        = 7
	SENSOR_POWERSUPPLY                      = 8
	SENSOR_POWERUNIT                        = 9
	SENSOR_COOLINGDEVICE                    = 10
	SENSOR_OTHERUNITSBASEDSENSOR            = 11
	SENSOR_MEMORY                           = 12
	SENSOR_DRIVEBAY                         = 13
	SENSOR_POSTMEMORYRESIZE                 = 14
	SENSOR_SYSTEMFIRMWAREPROGESS            = 15
	SENSOR_EVENTLOGGINGDISABLED             = 16
	SENSOR_WATCHDOG1                        = 17
	SENSOR_SYSTEMEVENT                      = 18
	SENSOR_CRITICALINTERRUPT                = 19
	SENSOR_BUTTONSWITCH                     = 20
	SENSOR_MODULEBOARD                      = 21
	SENSOR_MICROCONTROLLERCOPROCESSOR       = 22
	SENSOR_ADDINCARD                        = 23
	SENSOR_CHASSIS                          = 24
	SENSOR_CHIPSET                          = 25
	SENSOR_OTHERFRU                         = 26
	SENSOR_CABLEINTERCONNECT                = 27
	SENSOR_TERMINATOR                       = 28
	SENSOR_SYSTEMBOOT                       = 29
	SENSOR_BOOTERROR                        = 30
	SENSOR_OSBOOT                           = 31
	SENSOR_OSSTOP                           = 32
	SENSOR_SLOTCONNECTOR                    = 33
	SENSOR_SYSTEMACPIPOWERSTATE             = 34
	SENSOR_WATCHDOG2                        = 35
	SENSOR_PLATFORMALERT                    = 36
	SENSOR_ENTITYPRESENCE                   = 37
	SENSOR_MONITORASICIC                    = 38
	SENSOR_LAN                              = 39
	SENSOR_MANAGEMENTSUBSYSTEMHEALTH        = 40
	SENSOR_BATTERY                          = 41
	SENSOR_SESSIONAUDIT                     = 42
	SENSOR_VERSIONCHANGE                    = 43
	SENSOR_FRUSTATE                         = 44
	SENSOR_OEM                              = 192
	SENSOR_OEMRESERVED                      = 118
)
