package utils

import "log"

func ExtractMCCMNC(imsi string) (mcc string, mnc string, err error) {
	// Check if the IMSI has a valid length
	if len(imsi) < 6 {
		log.Println("Invalid IMSI")

		return "000", "00", nil
	}

	// Extract the MCC and MNC
	mcc = imsi[:3]
	if err != nil {
		return
	}

	// Determine the length of MNC (2 or 3 digits)
	mncLength := 3
	if imsi[3] == '0' {
		mncLength = 2
	}

	mnc = imsi[3 : 3+mncLength]

	return
}
