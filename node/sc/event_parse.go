package sc

import (
	"bytes"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/pkg/errors"
)

var (
	ErrUnknownEvent = errors.New("Unknown event type")
)

var RequestValueTransferEncodeABIs = map[uint]string{
	2: `[{
			"anonymous":false,
			"inputs": [{
				"name": "uri",
				"type": "string"
			}],
			"name": "packedURI",
			"type": "event"
		}]`,
}

func UnpackEncodedData(ver uint8, packed []byte) map[string]interface{} {
	switch ver {
	case 2:
		encodedEvent := map[string]interface{}{}
		abi, err := abi.JSON(strings.NewReader(RequestValueTransferEncodeABIs[2]))
		if err != nil {
			logger.Error("Failed to ABI setup", "err", err)
			return nil
		}
		if err := abi.UnpackIntoMap(encodedEvent, "packedURI", packed); err != nil {
			logger.Error("Failed to unpack the values", "err", err)
			return nil
		}
		return encodedEvent
	default:
		logger.Error(ErrUnknownEvent.Error(), "encodingVer", ver)
		return nil
	}
}

func GetURI(ev IRequestValueTransferEvent) string {
	switch evType := ev.(type) {
	case RequestValueTransferEncodedEvent:
		decoded := UnpackEncodedData(evType.EncodingVer, evType.EncodedData)
		uri, ok := decoded["uri"].(string)
		if !ok {
			return ""
		}
		if len(uri) <= 64 {
			return ""
		}
		return string(bytes.Trim([]byte(uri[64:]), "\x00"))
	}
	return ""
}
