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

package goxymemmory

import (
	"errors"
	"fmt"
	"github.com/Xustyx/w32"
	"unsafe"
)

//Exception type of ProcessHandler.
type ProcessException error

//Type of simple process.
type Process struct {
	Name string
	Pid  uint32
}

//This type handles the process.
type processHandler struct {
	process  *Process
	hProcess uintptr
}

//Constructor of ProcessHandler
//Param	  (processName)	   : The name of process to handle.
//Returns (*processHandler): A processHandler object.
//Errors  (err)		   : Error if don't exist process with passed name.
func ProcessHandler(processName string) (hProcess *processHandler, err ProcessException) {
	_hProcess := processHandler{}
	_hProcess.process, err = processFromName(processName)

	return &_hProcess, err
}

//This function returns a list of process.
func list() (processes []*Process) {
	processes = make([]*Process, 0)

	handle := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPPROCESS, 0)
	if handle == 0 {
		fmt.Printf("Warning, CreateToolhelp32Snapshot failed. Error: ")
		return
	}

	var pEntry w32.PROCESSENTRY32
	PROCESSENTRY32_SIZE := unsafe.Sizeof(pEntry)
	pEntry.Size = uint32(PROCESSENTRY32_SIZE)

	_err := w32.Process32First(handle, &pEntry) //Read frist element.
	if _err == nil {
		for {
			name := w32.UTF16PtrToString(&pEntry.ExeFile[0])
			processes = append(processes, &Process{name, pEntry.ProcessID})
			_err = w32.Process32Next(handle, &pEntry)
			if _err != nil {
				break
			}
		} //Loops until reach last process.
	} else {
		fmt.Printf("Warning, Process32First failed. Error: ", _err)
	}

	w32.CloseHandle(handle)

	return
}

//This function search a process with passed name in list() and returns it.
func processFromName(processName string) (*Process, ProcessException) {
	for _, process := range list() {
		if process.Name == processName {
			return process, nil
		}
	}

	err := errors.New("Invalid process name.")
	return nil, err
}

//Open the process of ProcessHandler in get self debug privileges.
//Public method of (processHandler) class.
//Errors (err): Error if don't exist process or cannot open with PAA.
func (ph *processHandler) Open() (err ProcessException) {

	if ph.process == nil {
		err = errors.New("The selected process does not exist")
		return
	}

	setDebugPrivilege()

	handle, _err := w32.OpenProcess(w32.PROCESS_ALL_ACCESS, false, ph.process.Pid)
	if _err != nil {
		err = errors.New("Cannot open this process. Reason: " + _err.Error())
		return
	}

	ph.hProcess = uintptr(handle)
	return
}

//This function try to set self process with debug privileges.
func setDebugPrivilege() bool {
	pseudoHandle, _err := w32.GetCurrentProcess()
	if _err != nil {
		fmt.Printf("Warning, GetCurrentProcess failed. Error: ", _err)
		return false
	}

	hToken := w32.HANDLE(0)
	if !w32.OpenProcessToken(w32.HANDLE(pseudoHandle), w32.TOKEN_ADJUST_PRIVILEGES|w32.TOKEN_QUERY, &hToken) {
		fmt.Printf("Warning, GetCurrentProcess failed.")
		return false
	}

	return setPrivilege(hToken, w32.SE_DEBUG_NAME, true)
}

//This function try to set privileges to a process.
func setPrivilege(hToken w32.HANDLE, lpszPrivilege string, bEnablePrivilege bool) bool {
	tPrivs := w32.TOKEN_PRIVILEGES{}
	TOKEN_PRIVILEGES_SIZE := uint32(unsafe.Sizeof(tPrivs))
	luid := w32.LUID{}

	if !w32.LookupPrivilegeValue(string(""), lpszPrivilege, &luid) {
		fmt.Printf("Warning, LookupPrivilegeValue failed.")
		return false
	}

	tPrivs.PrivilegeCount = 1
	tPrivs.Privileges[0].Luid = luid

	if bEnablePrivilege {
		tPrivs.Privileges[0].Attributes = w32.SE_PRIVILEGE_ENABLED
	} else {
		tPrivs.Privileges[0].Attributes = 0
	}

	if !w32.AdjustTokenPrivileges(hToken, 0, &tPrivs, TOKEN_PRIVILEGES_SIZE, nil, nil) {
		fmt.Printf("Warning, AdjustTokenPrivileges failed.")
		return false
	}

	return true
}

//This function search a module inside process.
func (ph *processHandler) GetModuleFromName(module string) (uintptr, error) {
	var (
		me32 w32.MODULEENTRY32
		snap w32.HANDLE
	)

	snap = w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE|w32.TH32CS_SNAPMODULE32, ph.process.Pid)
	me32.Size = uint32(unsafe.Sizeof(me32))

	for ok := w32.Module32First(snap, &me32); ok; ok = w32.Module32Next(snap, &me32) {
		szModule := w32.UTF16PtrToString(&me32.SzModule[0])

		if szModule == module {
			return (uintptr)(unsafe.Pointer(me32.ModBaseAddr)), nil
		}
	}

	return (uintptr)(unsafe.Pointer(me32.ModBaseAddr)), errors.New("module not found")
}

//Low level facade to Read memory.
//Public method of (processHandler) class.
//Param	  (address): The process memory addres in hexadecimal. EX: (0X0057F0F0).
//Param   (size)   : The size of bytes that we want to read.
//Returns (data)   : A byte array with data.
//Errors  (err)	   : This will be not nil if handle is not opened or cannot read the memory.
func (ph *processHandler) ReadBytes(address uint, size uint) (data []byte, err ProcessException) {
	if ph.hProcess == 0 {
		err = errors.New("No process handle.")
	}

	data, _err := w32.ReadProcessMemory(w32.HANDLE(ph.hProcess), uint32(address), size)
	if _err != nil {
		err = errors.New("Error reading memory. Reason: " + _err.Error())
	}

	return
}

//Low level facade to Write memory.
//Public method of (processHandler) class.
//Param	  (address): The process memory addres in hexadecimal. EX: (0X0057F0F0).
//Param   (data)   : A byte array with data.
//Errors  (err)	   : This will be not nil if handle is not opened or cannot write the memory.
func (ph *processHandler) WriteBytes(address uint, data []byte) (err ProcessException) {
	if ph.hProcess == 0 {
		err = errors.New("No process handle.")
	}

	_err := w32.WriteProcessMemory(w32.HANDLE(ph.hProcess), uint32(address), data, uint(len(data)))
	if _err != nil {
		err = errors.New("Error writing memory. Reason: " + _err.Error())
	}

	return
}
