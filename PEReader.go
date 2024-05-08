package main

import (
	"fmt"
	"os"
    "encoding/json"
	"io/ioutil"
)

type PeReader struct{}

func (pr *PeReader) Read(name string, data []byte) map[string]interface{} {
	JSON := make(map[string]map[string]interface{})
	fileContent, err := ioutil.ReadFile("json/" + name + ".json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(fileContent, &JSON); err != nil {
		panic(err)
	}

	reader := make(map[string]interface{})
	for key, value := range JSON {
		typ := value["Type"].(string)
		switch typ {
		case "string":
			reader[key] = pr.string(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "short":
			reader[key] = pr.short(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "uint":
			reader[key] = pr.uint(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "byte":
			reader[key] = pr.byte(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "ushort":
			reader[key] = pr.ushort(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "array_short":
			reader[key] = pr.arrayShort(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		case "array_char":
			reader[key] = pr.arrayChar(data, int(value["Offset"].(float64)), int(value["Size"].(float64)))
		}
	}
	return reader
}

func (pr *PeReader) string(data []byte, offset, size int) string {
	return string(data[offset : offset+size])
}

func (pr *PeReader) short(data []byte, offset, size int) int {
	return int(int16(data[offset]) | int16(data[offset+1])<<8)
}

func (pr *PeReader) uint(data []byte, offset, size int) int {
	return int(uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24)
}

func (pr *PeReader) byte(data []byte, offset, size int) byte {
	return data[offset]
}

func (pr *PeReader) ushort(data []byte, offset, size int) int {
	return int(uint16(data[offset]) | uint16(data[offset+1])<<8)
}

func (pr *PeReader) arrayShort(data []byte, offset, size int) []int {
	var result []int
	for i := offset; i < offset+size; i += 2 {
		result = append(result, int(int16(data[i])|int16(data[i+1])<<8))
	}
	return result
}

func (pr *PeReader) arrayChar(data []byte, offset, size int) []byte {
	return data[offset : offset+size]
}

const (
    SizeOfDosHeader       = 0x40
    SizeOfFileHeader      = 0x18
    SizeOfOptionalHeader  = 0x60
    SizeOfDataDirectories = 0x80
    SizeOfSectionHeader   = 0x28
)

func main() {
    pr := &PeReader{} // Crée une instance de la struct PeReader
    data, err := os.ReadFile("debug.exe")
    if err != nil {
        panic(err)
    }
   
    old :=0
    new := SizeOfDosHeader

    DOS_HEADER := pr.Read("IMAGE_DOS_HEADER", data[old:new])
    e_lfane := DOS_HEADER["e_lfanew"].(int)    
    old = e_lfane  + 0x04
    new = e_lfane + SizeOfFileHeader

    FileHeader := pr.Read("IMAGE_FILE_HEADER",  data[old:new]) 
    old = new

    new = old + SizeOfOptionalHeader
    OPTIONAL_HEADER32 := pr.Read("IMAGE_OPTIONAL_HEADER32", data[old: new])
    
    old = new
    new = old + SizeOfDataDirectories
    DATA_DIRECTORIES:= pr.Read("IMAGE_DATA_DIRECTORIES", data[old: new])


    fmt.Println(DOS_HEADER)
    fmt.Println(FileHeader)
    fmt.Println(OPTIONAL_HEADER32)
    fmt.Println(DATA_DIRECTORIES)
    /*
      
        

      
        
        $sections = [];
        for ($i=0; $i < $this->FileHeader->NumberOfSections ; $i++) { 
            $calc = $sectionbase + ($this->SizeOfSectionHeader * $i);

            fseek($file, $calc);

            $SECTION_HEADER = fread($file,  $this->SizeOfOptionalHeader);
            $sections[] = $this->read("IMAGE_SECTION_HEADER", $SECTION_HEADER);
        }
        $this->Sections = $sections;


    */
	_, _ = os.Stdin.Read([]byte{0})
}