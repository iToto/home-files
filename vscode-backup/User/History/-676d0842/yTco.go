package bulk-file-rename

// create a go application that renames all files in a directory to a new name with a sequenced number appended to the end
// the application should take in two arguments, the directory path and the new name
// the application should rename all files in the directory to the new name with a sequenced number appended to the end

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <directory> <new_name>")
		return
	}

	directory := os.Args[1]
	newName := os.Args[2]

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return
	}

	for i, file := range files {
		extension := filepath.Ext(file.Name())
		newFileName := fmt.Sprintf("%s_%d%s", newName, i+1, extension)
		oldPath := filepath.Join(directory, file.Name())
		newPath := filepath.Join(directory, newFileName)

		fmt.Printf("Renaming %s to %s\n", oldPath, newPath)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println("Failed to rename file:", err)
			return
		}
	}
}
