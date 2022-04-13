package sc

import (
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
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

func UnpackEncodedEvent(ver uint64, packed []byte) interface{} {
	switch ver {
	case 2:
		encodedEvent := map[string]interface{}{}
		abi, err := abi.JSON(strings.NewReader(RequestValueTransferEncodeABIs[2]))
		if err != nil {
			logger.Error("Failed to setup ABI of RequestValueTransferEncodedData", "err", err)
			return nil
		}
		if err := abi.UnpackIntoMap(encodedEvent, "packedURI", packed); err != nil {
			logger.Error("Failed to unpack the values", "err", err)
			return nil
		}
		return encodedEvent["uri"]
	default:
		logger.Error(ErrUnknownEvent.Error(), "ver", ver)
		return nil
	}
}
