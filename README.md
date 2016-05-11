# GoXyMemory
[GoXyMemory][git-goxymemory-url] is a port of [XyMemory][git-xymemory-url] (Python) to GoLang.
This project allows read and write process memory using a [fork][git-w32-fork-url] of syscalls wrapper called
[w32][git-w32-url] from [AllenDang][git-allen-url].

# Example of use
First, we need download and install this [repo][git-goxymemory-url] with go commands and next, use the DataManager
class like the [example][git-goxymemory-example-url] below.

```
$> go get github.com/Xustyx/goxymemory
$> go get github.com/Xustyx/goxymemory
```

```go
    //Init
    dm := goxymemmory.DataManager("StarCraft.exe") 	        //Get the DataManager with the process passed.
	if !dm.IsOpen{				                            //Check if process was opened.
		fmt.Printf("Failed opening process.\n")
		return
	}

	//(READ/WRITE) String example.
	err := dm.Write(0X0059B420,
		goxymemmory.Data{"TROLORO", goxymemmory.STRING})	//Write a string.
	if err != nil {							                //Check if not failed.
		fmt.Printf("Failed writing memory. %s", err)
	}
	data, err := dm.Read(0X0059B420, goxymemmory.STRING)	//Reads the string.
	if err != nil {						                    //Check if not failed.
		fmt.Printf("Failed reading memory. %s", err)
	}
	fmt.Println("STRING: ", data)					        //Print the value.

	//(READ/WRITE) Byte example.
	err = dm.Write(0X0059B420,
		goxymemmory.Data{256, goxymemmory.BYTE})		    //Write a byte.
	if err != nil {							                //Check if not failed.
		fmt.Printf("Failed writing memory. %s", err)
	}
	data, err = dm.Read(0X0059B420 ,goxymemmory.BYTE)		//Reads the byte.
	if err != nil {							                //Check if not failed.
		fmt.Printf("Failed reading memory. %s", err)
	}
	fmt.Println("BYTE: ", data)				            	//Print the value.


	//(READ/WRITE) Int example.
	err = dm.Write(0X0057F0F4,
		goxymemmory.Data{-1, goxymemmory.INT})			     //Write an int.
	if err != nil {							                 //Check if not failed.
		fmt.Printf("Failed writing memory. %s", err)

	}
	data, err = dm.Read(0X0057F0F4, goxymemmory.INT)		//Reads the int.
	if err != nil {						                	//Check if not failed.
		fmt.Printf("Failed reading memory. %s", err)
	}
	fmt.Println("INT: ", data)					            //Print the value.

	//(READ/WRITE) Uint example.
	err = dm.Write(0X0057F0F0,
		goxymemmory.Data{500, goxymemmory.UINT})	    	//Write an uint.
	if err != nil {							                //Check if not failed.
		fmt.Printf("Failed writing memory. %s", err)
	}
	data, err = dm.Read(0X0057F0F0, goxymemmory.UINT)		//Reads the uint.
	if err != nil {						                	//Check if not failed.
		fmt.Printf("Failed reading memory. %s", err)
	}
	fmt.Println("UINT: ", data)					            //Print the value.
```

# Actually supported types and methods
### Types
- BYTE: 'byte'
- STRING: 'string'
- INT: 'int'
- UINT: 'uint'

### Methods (DataManager)
* DataManager(processName string)
  * processName: The process name to handle.
* read(address uint, dataType DataType)
  * address: Memory address in hexadecimal format.
  * dataType: The enum value of desired Type.
* write(address uint, data Data)
  * address: Memory address in hexadecimal format.
  * data: Struct that contains the data to add and the enum value of desired Type.

#Disclaimer
The author can not be held liable for any use of this code.

[git-goxymemory-example-url]: <https://github.com/Xustyx/goxymemory/tree/master/example>
[git-goxymemory-url]: <https://github.com/Xustyx/goxymemory>
[git-xymemory-url]: <https://github.com/Xustyx/xymemory>
[git-w32-url]: <https://github.com/AllenDang/w32>
[git-allen-url]: <https://github.com/AllenDang/w32>
[git-w32-fork-url]: <https://github.com/Xustyx/w32>