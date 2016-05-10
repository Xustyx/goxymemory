package xymemmory

import (
	"fmt"
	"errors"
	"encoding/binary"
)

type DataType int

const (
	UINT DataType = iota
	INT
	BYTE
	STRING
)

var data_types = [...]string {
	"uint",
	"int",
	"byte",
	"string",
}

func (data_type DataType) String() string {
	return data_types[data_type]
}

type DataException error

type Data struct {
	Value interface{}
	DataType DataType
}

type dataManager struct {
	ProcessName string
	process *processHandler
	IsOpen bool
}

func DataManager(processName string) *dataManager {
	_err := error(nil)
	dm := &dataManager{}
	dm.ProcessName = processName

	dm.process, _err = ProcessHandler(processName)
	if _err != nil {
		fmt.Errorf("Error in processHandler: %s\n", _err)

		return dm
	}

	_err = dm.process.Open()
	if _err != nil {
		fmt.Errorf("Error in processHandler Open: %s\n", _err)
		return dm
	} else {
		dm.IsOpen = true
	}

	return dm
}

func (dm *dataManager) Read(address uint, dataType DataType) (data Data, err DataException) {
	_err := error(nil)

	if !dm.IsOpen {
		err = errors.New("Process is not open.")
		return
	}

	switch dataType {
		case UINT:
			data, _err = dm.readUint(address)
		case INT:
			data, _err = dm.readInt(address)
		case BYTE:
			data, _err = dm.readByte(address)
		case STRING:
			data, _err = dm.readString(address)
		default:
			err = errors.New("Invalid data type.")
	}

	if _err != nil {
		fmt.Errorf("Error in processHandler Read: %s\n", _err)
	}

	return
}


func (dm *dataManager) readByte(address uint) (data Data, err ProcessException) {
	data.DataType = BYTE

	_data, err := dm.process.ReadBytes(address, 1)
	data.Value = _data[0]
	return
}

func (dm *dataManager) readString(address uint) (data Data, err ProcessException) {
	data.DataType = STRING

	 wordBytes := make([]byte, 0)
	_address := address

	for {
		_data, _err := dm.readByte(_address)
		if _err != nil{
			fmt.Errorf("Error in DataManager readByte: %s\n", _err)
			break
		}

		value :=  _data.Value.(byte)

		if value == 0 {
			break
		}

		_address += 0x01
		wordBytes = append(wordBytes, value)
	}

	data.Value = string(wordBytes[:])

	return
}

func (dm *dataManager) readInt(address uint) (data Data, err ProcessException)   {
	data.DataType = INT

	_data, err := dm.process.ReadBytes(address, 4)
	data.Value = int(binary.LittleEndian.Uint32(_data))
	return
}

func (dm *dataManager) readUint(address uint) (data Data, err ProcessException)   {
	data.DataType = UINT

	_data, err := dm.process.ReadBytes(address, 4)
	data.Value = binary.LittleEndian.Uint32(_data)
	return
}

func (dm *dataManager) Write(address uint, data Data) (err DataException){
	_err := error(nil)

	if !dm.IsOpen {
		err = errors.New("Process is not open.")
		return
	}

	switch data.DataType {
		case UINT:
			_err = dm.writeUint(address, uint(data.Value.(int)))
		case INT:
			_err = dm.writeInt(address, data.Value.(int))
		case BYTE:
			_err = dm.writeByte(address, byte(data.Value.(int)))
		case STRING:
			_err = dm.writeString(address, data.Value.(string))
		default:
			err = errors.New("Invalid data type.")
	}

	if _err != nil {
		fmt.Errorf("Error in processHandler Write: %s\n", _err)
	}

	return
}

func (dm *dataManager) writeByte(address uint, b byte) (err ProcessException) {
	data := []byte{b}

	err = dm.process.WriteBytes(address, data)

	return
}

func (dm *dataManager) writeString(address uint, str string) (err ProcessException) {
	data := []byte(str)

	err = dm.process.WriteBytes(address, data)

	return
}

func (dm *dataManager) writeInt(address uint, i int) (err ProcessException)   {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(i))

	err = dm.process.WriteBytes(address, data)

	return
}

func (dm *dataManager) writeUint(address uint, u uint) (err ProcessException)   {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(u))

	err = dm.process.WriteBytes(address, data)

	return
}
