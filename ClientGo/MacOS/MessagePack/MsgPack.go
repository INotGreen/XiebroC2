package MessagePack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"encoding/binary"
)

type MsgPackType int

const (
	Unknown  MsgPackType = 0
	Null     MsgPackType = 1
	Map      MsgPackType = 2
	Array    MsgPackType = 3
	String   MsgPackType = 4
	Integer  MsgPackType = 5
	UInt64   MsgPackType = 6
	Boolean  MsgPackType = 7
	Float    MsgPackType = 8
	Single   MsgPackType = 9
	DateTime MsgPackType = 10
	Binary   MsgPackType = 11
)

type MsgPack struct {
	name       string
	lowerName  string
	innerValue interface{}
	parent     *MsgPack
	children   []*MsgPack
	valueType  MsgPackType
}

type MsgPackEnum struct {
	children []*MsgPack
	position int
}

func (e *MsgPackEnum) Current() interface{} {
	return e.children[e.position]
}

func (e *MsgPackEnum) MoveNext() bool {
	e.position++
	return e.position < len(e.children)
}

func (e *MsgPackEnum) Reset() {
	e.position = -1
}

func (m *MsgPack) ValueType() MsgPackType {
	return m.valueType
}

func (m *MsgPack) ForcePathObject(path string) *MsgPack {
	tmpParent := m
	pathList := strings.FieldsFunc(path, func(r rune) bool {
		return r == '.' || r == '/' || r == '\\'
	})
	if len(pathList) == 0 {
		return nil
	}
	if len(pathList) > 1 {
		for _, p := range pathList[:len(pathList)-1] {
			tmpObject := tmpParent.FindObject(p)
			if tmpObject == nil {
				tmpParent = tmpParent.InnerAddMapChild()
				tmpParent.SetName(p)
			} else {
				tmpParent = tmpObject
			}
		}
	}
	last := pathList[len(pathList)-1]
	j := tmpParent.IndexOf(last)
	if j > -1 {
		return tmpParent.children[j]
	}
	tmpParent = tmpParent.InnerAddMapChild()
	tmpParent.SetName(last)
	return tmpParent
}
func (m *MsgPack) IndexOf(name string) int {
	tmp := strings.ToLower(name)
	for i, child := range m.children {
		if tmp == child.lowerName {
			return i
		}
	}
	return -1
}
func (m *MsgPack) SetName(value string) {
	m.name = value
	m.lowerName = strings.ToLower(value)
}
func (m *MsgPack) InnerAddMapChild() *MsgPack {
	if m.valueType != Map {
		m.Clear()
		m.valueType = Map
	}
	return m.InnerAdd()
}
func (m *MsgPack) InnerAdd() *MsgPack {
	r := &MsgPack{
		parent: m,
	}
	m.children = append(m.children, r)
	return r
}

func (m *MsgPack) SetAsString(value string) {
	m.innerValue = value
	m.valueType = String
}

func (m *MsgPack) AsString() string {
	return m.GetAsString()
}

func (m *MsgPack) GetAsString() string {
	if m.innerValue == nil {
		return ""
	}
	return m.innerValue.(string)
}

func (m *MsgPack) Clear() {
	for _, child := range m.children {
		child.Clear()
	}
	m.children = nil
}

func (m *MsgPack) InnerAddArrayChild() *MsgPack {
	if m.valueType != Array {
		m.Clear()
		m.valueType = Array
	}
	return m.InnerAdd()
}

func (m *MsgPack) AddArrayChild() *MsgPack {
	return m.InnerAddArrayChild()
}

func (m *MsgPack) FindObject(name string) *MsgPack {
	i := m.IndexOf(name)
	if i == -1 {
		return nil
	}
	return m.children[i]
}

func (m *MsgPack) GetEnumerator() *MsgPackEnum {
	return &MsgPackEnum{
		children: m.children,
		position: -1,
	}
}

//ReadTools

func GetUtf8Bytes(s string) []byte {
	if s == "" {
		return nil
	}

	return []byte(s)
}

func WriteNull(ms *bytes.Buffer) {
	ms.WriteByte(0xC0)
}

func WriteString(ms *bytes.Buffer, strVal string) {
	rawBytes := []byte(strVal)
	lenBytes := make([]byte, 0)

	len := len(rawBytes)
	if len <= 31 {
		b := byte(0xA0 + len)
		ms.WriteByte(b)
	} else if len <= 255 {
		b := byte(0xD9)
		ms.WriteByte(b)
		b = byte(len)
		ms.WriteByte(b)
	} else if len <= 65535 {
		b := byte(0xDA)
		ms.WriteByte(b)
		lenBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(lenBytes, uint16(len))
		ms.Write(lenBytes)
	} else {
		b := byte(0xDB)
		ms.WriteByte(b)
		lenBytes = make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(len))
		ms.Write(lenBytes)
	}

	ms.Write(rawBytes)
}

func WriteBinary(ms *bytes.Buffer, rawBytes []byte) {
	lenBytes := make([]byte, 0)

	len := len(rawBytes)
	if len <= 255 {
		b := byte(0xC4)
		ms.WriteByte(b)
		b = byte(len)
		ms.WriteByte(b)
	} else if len <= 65535 {
		b := byte(0xC5)
		ms.WriteByte(b)
		lenBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(lenBytes, uint16(len))
		ms.Write(lenBytes)
	} else {
		b := byte(0xC6)
		ms.WriteByte(b)
		lenBytes = make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(len))
		ms.Write(lenBytes)
	}

	ms.Write(rawBytes)
}

func WriteFloat(ms *bytes.Buffer, fVal float64) {
	ms.WriteByte(0xCB)
	binary.Write(ms, binary.BigEndian, math.Float64bits(fVal))
}

func WriteSingle(ms *bytes.Buffer, fVal float32) {
	ms.WriteByte(0xCA)
	binary.Write(ms, binary.BigEndian, math.Float32bits(fVal))
}

func WriteBoolean(ms *bytes.Buffer, bVal bool) {
	if bVal {
		ms.WriteByte(0xC3)
	} else {
		ms.WriteByte(0xC2)
	}
}

func WriteUInt64(ms *bytes.Buffer, iVal uint64) {
	ms.WriteByte(0xCF)
	binary.Write(ms, binary.BigEndian, iVal)
}

func WriteInteger(ms *bytes.Buffer, iVal int64) {
	if iVal >= 0 {
		// Positive integer
		if iVal <= 127 {
			ms.WriteByte(byte(iVal))
		} else if iVal <= 255 {
			ms.WriteByte(0xCC)
			ms.WriteByte(byte(iVal))
		} else if iVal <= 65535 {
			ms.WriteByte(0xCD)
			lenBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(lenBytes, uint16(iVal))
			ms.Write(lenBytes)
		} else if iVal <= 0xFFFFFFFF {
			ms.WriteByte(0xCE)
			lenBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(lenBytes, uint32(iVal))
			ms.Write(lenBytes)
		} else {
			ms.WriteByte(0xD3)
			lenBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(lenBytes, uint64(iVal))
			ms.Write(lenBytes)
		}
	} else {
		// Negative integer
		if iVal >= int64(math.MinInt32) {
			ms.WriteByte(0xD2)
			lenBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(lenBytes, uint32(iVal))
			ms.Write(lenBytes)
		} else if iVal >= int64(math.MinInt16) {
			ms.WriteByte(0xD1)
			lenBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(lenBytes, uint16(iVal))
			ms.Write(lenBytes)
		} else if iVal >= -128 {
			ms.WriteByte(0xD0)
			ms.WriteByte(byte(iVal))
		} else {
			ms.WriteByte(byte(iVal))
		}
	}
}

func (mp *MsgPack) WriteMap(buf *bytes.Buffer) {
	b := byte(0)
	lenBytes := make([]byte, 4)
	len := len(mp.children)
	if len <= 15 {
		b = byte(0x80 + byte(len))
		buf.WriteByte(b)
	} else if len <= 65535 {
		b = 0xDE
		buf.WriteByte(b)
		binary.BigEndian.PutUint16(lenBytes, uint16(len))
		buf.Write(lenBytes)
	} else {
		b = 0xDF
		buf.WriteByte(b)
		binary.BigEndian.PutUint32(lenBytes, uint32(len))
		buf.Write(lenBytes)
	}

	for _, child := range mp.children {
		WriteString(buf, child.name)
		child.Encode2Stream(buf)
	}
}

func (mp *MsgPack) WriteArray(buf *bytes.Buffer) {
	b := byte(0)
	lenBytes := make([]byte, 4)
	len := len(mp.children)
	if len <= 15 {
		b = byte(0x90 + byte(len))
		buf.WriteByte(b)
	} else if len <= 65535 {
		b = 0xDC
		buf.WriteByte(b)
		binary.BigEndian.PutUint16(lenBytes, uint16(len))
		buf.Write(lenBytes)
	} else {
		b = 0xDD
		buf.WriteByte(b)
		binary.BigEndian.PutUint32(lenBytes, uint32(len))
		buf.Write(lenBytes)
	}

	for _, child := range mp.children {
		child.Encode2Stream(buf)
	}
}

func GetString(utf8Bytes []byte) string {
	return string(utf8Bytes)
}

func ReadString(ms io.Reader, length int) (string, error) {
	rawBytes := make([]byte, length)
	_, err := ms.Read(rawBytes)
	if err != nil {
		return "", err
	}
	return GetString(rawBytes), nil
}

func ReadStringWithFlag(ms io.Reader) (string, error) {
	strFlag := make([]byte, 1)
	_, err := ms.Read(strFlag)
	if err != nil {
		return "", err
	}
	return ReadStringWithBytes(strFlag[0], ms)
}

func SwapBytes(v []byte) []byte {
	r := make([]byte, len(v))
	j := len(v) - 1
	for i := 0; i < len(r); i++ {
		r[i] = v[j]
		j--
	}
	return r
}

func SwapInt64(v int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return SwapBytes(buf)
}

func ReadStringWithBytes(strFlag byte, ms io.Reader) (string, error) {
	var len int
	var rawBytes []byte
	if strFlag >= 0xA0 && strFlag <= 0xBF {
		len = int(strFlag) - 0xA0
	} else if strFlag == 0xD9 {
		lenByte := make([]byte, 1)
		_, err := ms.Read(lenByte)
		if err != nil {
			return "", err
		}
		len = int(lenByte[0])
	} else if strFlag == 0xDA {
		rawBytes = make([]byte, 2)
		_, err := ms.Read(rawBytes)
		if err != nil {
			return "", err
		}
		rawBytes = SwapBytes(rawBytes)
		len = int(binary.BigEndian.Uint16(rawBytes))
	} else if strFlag == 0xDB {
		rawBytes = make([]byte, 4)
		_, err := ms.Read(rawBytes)
		if err != nil {
			return "", err
		}
		rawBytes = SwapBytes(rawBytes)
		len = int(binary.BigEndian.Uint32(rawBytes))
	} else {
		return "", errors.New("unsupported strFlag value")
	}

	rawBytes = make([]byte, len)
	_, err := ms.Read(rawBytes)
	if err != nil {
		return "", err
	}

	return GetString(rawBytes), nil
}

// SetAsInteger sets the innerValue as an integer and valueType as Integer
func (m *MsgPack) SetAsInteger(value int64) {
	m.innerValue = value
	m.valueType = Integer
}

// SetAsUInt64 sets the innerValue as an uint64 and valueType as UInt64
func (m *MsgPack) SetAsUInt64(value uint64) {
	m.innerValue = value
	m.valueType = UInt64
}

// GetAsUInt64 gets the innerValue as an uint64
func (m *MsgPack) GetAsUInt64() uint64 {
	switch m.valueType {
	case Integer:
		return uint64(m.innerValue.(int64))
	case UInt64:
		return m.innerValue.(uint64)
	case String:
		v, _ := strconv.ParseUint(strings.TrimSpace(m.innerValue.(string)), 10, 64)
		return v
	case Float:
		return uint64(m.innerValue.(float64))
	case Single:
		return uint64(m.innerValue.(float32))
	case DateTime:
		// Needs additional handling
	default:
		return 0
	}
	return 0
}

// GetAsInteger gets the innerValue as an int64
func (m *MsgPack) GetAsInteger() int64 {
	switch m.valueType {
	case Integer:
		return m.innerValue.(int64)
	case UInt64:
		return int64(m.innerValue.(uint64))
	case String:
		v, _ := strconv.ParseInt(strings.TrimSpace(m.innerValue.(string)), 10, 64)
		return v
	case Float:
		return int64(m.innerValue.(float64))
	case Single:
		return int64(m.innerValue.(float32))
	case DateTime:
		// Needs additional handling
	default:
		return 0
	}
	return 0
}

// GetAsFloat gets the innerValue as a float64
func (m *MsgPack) GetAsFloat() float64 {
	switch m.valueType {
	case Integer:
		return float64(m.innerValue.(int64))
	case String:
		v, _ := strconv.ParseFloat(m.innerValue.(string), 64)
		return v
	case Float:
		return m.innerValue.(float64)
	case Single:
		return float64(m.innerValue.(float32))
	case DateTime:
		// Needs additional handling
	default:
		return 0
	}
	return 0
}

// SetAsNull sets the innerValue as nil and valueType as Null
func (m *MsgPack) SetAsNull() {
	m.innerValue = nil
	m.valueType = Null
}

// SetAsString sets the innerValue as a string and valueType as String
func (m *MsgPack) SetAsStringA(value string) {
	m.innerValue = value
	m.valueType = String
}

// GetAsString gets the innerValue as a string
func (m *MsgPack) GetAsStringA() string {
	if m.innerValue == nil {
		return ""
	}

	return fmt.Sprintf("%v", m.innerValue)
}

// SetAsBoolean sets the innerValue as a boolean and valueType as Boolean
func (m *MsgPack) SetAsBoolean(bVal bool) {
	m.innerValue = bVal
	m.valueType = Boolean
}

// SetAsSingle sets the innerValue as a float32 and valueType as Single
func (m *MsgPack) SetAsSingle(fVal float32) {
	m.innerValue = fVal
	m.valueType = Single
}

// SetAsFloat sets the innerValue as a float64 and valueType as Float
func (m *MsgPack) SetAsFloat(fVal float64) {
	m.innerValue = fVal
	m.valueType = Float
}

// SetAsBytes sets the innerValue as a byte array and valueType as Binary
func (m *MsgPack) SetAsBytes(value []byte) {
	m.innerValue = value
	m.valueType = Binary
}

// GetAsBytes gets the innerValue as a byte array
func (m *MsgPack) GetAsBytes() []byte {
	switch m.valueType {
	case Integer:
		return int64ToBytes(m.innerValue.(int64))
	case String:
		return []byte(m.innerValue.(string))
	case Float:
		return float64ToBytes(m.innerValue.(float64))
	case Single:
		return float32ToBytes(m.innerValue.(float32))
	case DateTime:
		// Needs additional handling
	case Binary:
		return m.innerValue.([]byte)
	default:
		return []byte{}
	}
	return []byte{}
}

// Conversion functions go here. For example:
func int64ToBytes(i int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes()
}

//DecodeFromStrem

func (mp *MsgPack) DecodeFromStream(ms io.Reader) error {
	lvByte := make([]byte, 1)
	_, err := ms.Read(lvByte)
	if err != nil {
		return err
	}

	var rawByte []byte
	var msgPack *MsgPack
	var len, i int

	switch {
	case lvByte[0] <= 0x7F:
		mp.SetAsInteger(int64(lvByte[0]))
	case (lvByte[0] >= 0x80) && (lvByte[0] <= 0x8F):
		mp.Clear()
		mp.valueType = Map
		len = int(lvByte[0]) - 0x80
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			name, err := ReadStringWithFlag(ms)
			if err != nil {
				return err
			}
			msgPack.SetName(name)
			err = msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case (lvByte[0] >= 0x90) && (lvByte[0] <= 0x9F):
		mp.Clear()
		mp.valueType = Array
		len = int(lvByte[0]) - 0x90
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			err := msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case (lvByte[0] >= 0xA0) && (lvByte[0] <= 0xBF):
		len = int(lvByte[0]) - 0xA0
		str, err := ReadString(ms, len)
		if err != nil {
			return err
		}
		mp.SetAsString(str)
	case (lvByte[0] >= 0xE0) && (lvByte[0] <= 0xFF):
		mp.SetAsInteger(int64(int8(lvByte[0])))
	case lvByte[0] == 0xC0:
		mp.SetAsNull()
	case lvByte[0] == 0xC1:
		return errors.New("(never used) type $c1")
	case lvByte[0] == 0xC2:
		mp.SetAsBoolean(false)
	case lvByte[0] == 0xC3:
		mp.SetAsBoolean(true)
	case lvByte[0] == 0xC4:
		lenByte := make([]byte, 1)
		_, err := ms.Read(lenByte)
		if err != nil {
			return err
		}
		len := int(lenByte[0])
		rawByte = make([]byte, len)
		_, err = ms.Read(rawByte)
		if err != nil {
			return err
		}
		mp.SetAsBytes(rawByte)
	case lvByte[0] == 0xC5:
		rawByte = make([]byte, 2)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		len := int(binary.BigEndian.Uint16(rawByte))
		rawByte = make([]byte, len)
		_, err = ms.Read(rawByte)
		if err != nil {
			return err
		}
		mp.SetAsBytes(rawByte)
	case lvByte[0] == 0xC6:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		len := int(binary.BigEndian.Uint32(rawByte))
		rawByte = make([]byte, len)
		_, err = ms.Read(rawByte)
		if err != nil {
			return err
		}
		mp.SetAsBytes(rawByte)
	case (lvByte[0] == 0xC7) || (lvByte[0] == 0xC8) || (lvByte[0] == 0xC9):
		return errors.New("(ext8,ext16,ex32) type $c7,$c8,$c9")
	case lvByte[0] == 0xCA:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsSingle(float32FromBytes(rawByte))
	case lvByte[0] == 0xCB:
		rawByte = make([]byte, 8)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		math.Float32frombits(binary.LittleEndian.Uint32(rawByte))
	case lvByte[0] == 0xCC:
		lvByte = make([]byte, 1)
		_, err := ms.Read(lvByte)
		if err != nil {
			return err
		}
		mp.SetAsInteger(int64(lvByte[0]))
	case lvByte[0] == 0xCD:
		rawByte = make([]byte, 2)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsInteger(int64(binary.BigEndian.Uint16(rawByte)))
	case lvByte[0] == 0xCE:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsInteger(int64(binary.BigEndian.Uint32(rawByte)))
	case lvByte[0] == 0xCF:
		rawByte = make([]byte, 8)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsUInt64(binary.BigEndian.Uint64(rawByte))
	case lvByte[0] == 0xDC:
		rawByte = make([]byte, 2)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}

		rawByte = SwapBytes(rawByte)
		len := int(binary.BigEndian.Uint16(rawByte))
		mp.Clear()
		mp.valueType = Array
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			err := msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case lvByte[0] == 0xDD:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		len := int(binary.BigEndian.Uint16(rawByte))
		mp.Clear()
		mp.valueType = Array
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			err := msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case lvByte[0] == 0xD9:
		str, err := ReadStringWithBytes(lvByte[0], ms)
		if err != nil {
			return err
		}
		mp.SetAsString(str)
	case lvByte[0] == 0xDE:
		rawByte = make([]byte, 2)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		len := int(binary.BigEndian.Uint16(rawByte))
		mp.Clear()
		mp.valueType = Map
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			name, err := ReadStringWithFlag(ms)
			if err != nil {
				return err
			}
			msgPack.SetName(name)
			err = msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case lvByte[0] == 0xDF:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		len := int(binary.BigEndian.Uint32(rawByte))
		mp.Clear()
		mp.valueType = Map
		for i = 0; i < len; i++ {
			msgPack = mp.InnerAdd()
			name, err := ReadStringWithFlag(ms)
			if err != nil {
				return err
			}
			msgPack.SetName(name)
			err = msgPack.DecodeFromStream(ms)
			if err != nil {
				return err
			}
		}
	case lvByte[0] == 0xDA:
		str, err := ReadStringWithBytes(lvByte[0], ms)
		if err != nil {
			return err
		}
		mp.SetAsString(str)
	case lvByte[0] == 0xDB:
		str, err := ReadStringWithBytes(lvByte[0], ms)
		if err != nil {
			return err
		}
		mp.SetAsString(str)
	case lvByte[0] == 0xD0:
		byteVal := make([]byte, 1)
		_, err := ms.Read(byteVal)
		if err != nil {
			return err
		}
		mp.SetAsInteger(int64(int8(byteVal[0])))
	case lvByte[0] == 0xD1:
		rawByte = make([]byte, 2)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsInteger(int64(binary.BigEndian.Uint16(rawByte)))
	case lvByte[0] == 0xD2:
		rawByte = make([]byte, 4)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsInteger(int64(binary.BigEndian.Uint32(rawByte)))
	case lvByte[0] == 0xD3:
		rawByte = make([]byte, 8)
		_, err := ms.Read(rawByte)
		if err != nil {
			return err
		}
		rawByte = SwapBytes(rawByte)
		mp.SetAsInteger(int64(binary.BigEndian.Uint64(rawByte)))
	}

	return nil
}

func (mp *MsgPack) DecodeFromBytes(data []byte) {
	// In Go, the bytes.Buffer can be used similar to MemoryStream in C#.
	buf := bytes.NewBuffer(data)
	mp.DecodeFromStream(buf)
	// No need to set the position as 0 as the bytes.Buffer will start reading
	// from the start of the provided bytes.
}

func float32ToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func float64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}
func (mp *MsgPack) Encode2Bytes() []byte {
	var buf bytes.Buffer
	mp.Encode2Stream(&buf)
	return buf.Bytes()
}

func (mp *MsgPack) Encode2Stream(buf *bytes.Buffer) {
	switch mp.valueType {
	case Unknown, Null:
		WriteNull(buf)
	case String:
		WriteString(buf, mp.innerValue.(string))
	case Integer:
		WriteInteger(buf, mp.innerValue.(int64))
	case UInt64:
		WriteUInt64(buf, mp.innerValue.(uint64))
	case Boolean:
		WriteBoolean(buf, mp.innerValue.(bool))
	case Float:
		WriteFloat(buf, mp.innerValue.(float64))
	case Single:
		WriteSingle(buf, mp.innerValue.(float32))
	case DateTime:
		WriteInteger(buf, mp.GetAsInteger())
	case Binary:
		WriteBinary(buf, mp.innerValue.([]byte))
	case Map:
		mp.WriteMap(buf)
	case Array:
		mp.WriteArray(buf)
	default:
		WriteNull(buf)
	}
}
