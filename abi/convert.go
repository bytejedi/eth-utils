package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

func hexToAddress(v string) (ethcmn.Address, error) {
	if !ethcmn.IsHexAddress(v) {
		return ethcmn.Address{}, fmt.Errorf("invalid address %s", v)
	}
	return ethcmn.HexToAddress(v), nil
}

func stringToInt(ty ethabi.Type, src string) (interface{}, error) {
	var dst interface{}
	if ty.T == ethabi.IntTy && ty.Size <= 64 {
		tmp, err := strconv.ParseInt(src, 10, ty.Size)
		if err != nil {
			return nil, err
		}
		switch ty.Size {
		case 8:
			dst = int8(tmp)
		case 16:
			dst = int16(tmp)
		case 32:
			dst = int32(tmp)
		case 64:
			dst = tmp
		}
	} else if ty.T == ethabi.UintTy && ty.Size <= 64 {
		tmp, err := strconv.ParseUint(src, 10, ty.Size)
		if err != nil {
			return nil, err
		}
		switch ty.Size {
		case 8:
			dst = uint8(tmp)
		case 16:
			dst = uint16(tmp)
		case 32:
			dst = uint32(tmp)
		case 64:
			dst = tmp
		}
	} else {
		var ok bool
		dst, ok = new(big.Int).SetString(src, 10)
		if !ok {
			return nil, errors.New("big.Int SetString failed")
		}
	}
	return dst, nil
}

func JsonToGoType(ty ethabi.Type, j string) (interface{}, error) {
	var src interface{}
	err := json.Unmarshal([]byte(j), &src)
	if err != nil {
		return nil, err
	}

	if ty.T == ethabi.SliceTy || ty.T == ethabi.ArrayTy {

		if ty.Elem.T == ethabi.AddressTy {
			var tmp = src.([]string)
			var dst []ethcmn.Address
			for _, v := range tmp {
				addr, err := hexToAddress(v)
				if err != nil {
					return nil, err
				}
				dst = append(dst, addr)
			}
			return dst, nil
		}

		if (ty.Elem.T == ethabi.IntTy || ty.Elem.T == ethabi.UintTy) && reflect.TypeOf(src).Elem().Kind() == reflect.Interface {
			var tmp = src.([]string)
			var dst []interface{}
			for _, v := range tmp {
				i, err := stringToInt(ty, v)
				if err != nil {
					return nil, err
				}
				dst = append(dst, i)
			}
			return dst, nil
		}
	}

	if ty.T == ethabi.AddressTy {
		return hexToAddress(src.(string))
	}

	if (ty.T == ethabi.IntTy || ty.T == ethabi.UintTy) && reflect.TypeOf(src).Kind() == reflect.String {
		return stringToInt(ty, src.(string))
	}

	return nil, fmt.Errorf("cannot convert %T to Golang type", src)
}
