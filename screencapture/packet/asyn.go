package packet

import (
	"encoding/binary"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/dict"
)

//Async Packet types
const (
	AsynPacketMagic uint32 = 0x6173796E
	FEED            uint32 = 0x66656564 //These contain CMSampleBufs which contain raw h264 Nalus
	TJMP            uint32 = 0x746A6D70
	SRAT            uint32 = 0x73726174 //CMTimebaseSetRateAndAnchorTime https://developer.apple.com/documentation/coremedia/cmtimebase?language=objc
	SPRP            uint32 = 0x73707270 // Set Property
	TBAS            uint32 = 0x74626173 //TimeBase https://developer.apple.com/library/archive/qa/qa1643/_index.html
	RELS            uint32 = 0x72656C73
	HPD1            uint32 = 0x68706431 //hpd1 - 1dph | For specifying/requesting the video format
	HPA1            uint32 = 0x68706131 //hpa1 - 1aph | For specifying/requesting the audio format
	NEED            uint32 = 0x6E656564 //need - deen
	EAT             uint32 = 0x65617421 //contains audio sbufs
	HPD0            uint32 = 0x68706430
	HPA0            uint32 = 0x68706130
)

//NewAsynHpd1Packet creates a []byte containing a valid ASYN packet with the Hpd1 dictionary
func NewAsynHpd1Packet(stringKeyDict dict.StringKeyDict) []byte {
	return newAsynDictPacket(stringKeyDict, HPD1, EmptyCFType)
}

//NewAsynHpa1Packet creates a []byte containing a valid ASYN packet with the Hpa1 dictionary
func NewAsynHpa1Packet(stringKeyDict dict.StringKeyDict, clockRef CFTypeID) []byte {
	return newAsynDictPacket(stringKeyDict, HPA1, clockRef)
}

func newAsynDictPacket(stringKeyDict dict.StringKeyDict, subtypeMarker uint32, asynTypeHeader uint64) []byte {
	serialize := dict.SerializeStringKeyDict(stringKeyDict)
	length := len(serialize) + 20
	header := make([]byte, 20)
	binary.LittleEndian.PutUint32(header, uint32(length))
	binary.LittleEndian.PutUint32(header[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(header[8:], asynTypeHeader)
	binary.LittleEndian.PutUint32(header[16:], subtypeMarker)
	return append(header, serialize...)
}

//AsynNeedPacketBytes can be used to create the NEED message as soon as the clockRef from SYNC CVRP has been received.
func AsynNeedPacketBytes(clockRef CFTypeID) []byte {
	needPacketLength := 20
	packet := make([]byte, needPacketLength)
	binary.LittleEndian.PutUint32(packet, uint32(needPacketLength))
	binary.LittleEndian.PutUint32(packet[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(packet[8:], clockRef)
	binary.LittleEndian.PutUint32(packet[16:], NEED) //need - deen
	return packet
}

//CreateHpd1DeviceInfoDict creates a dict.StringKeyDict that needs to be sent to the device before receiving a feed
func CreateHpd1DeviceInfoDict() dict.StringKeyDict {
	resultDict := dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 3)}
	displaySizeDict := dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 2)}
	resultDict.Entries[0] = dict.StringKeyEntry{
		Key:   "Valeria",
		Value: true,
	}
	resultDict.Entries[1] = dict.StringKeyEntry{
		Key:   "HEVCDecoderSupports444",
		Value: true,
	}

	displaySizeDict.Entries[0] = dict.StringKeyEntry{
		Key:   "Width",
		Value: common.NewNSNumberFromUFloat64(1920),
	}
	displaySizeDict.Entries[1] = dict.StringKeyEntry{
		Key:   "Height",
		Value: common.NewNSNumberFromUFloat64(1200),
	}

	resultDict.Entries[2] = dict.StringKeyEntry{
		Key:   "DisplaySize",
		Value: displaySizeDict,
	}

	return resultDict
}

//CreateHpa1DeviceInfoDict creates a dict.StringKeyDict that needs to be sent to the device before receiving a feed
func CreateHpa1DeviceInfoDict() dict.StringKeyDict {
	resultDict := dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 6)}
	resultDict.Entries[0] = dict.StringKeyEntry{
		Key:   "BufferAheadInterval",
		Value: common.NewNSNumberFromUFloat64(0.07300000000000001),
	}

	resultDict.Entries[1] = dict.StringKeyEntry{
		Key:   "deviceUID",
		Value: "Valeria",
	}

	resultDict.Entries[2] = dict.StringKeyEntry{
		Key:   "ScreenLatency",
		Value: common.NewNSNumberFromUFloat64(0.04),
	}

	resultDict.Entries[3] = dict.StringKeyEntry{
		Key:   "formats",
		Value: coremedia.CreateLpcmInfo(),
	}

	resultDict.Entries[4] = dict.StringKeyEntry{
		Key:   "EDIDAC3Support",
		Value: common.NewNSNumberFromUInt32(0),
	}

	resultDict.Entries[5] = dict.StringKeyEntry{
		Key:   "deviceName",
		Value: "Valeria",
	}
	return resultDict
}