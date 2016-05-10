package main

import (
	"fmt"
	"github.com/Xustyx/goxymemory"
	"time"
)

func main() {
	dm := xymemmory.DataManager("StarCraft.exe")
	if !dm.IsOpen{
		fmt.Printf("Failed opening process.\n")
		return
	}

	for {
		//(READ/WRITE) String
		err := dm.Write(0X0059B420, xymemmory.Data{"TROLORO", xymemmory.STRING})
		if err != nil {
			fmt.Printf("Failed writing memory. %s", err)

		}
		data, err := dm.Read(0X0059B420, xymemmory.STRING)
		if err != nil {
			fmt.Printf("Failed reading memory. %s", err)

		}
		fmt.Println("STRING: ", data)

		//(READ/WRITE) Byte
		err = dm.Write(0X0059B420, xymemmory.Data{256, xymemmory.BYTE})
		if err != nil {
			fmt.Printf("Failed writing memory. %s", err)

		}
		data, err = dm.Read(0X0059B420 ,xymemmory.BYTE)
		if err != nil {
			fmt.Printf("Failed reading memory. %s", err)

		}
		fmt.Println("BYTE: ", data)


		//(READ/WRITE) Int
		err = dm.Write(0X0057F0F4, xymemmory.Data{-1, xymemmory.INT})
		if err != nil {
			fmt.Printf("Failed writing memory. %s", err)

		}
		data, err = dm.Read(0X0057F0F4, xymemmory.INT)
		if err != nil {
			fmt.Printf("Failed reading memory. %s", err)

		}
		fmt.Println("INT: ", data)

		//(READ/WRITE) Uint
		err = dm.Write(0X0057F0F0, xymemmory.Data{500, xymemmory.UINT})
		if err != nil {
			fmt.Printf("Failed writing memory. %s", err)

		}
		data, err = dm.Read(0X0057F0F0, xymemmory.UINT)
		if err != nil {
			fmt.Printf("Failed reading memory. %s", err)

		}
		fmt.Println("UINT: ", data)

		time.Sleep(1000 * time.Millisecond)
	}

	/*proc, _err := xymemmory.ProcessHandler("StarCraft.exe")
	if _err != nil {
		fmt.Printf("Error in processHandler: %s\n", _err)
		return
	}else {
		_err = proc.Open()
		if _err != nil {
			fmt.Printf("Error in processHandler: %s\n", _err)
			return
		}else {
			for {
				data,_err := proc.ReadBytes(0X0057F0F0, 4)
				if _err != nil {
					fmt.Printf("Error in processHandler: %s\n", _err)
				} else {
					fmt.Println("Data: ", data)
				}

				time.Sleep(1000 * time.Millisecond)
			}
		}
	}*/

	return
}
