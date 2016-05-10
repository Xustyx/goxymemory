package xymemmory

import (
	"unsafe"
	"fmt"
	"errors"
	"github.com/Xustyx/w32"
)

type ProcessException error

type Process struct {
	Name string
	Pid uint32
}

type processHandler struct {
	process *Process
	hProcess uintptr
}

func ProcessHandler(processName string) (hProcess *processHandler, err ProcessException)  {
	_hProcess := processHandler{}
	_hProcess.process, err = processFromName(processName)

	return &_hProcess, err
}

func list() (processes []*Process) {
	processes = make([]*Process,0)

	handle := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPPROCESS, 0)
	if handle == 0 {
		fmt.Printf("Warning, CreateToolhelp32Snapshot failed. Error: ")
		return
	}

	var pEntry w32.PROCESSENTRY32
	PROCESSENTRY32_SIZE := unsafe.Sizeof(pEntry)
	pEntry.Size = uint32(PROCESSENTRY32_SIZE)

	_err := w32.Process32First(handle, &pEntry)
	if _err == nil {
		for {
			name := w32.UTF16PtrToString(&pEntry.ExeFile[0])
			processes = append(processes, &Process{ name,  pEntry.ProcessID})
			_err = w32.Process32Next(handle, &pEntry)
			if _err != nil {
				break
			}
		}
	} else {
		fmt.Printf("Warning, Process32First failed. Error: ", _err)
	}

	w32.CloseHandle(handle)

	return
}

func processFromName(processName string) ( *Process, ProcessException) {
	for _,process := range list() {
		if process.Name == processName {
			return  process, nil
		}
	}

	err := errors.New("Invalid process name.")
	return nil, err
}

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

func setDebugPrivilege() bool {
	pseudoHandle, _err := w32.GetCurrentProcess()
	if _err != nil {
		fmt.Printf("Warning, GetCurrentProcess failed. Error: ", _err)
		return false
	}

	hToken := w32.HANDLE(0)
	if !w32.OpenProcessToken(w32.HANDLE(pseudoHandle),  w32.TOKEN_ADJUST_PRIVILEGES | w32.TOKEN_QUERY, &hToken) {
		fmt.Printf("Warning, GetCurrentProcess failed.")
		return false
	}

	return setPrivilege(hToken, w32.SE_DEBUG_NAME, true)
}

func setPrivilege (hToken w32.HANDLE, lpszPrivilege string, bEnablePrivilege bool) bool {
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

	if !w32.AdjustTokenPrivileges(hToken, 0, &tPrivs, TOKEN_PRIVILEGES_SIZE, nil, nil){
		fmt.Printf("Warning, AdjustTokenPrivileges failed.")
		return false
	}

	return true
}

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