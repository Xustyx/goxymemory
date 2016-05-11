package main

import (
	"fmt"
	"github.com/Xustyx/goxymemory"
	"time"
)

//This is and example of use goxymemory.
func main() {
	dm := goxymemmory.DataManager("StarCraft.exe") 	//Get the DataManager with the process passed.
	if !dm.IsOpen{					//Check if process was opened.
		fmt.Printf("Failed opening process.\n")
		return
	}

	for {
		//(READ/WRITE) String example.
		err := dm.Write(0X0059B420,
			goxymemmory.Data{"TROLORO", goxymemmory.STRING})	//Write a string.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed writing memory. %s", err)
		}
		data, err := dm.Read(0X0059B420, goxymemmory.STRING)		//Reads the string.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed reading memory. %s", err)
		}
		fmt.Println("STRING: ", data)					//Print the value.

		//(READ/WRITE) Byte example.
		err = dm.Write(0X0059B420,
			goxymemmory.Data{256, goxymemmory.BYTE})		//Write a byte.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed writing memory. %s", err)
		}
		data, err = dm.Read(0X0059B420 ,goxymemmory.BYTE)		//Reads the byte.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed reading memory. %s", err)
		}
		fmt.Println("BYTE: ", data)					//Print the value.
		
		//(READ/WRITE) Int example.
		err = dm.Write(0X0057F0F4,
			goxymemmory.Data{-1, goxymemmory.INT})			//Write an int.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed writing memory. %s", err)
		}
		data, err = dm.Read(0X0057F0F4, goxymemmory.INT)		//Reads the int.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed reading memory. %s", err)
		}
		fmt.Println("INT: ", data)					//Print the value.

		//(READ/WRITE) Uint example.
		err = dm.Write(0X0057F0F0,
			goxymemmory.Data{500, goxymemmory.UINT})		//Write an uint.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed writing memory. %s", err)
		}
		data, err = dm.Read(0X0057F0F0, goxymemmory.UINT)		//Reads the uint.
		if err != nil {							//Check if not failed.
			fmt.Printf("Failed reading memory. %s", err)
		}
		fmt.Println("UINT: ", data)					//Print the value.

		time.Sleep(1000 * time.Millisecond)				//Wait a second and repeat.
	} //This loops runs, and runs, and runs... until ctrl+c.


	// This is another example using directly processHandler.
	/*proc, _err := goxymemmory.ProcessHandler("StarCraft.exe")
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
