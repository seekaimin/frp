package commontools

import (
	"fmt"
	"strconv"
)

func ToUInt8(src string) uint8 {
	r, _ := strconv.ParseUint(src, 0, 8)
	return uint8(r)
}
func ToUInt16(src string) uint16 {
	r, _ := strconv.ParseUint(src, 0, 16)
	return uint16(r)
}
func ToUInt32(src string) uint32 {
	r, _ := strconv.ParseUint(src, 0, 32)
	return uint32(r)
}
func ToUInt64(src string) uint64 {
	r, _ := strconv.ParseUint(src, 0, 64)
	return uint64(r)
}
func HextoUInt8(src string) uint8 {
	r, _ := strconv.ParseUint(src, 16, 16)
	return uint8(r)
}
func HextoUInt16(src string) uint16 {
	r, _ := strconv.ParseUint(src, 16, 16)
	return uint16(r)
}
func HextoUInt32(src string) uint32 {
	r, _ := strconv.ParseUint(src, 16, 32)
	return uint32(r)
}
func HextoUInt64(src string) uint64 {
	r, _ := strconv.ParseUint(src, 16, 64)
	return r
}

func InterfaceToUInt8(item interface{}) (uint8, error) {
	str := fmt.Sprintf("%.0f", item)
	v, err := interfaceToUInt64(str, 8)
	return uint8(v), err
}
func InterfaceToUInt16(item interface{}) (uint16, error) {
	str := fmt.Sprintf("%.0f", item)
	v, err := interfaceToUInt64(str, 16)
	return uint16(v), err

}
func InterfaceToUInt32(item interface{}) (uint32, error) {
	str := fmt.Sprintf("%.0f", item)
	v, err := interfaceToUInt64(str, 32)
	return uint32(v), err

}
func InterfaceToUInt64(item interface{}) (uint64, error) {
	str := fmt.Sprintf("%.0f", item)
	v, err := interfaceToUInt64(str, 64)
	return uint64(v), err

}

func interfaceToUInt64(src string, bitSize int) (uint64, error) {
	return strconv.ParseUint(src, 0, bitSize)
}
func hexToUInt64(src string, bitSize int) (uint64, error) {
	return strconv.ParseUint(src, 16, bitSize)
}
