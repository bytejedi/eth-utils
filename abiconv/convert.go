package abiconv

import (
	"errors"
	"fmt"
	"math/big"
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

func StringToGoType(ty ethabi.Type, src string) (interface{}, error) {
	switch ty.T {
	case ethabi.AddressTy:
		return hexToAddress(src)
	case ethabi.IntTy, ethabi.UintTy:
		return stringToInt(ty, src)
	default:
		return nil, fmt.Errorf("cannot convert %T to Golang type", src)
	}
}

func StringSliceToGoType(ty ethabi.Type, src []string) (interface{}, error) {
	switch ty.T {
	case ethabi.SliceTy, ethabi.ArrayTy:
		var dst []interface{}
		for _, str := range src {
			data, err := StringToGoType(*ty.Elem, str)
			if err != nil {
				return nil, err
			}
			dst = append(dst, data)
		}
		return dst, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to Golang type", src)
	}
}
