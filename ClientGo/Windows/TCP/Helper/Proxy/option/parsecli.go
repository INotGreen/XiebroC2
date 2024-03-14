package option

import (
	"encoding/hex"
	"errors"
	"main/Helper/Proxy/crypto"

	"strconv"
)

var (
	errUnrecognizedMode    = errors.New("Unrecognized mode. Must choose a working mode in [fwd/proxy]")
	errHexDecodeError      = errors.New("KEY must be a hexadecimal string")
	PrintUsage             = errors.New("")
	errUnrecognizedSubMode = errors.New("Malformed args. Incorrect number of `-l/-r` params")
	errNoSecretKey         = errors.New("Encryption enabled, must specify a KEY by `-k` param")
	errNotANumber          = errors.New("Timeout param must be a number")
	errUDPMode             = errors.New("UDP mode only support fwd mode")
)

const (
	SUBMODE_L2L = iota
	SUBMODE_R2R
	SUBMODE_L2R

	SUBMODE_LP
	SUBMODE_RP
	SUBMODE_RPL2L
)

// Dont need flag-lib
func ParseCli(args []string) (
	mode string,
	submode int,
	local []string,
	remote []string,
	lenc []bool,
	renc []bool,
	err error) {

	if len(args) == 0 {
		err = PrintUsage
		return
	}

	mode = args[0]

	switch mode {
	case "fwd", "proxy":
	case "-h", "--help":
		err = PrintUsage
		return
	default:
		err = errUnrecognizedMode
		return
	}

	args = args[1:]
	ptr := 0

	for {
		if ptr == len(args) {
			break
		}

		switch args[ptr] {
		case "-l", "--local":
			l := args[ptr+1]
			if l[0] == '*' {
				lenc = append(lenc, true)
				l = l[1:]
			} else {
				lenc = append(lenc, false)
			}

			if _, err := strconv.Atoi(l); err == nil {
				local = append(local, "0.0.0.0:"+l)
			} else {
				if l[0] == ':' {
					local = append(local, "0.0.0.0"+l)
				} else {
					local = append(local, l)
				}
			}
			ptr++

		case "-r", "--remote":
			r := args[ptr+1]
			if r[0] == '*' {
				renc = append(renc, true)
				r = r[1:]
			} else {
				renc = append(renc, false)
			}

			remote = append(remote, r)
			ptr++

		case "-u", "--udp":
			PROTOCOL = "UDP"

		case "-k", "--key":
			var key []byte
			key, err = hex.DecodeString(args[ptr+1])
			if err != nil {
				err = errHexDecodeError
				return
			}
			crypto.ExpandKey(key)
			ptr++

		case "-t", "--timeout":
			TIMEOUT, err = strconv.Atoi(args[ptr+1])
			if err != nil {
				err = errNotANumber
				return
			}
			ptr++
		case "-v", "--verbose":
			VERBOSE = true
		case "-h", "--help":
			err = PrintUsage
			return
		}

		ptr++
	}

	if mode == "fwd" {
		switch {
		case len(local) == 0 && len(remote) == 2:
			submode = SUBMODE_R2R
		case len(local) == 1 && len(remote) == 1:
			submode = SUBMODE_L2R
		case len(local) == 2 && len(remote) == 0:
			submode = SUBMODE_L2L
		default:
			err = errUnrecognizedSubMode
			return
		}
	} else {
		switch {
		case len(local) == 0 && len(remote) == 1:
			submode = SUBMODE_RP
		case len(local) == 1 && len(remote) == 0:
			submode = SUBMODE_LP
		case len(local) == 2 && len(remote) == 0:
			submode = SUBMODE_RPL2L
		default:
			err = errUnrecognizedSubMode
			return
		}
	}

	if len(lenc) != len(local) || len(renc) != len(remote) {
		err = errUnrecognizedSubMode
		return
	}

	if crypto.SECRET_KEY == nil {
		for i, _ := range lenc {
			if lenc[i] {
				err = errNoSecretKey
				return
			}
		}

		for i, _ := range renc {
			if renc[i] {
				err = errNoSecretKey
				return
			}
		}
	}

	if PROTOCOL == "UDP" && mode == "proxy" {
		err = errUDPMode
		return
	}

	shouldFwdWithoutDec(lenc, renc)

	return
}

func shouldFwdWithoutDec(lenc []bool, renc []bool) {
	if len(lenc)+len(renc) != 2 {
		return
	}

	var result uint8
	for i, _ := range lenc {
		if lenc[i] {
			result++
		}
	}

	for i, _ := range renc {
		if renc[i] {
			result++
		}
	}

	if result == 2 {
		FORWARD_WITHOUT_DEC = true
	}
}
