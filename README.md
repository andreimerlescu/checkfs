# CheckFS

Utility for checking File and Directory with Options.

## Installation

```bash
go get -u github.com/andreimerlescu/checkfs
```

This package is built for **Go 1.18+**!

## Usage

### Check File

```go
package main

import (
	"fmt"
	
	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

func main() {
	err := check.File("/path/to/file.txt", file.Options{
		Exists: true,
	})
	if err != nil {
		fmt.Printf("File validation failed: %v\n", err)
	} else {
		fmt.Println("File exists!")
	}
}

```

### Check Directory

```go
package main

import (
	"fmt"
	
	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/directory"
)

func main() {
	err := check.Directory("/path/to/directory", directory.Options{
		RequireWrite: true,
	})
	if err != nil {
		fmt.Printf("Directory validation failed: %v\n", err)
	} else {
		fmt.Println("Directory is writable!")
	}
}

```

## Configurations

### `file.Options`

| **Field**        | **Type**      | **Description**                                             |
|------------------|---------------|-------------------------------------------------------------|
| `ReadOnly`       | `bool`        | Check if the file is read-only                              |
| `RequireWrite`   | `bool`        | Check if the file is writable                               |
| `RequireOwner`   | `string`      | Ensure the file is owned by a specific user (UID as string) |
| `RequireGroup`   | `string`      | Ensure the file belongs to a specific group (GID as string) |
| `RequireBaseDir` | `string`      | Check if the file resides inside a specific base directory  |
| `CreatedBefore`  | `time.Time`   | Verify the file was created before a specific time          |
| `ModifiedBefore` | `time.Time`   | Verify the file was modified before a specific time         |
| `RequireExt`     | `string`      | Ensure the file has a specific extension                    |
| `RequirePrefix`  | `string`      | Ensure the file name begins with a specific prefix          |
| `IsLessThan`     | `int64`       | Verify the file size is less than this value                |
| `IsSize`         | `int64`       | Verify the file size matches this exact value               |
| `IsGreaterThan`  | `int64`       | Verify the file size is greater than this value             |
| `IsBaseNameLen`  | `int`         | Verify the file base name is exactly this length            |
| `IsFileMode`     | `os.FileMode` | Verify the file permissions match this mode                 |
| `WriteOnly`      | `bool`        | Check if the file is write-only                             |
| `Exists`         | `bool`        | Verify whether the file exists or not                       |
| `Create`         | `Create{}`    | Creates the resource.                                       | 


### `file.Create{}`

When you want to use `checkfs` to `Create` a new `File` or `Directory`, you can use:

| Property   | Type                     | Default                        |
|------------|--------------------------|--------------------------------|
| `Kind`     | `uint8`                  | `file.NoAction`                | 
| `FileMode` | `os.FileMode` / `uint32` | `0`                            |
| `OpenFlag` | `int`                    | `0`                            | 
| `Path`     | `string`                 | Uses path from original call\* | 
| `Size`     | `int64`                  | `0`                            | 

\*  See the usage of the `.Path` property in `file.Create{}`:

```go
package main

import (
    "fmt"
    "filepath"
    "os"
    
    check "github.com/andreimerlescu/checkfs"
    "github.com/andreimerlescu/checkfs/file"
)

func main() {
    // USING CLI ARG PARAM 1 AS PATH `go run . /my/path/as/args/1`
    path := os.Args[1]
	err := check.File(path, file.Options{
		file.Create: file.Create{
			Kind:     file.IfNotExists,
			OpenFlag: os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			FileMode: 0644,
		},
	})
    // OR 
    testFile := file.Create{
        Path: filepath.Join(os.TempDir(), "test-file.yaml"),
        OpenFlag: os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
        Kind: file.IfNotExists, // will only run if /tmp/test-file.yaml doesn't exist
        FileMode: 0644,
        // Size: 1776, // no size defined here, so it will be an empty file
    }
    err := testFile.Run()
    if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    // OR
    anotherFile := file.Create{
		Path: filepath.Join(os.TempDir(), "test-file.yaml"),
		OpenFlag: os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		Kind: file.IfExists, // deletes old empty file then creates new file
		FileMode: 0644,
        Size: 1776, // new file will be populated with 1776 bytes
    }
	err := anotherFile.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

### `directory.Options`

| **Field**        | **Type**    | **Description**                                                  |
|------------------|-------------|------------------------------------------------------------------|
| `ReadOnly`       | `bool`      | Check if the directory is read-only                              |
| `RequireWrite`   | `bool`      | Check if the directory is writable                               |
| `RequireOwner`   | `string`    | Ensure the directory is owned by a specific user (UID as string) |
| `RequireGroup`   | `string`    | Ensure the directory belongs to a specific group (GID as string) |
| `RequireBaseDir` | `string`    | Check if the directory resides inside a specific base directory  |
| `CreatedBefore`  | `time.Time` | Verify the directory was created before a specific time          |
| `ModifiedBefore` | `time.Time` | Verify the directory was modified before a specific time         |
| `RequirePrefix`  | `string`    | Ensure the directory name begins with a specific prefix          |
| `WillCreate`     | `bool`      | Verify ability to create the directory if it doesn't exist       |
| `Exists`         | `bool`      | Verify whether the directory exists or not                       |
| `Create`         | `Create{}`  | Creates the resource.                                            | 

### `directory.Create{}`

When you want to use `checkfs` to `Create` a new `File` or `Directory`, you can use: 

| Property   | Type                     | Default                        |
|------------|--------------------------|--------------------------------|
| `Kind`     | `uint8`                  | `file.NoAction`                | 
| `FileMode` | `os.FileMode` / `uint32` | `0`                            |
| `Path`     | `string`                 | Uses path from original call\* | 
| `Size`     | `int64`                  | `0`                            | 

\*  See the usage of the `.Path` property in `directory.Create{}`: 

```go
package main

import (
    "fmt"
    "os"
    
    check "github.com/andreimerlescu/checkfs"
    "github.com/andreimerlescu/checkfs/directory"
)

func main() {
    path := os.Args[1]
	err := check.Directory(path, directory.Options{
		directory.Create: directory.Create{
			Kind:     directory.IfNotExists,
			FileMode: 0755,
		},
	})
    if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
    // OR
    err := directory.Create{
		Kind: directory.IfNotExists, // only if the path doesnt exist shall this run
		Path: filepath.Join(os.TempDir(), "test-directory"), // create a new path
		FileMode: 0755, // set its mode to 0755 (standard) default mode is 0
    }
    // OR replace that directory again with...
	err := directory.Create{
		Kind: directory.IfExists, // deletes old directory first, then creates a new one
		Path: filepath.Join(os.TempDir(), "test-directory"), // reuse a path that already exists
		FileMode: 0755, // uses these permissions for the new directory
	}
}
```

The `directory.Create{}` struct has `.Run() error` exposed that you can run outside of the `.Check() error` func.

Throughout the `.Check() error` functionality, the `directory.Create{}` struct is processed in the `directory.Options{}`
structure, but the default `directory.Create.Kind` is `directory.NoAction` which is a `uint8` set to `0`. No actions
take by `.Run() error` are performed without `directory.NoAction` set to `0`. When you change this value, you are
telling `checkfs` that it is okay to **destroy** the path and its contents in a non-reversible manner.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

```plaintext
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

### Contributions

Contributions are welcome! Please fork this repository, make your changes, and submit a pull request.
