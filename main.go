package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/transformations/pkg/transform"
	"sigs.k8s.io/yaml"
)

func main() {

	// Read Test Data folder will paramertize for CLI
	// pass the list of unstructured/unstructured or JSON (TBD) to a pkg.
	// map[gk][]json
	// codify the rules for determining what needs to transformed
	// Determine the files that need to be skipped
	// for each file to be skipped create a wh_file_name file to be created.
	// This will allow the user to see what is being removed as well as choose that it is not removed.

	files, err := ioutil.ReadDir("./test-data")
	if err != nil {
		fmt.Printf("%v", err)
	}

	jsonArray := readFiles("./test-data", files)
	transform.OutputTransforms(jsonArray, transform.TransformOptions{
		Annotations:   map[string]string{"key": "value"},
		IsOpenshift:   true,
		StartingPath:  "./test-data",
		OutputDirPath: "./transforms",
	})

	// Now we will add code to apply the transforms. This should eventually be two commands

	// Read the "export" dir or test-data
	// For each file found, look for the corresponding whiteout file or the transform file
	// if whiteout remove skip the file
	// if transform apply the jsonpatch and save file in the same spot in the "output" dir
	// if nothing found keep the file the same.

	for _, file := range jsonArray {
		fmt.Printf("\n")
		fname, _ := transform.GetTransformPath("./transforms", "./test-data", file.Path)
		whfname, _ := transform.GetWhiteOutFilePath("./transforms", "./test-data", file.Path)

		// if white out alert user, and continue
		_, err := os.Stat(whfname)
		if !os.IsNotExist(err) {
			fmt.Printf("\nSkipping file: %v becuase it should be deleted", file.Path)
			continue
		}

		fmt.Printf("testing-testing")

		// Get transform
		patchesJSON, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		pa, err := jsonpatch.DecodePatch(patchesJSON)
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		// apply transformation
		output, err := pa.Apply(patchesJSON)
		if err != nil {
			fmt.Printf("err: %v", err)
		}

		// write file to output
		// Get fileoutput
		dir, newName := filepath.Split(file.Path)
		dir = strings.Replace(dir, "./test-data", "./output", 1)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Printf("err: %v", err)
		}
		err = ioutil.WriteFile(filepath.Join(dir, newName), output, 0664)
		if err != nil {
			fmt.Printf("err: %v", err)
		}
	}
}

func readFiles(path string, files []os.FileInfo) []transform.TransformFile {
	jsonFiles := []transform.TransformFile{}
	for _, file := range files {
		filePath := fmt.Sprintf("%v/%v", path, file.Name())
		if file.IsDir() {
			newFiles, err := ioutil.ReadDir(filePath)
			if err != nil {
				fmt.Printf("%v", err)
			}
			files := readFiles(filePath, newFiles)
			jsonFiles = append(jsonFiles, files...)
		} else {
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("%v", err)
			}
			json, err := yaml.YAMLToJSON(data)
			if err != nil {
				fmt.Printf("%v", err)
			}

			jsonFiles = append(jsonFiles, transform.TransformFile{
				FileInfo: file,
				Path:     filePath,
				JSONData: string(json),
			})
		}
	}
	return jsonFiles
}
