# Go Check FS

Utility for checking File and Directory with Options.

## Installation

```bash
go get -u github.com/andreimerlescu/go-checkfs
```

## Usage

### Check File

```go
package main

import (
	"fmt"
	check "github.com/andreimerlescu/go-checkfs/file"
)

func main() {
	err := check.File("/path/to/file.txt", check.Options{
		ReadOnly: true,
	})
	if err != nil {
		fmt.Printf("File validation failed: %v\n", err)
	} else {
		fmt.Println("File is read-only!")
	}
}

```

### Check Directory

```go
package main

import (
	"fmt"
	check "github.com/andreimerlescu/go-checkfs/directory"
)

func main() {
	err := check.Directory("/path/to/directory", check.Options{
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

| **Field**        | **Type** | **Description**                                              |
|------------------|----------|--------------------------------------------------------------|
| `ReadOnly`       | `bool`   | Check if the file is read-only.                              |
| `RequireWrite`   | `bool`   | Check if the file is writable.                               |
| `RequireOwner`   | `string` | Ensure the file is owned by a specific user (UID as string). |
| `RequireGroup`   | `string` | Ensure the file belongs to a specific group (GID as string). |
| `RequireBaseDir` | `string` | Check if the file resides inside a specific base directory.  |

### `directory.Options`

| **Field**        | **Type** | **Description**                                                   |
|------------------|----------|-------------------------------------------------------------------|
| `ReadOnly`       | `bool`   | Check if the directory is read-only.                              |
| `RequireWrite`   | `bool`   | Check if the directory is writable.                               |
| `RequireOwner`   | `string` | Ensure the directory is owned by a specific user (UID as string). |
| `RequireGroup`   | `string` | Ensure the directory belongs to a specific group (GID as string). |
| `RequireBaseDir` | `string` | Check if the directory resides inside a specific base directory.  |

## Test Results

### Common

Functions shared between `file.File` and `directory.Directory`

#### Unit Test

```log
=== RUN   TestIsPathInBase
--- PASS: TestIsPathInBase (0.00s)
=== RUN   TestIsPathInBase/Valid_path_in_base
    --- PASS: TestIsPathInBase/Valid_path_in_base (0.00s)
=== RUN   TestIsPathInBase/Path_outside_base
    --- PASS: TestIsPathInBase/Path_outside_base (0.00s)
=== RUN   TestIsPathInBase/Path_escaping_base
    --- PASS: TestIsPathInBase/Path_escaping_base (0.00s)
=== RUN   TestIsPathInBase/Empty_path
    --- PASS: TestIsPathInBase/Empty_path (0.00s)
=== RUN   TestIsPathInBase/Empty_base_directory
    --- PASS: TestIsPathInBase/Empty_base_directory (0.00s)
PASS

=== RUN   TestRelStartsWithParent
--- PASS: TestRelStartsWithParent (0.00s)
=== RUN   TestRelStartsWithParent/Relative_path_escapes
    --- PASS: TestRelStartsWithParent/Relative_path_escapes (0.00s)
=== RUN   TestRelStartsWithParent/Relative_path_inside
    --- PASS: TestRelStartsWithParent/Relative_path_inside (0.00s)
=== RUN   TestRelStartsWithParent/Current_directory
    --- PASS: TestRelStartsWithParent/Current_directory (0.00s)
=== RUN   TestRelStartsWithParent/Escaping_with_separator
    --- PASS: TestRelStartsWithParent/Escaping_with_separator (0.00s)
PAS
Process finished with the exit code 
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/go-checkfs/common
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkIsPathInBase
BenchmarkIsPathInBase-16    	    7075672	       162.7 ns/op
BenchmarkRelStartsWithParent
BenchmarkRelStartsWithParent-16    	53051914	   23.57 ns/op
PASS

Process finished with the exit code 0

```

### File

#### Unit Test

```log
=== RUN   TestFile
--- PASS: TestFile (0.00s)
=== RUN   TestFile/Valid_file
    --- PASS: TestFile/Valid_file (0.00s)
=== RUN   TestFile/Non-file_path
    --- PASS: TestFile/Non-file_path (0.00s)
=== RUN   TestFile/Invalid_base_directory
    --- PASS: TestFile/Invalid_base_directory (0.00s)
PASS

Process finished with the exit code 0
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/go-checkfs/file
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkFile
BenchmarkFile-16    	 1229371	       982.5 ns/op
PASS

Process finished with the exit code 0
```

### Directory

#### Unit Test

```log
=== RUN   TestDirectory
=== RUN   TestDirectory/Valid_directory
=== RUN   TestDirectory/Non-directory_path
=== RUN   TestDirectory/Invalid_base_directory
--- PASS: TestDirectory (0.00s)
    --- PASS: TestDirectory/Valid_directory (0.00s)
    --- PASS: TestDirectory/Non-directory_path (0.00s)
    --- PASS: TestDirectory/Invalid_base_directory (0.00s)
PASS

Process finished with the exit code 0
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/go-checkfs/directory
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkDirectory
BenchmarkDirectory-16    	 1318100	       907.8 ns/op
PASS

Process finished with the exit code 
```

## CheckFS Package

### Unit Test

```log
=== RUN   TestFile
--- PASS: TestFile (0.00s)
=== RUN   TestDirectory
--- PASS: TestDirectory (0.00s)
PASS

Process finished with the exit code 0
```

### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/go-checkfs
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkFile
BenchmarkFile-16         	 1235120	       957.5 ns/op
BenchmarkDirectory
BenchmarkDirectory-16    	 1365142	       879.2 ns/op
PASS

Process finished with the exit code 
```


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
