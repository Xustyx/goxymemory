//The MIT License (MIT)
//
//Copyright (c) 2016 Xustyx
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

//This package implements functions to reads and writes process memory more easily.
//
//Remember to execute with administrator privileges to grant debug on other process.
package goxymemmory

import (
	"encoding/binary"
	"errors"
	"fmt"
)

//Type of the data.
type DataType int

//Enum of data types.
const (
	UINT DataType = iota
	INT
	BYTE
	STRING
)

//String representation of data types.
var data_types = [...]string{
	"uint",
	"int",
	"byte",
	"string",
}

//Get the string value from enum value.
func (data_type DataType) String() string {
	return data_types[data_type]
}

//Exception type of DataManager.
type DataException error

//This type warp the read and write values.
type Data struct {
	Value    interface{} //Any type value.
	DataType DataType    //Unwarp value.
}

//This type is the Facade for read and write.
type dataManager struct {
	ProcessName string          //Name of the process.
	process     *processHandler //This handles the low level facade.
	IsOpen      bool            //True if we are in process.
}

//Constructor of DataManager
//Param	  (processName)	: The name of process to handle.
//Returns (*dataManager): A dataManager object.
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

//Facade to Read methods.
//Public method of (dataManager) class.
//Param	  (address) : The process memory addres in hexadecimal. EX: (0X0057F0F0).
//Param   (dataType): The type of data that want to retrieve.
//Returns (data)    : The data from memory. If low level facade fails, this will be nil.
//Errors  (err)	    : This will be not nil if handle is not opened or the type is invalid.
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

//Specific method for read a byte.
func (dm *dataManager) readByte(address uint) (data Data, err ProcessException) {
	data.DataType = BYTE

	_data, err := dm.process.ReadBytes(address, 1)
	data.Value = _data[0]
	return
}

//Specific method for read a String.
func (dm *dataManager) readString(address uint) (data Data, err ProcessException) {
	data.DataType = STRING

	wordBytes := make([]byte, 0)
	_address := address

	for {
		_data, _err := dm.readByte(_address)
		if _err != nil {
			fmt.Errorf("Error in DataManager readByte: %s\n", _err)
			break
		}

		value := _data.Value.(byte)

		if value == 0 {
			break
		}

		_address += 0x01
		wordBytes = append(wordBytes, value)
	}

	data.Value = string(wordBytes[:])

	return
}

//Specific method for read an int.
func (dm *dataManager) readInt(address uint) (data Data, err ProcessException) {
	data.DataType = INT

	_data, err := dm.process.ReadBytes(address, 4)
	data.Value = int(binary.LittleEndian.Uint32(_data))
	return
}

//Specific method for read an uint.
func (dm *dataManager) readUint(address uint) (data Data, err ProcessException) {
	data.DataType = UINT

	_data, err := dm.process.ReadBytes(address, 4)
	data.Value = binary.LittleEndian.Uint32(_data)
	return
}

//Facade to Write methods.
//Public method of (dataManager) class.
//Param	  (address) : The process memory addres in hexadecimal. EX: (0X0057F0F0).
//Param   (data)    : The data to write.
//Errors  (err)	    : This will be not nil if handle is not opened or the type is invalid.
func (dm *dataManager) Write(address uint, data Data) (err DataException) {
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

//Specific method for write a byte.
func (dm *dataManager) writeByte(address uint, b byte) (err ProcessException) {
	data := []byte{b}

	err = dm.process.WriteBytes(address, data)

	return
}

//Specific method for write a string.
func (dm *dataManager) writeString(address uint, str string) (err ProcessException) {
	data := []byte(str)

	err = dm.process.WriteBytes(address, data)

	return
}

//Specific method for write an int.
func (dm *dataManager) writeInt(address uint, i int) (err ProcessException) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(i))

	err = dm.process.WriteBytes(address, data)

	return
}

//Specific method for write an uint.
func (dm *dataManager) writeUint(address uint, u uint) (err ProcessException) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(u))

	err = dm.process.WriteBytes(address, data)

	return
}
