package main
 
import (
	//"bufio"
	"fmt"
	//"io"
	"unsafe"
	"io/ioutil"
	"os"
	"strings"
)
 
func main() {
	if len(os.Args) != 4 {
		fmt.Println("lack of config file, eg: ${path_of_file} ${old_string} ${new_string}")
		os.Exit(-1)
	}
	EditDirandFile(os.Args[1], os.Args[2], os.Args[3])
}

func EditDirandFile(pathName string, ustr string, nstr string){
	rd, _ := ioutil.ReadDir(pathName)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathName+""+fi.Name())
			EditDirandFile(pathName + "/" +fi.Name() + "", ustr, nstr)
		} else {
			fmt.Println("file",pathName+"/"+fi.Name())
			if str, err := ioutil.ReadFile(pathName+"/"+fi.Name()); err == nil {
				new_str := strings.Replace(*(*string)(unsafe.Pointer(&str)), ustr, nstr, -1)
				if ioutil.WriteFile(pathName+"/"+fi.Name(), *(*[]byte)(unsafe.Pointer(&new_str)), os.ModePerm) == nil {
					fmt.Println("OK file",pathName+"/"+fi.Name())
				}
			}
		}
	}
}







