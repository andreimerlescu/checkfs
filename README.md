# CheckFS

Utility for checking File and Directory with Options.

## Installation

```bash
go get -u github.com/andreimerlescu/checkfs
```

## Usage

### Check File

```go
package main

import (
	"fmt"
	check "github.com/andreimerlescu/checkfs/file"
)

func main() {
	err := check.File("/path/to/file.txt", check.Options{
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
	check "github.com/andreimerlescu/checkfs/directory"
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
PASS
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/checkfs/common
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkIsPathInBase
BenchmarkIsPathInBase-16           	 1000000	      1137 ns/op
BenchmarkRelStartsWithParent
BenchmarkRelStartsWithParent-16    	 7390912	       159.9 ns/op
PASS
```

### File

#### Unit Test

```log
=== RUN   TestFile
--- PASS: TestFile (0.00s)
=== RUN   TestFile/Valid_regular_file
    --- PASS: TestFile/Valid_regular_file (0.00s)
=== RUN   TestFile/Non-existent_file_with_Exists=false
    --- PASS: TestFile/Non-existent_file_with_Exists=false (0.00s)
=== RUN   TestFile/Non-existent_file_with_Exists=true
    --- PASS: TestFile/Non-existent_file_with_Exists=true (0.00s)
=== RUN   TestFile/Directory_path
    --- PASS: TestFile/Directory_path (0.00s)
=== RUN   TestFile/Valid_base_directory
    --- PASS: TestFile/Valid_base_directory (0.00s)
=== RUN   TestFile/Invalid_base_directory
    --- PASS: TestFile/Invalid_base_directory (0.00s)
=== RUN   TestFile/Valid_extension
    --- PASS: TestFile/Valid_extension (0.00s)
=== RUN   TestFile/Invalid_extension
    --- PASS: TestFile/Invalid_extension (0.00s)
=== RUN   TestFile/Valid_prefix
    --- PASS: TestFile/Valid_prefix (0.00s)
=== RUN   TestFile/Invalid_prefix
    --- PASS: TestFile/Invalid_prefix (0.00s)
=== RUN   TestFile/Valid_creation_time
    --- PASS: TestFile/Valid_creation_time (0.00s)
=== RUN   TestFile/Invalid_creation_time
    --- PASS: TestFile/Invalid_creation_time (0.00s)
=== RUN   TestFile/Valid_modification_time
    --- PASS: TestFile/Valid_modification_time (0.00s)
=== RUN   TestFile/Invalid_modification_time
    --- PASS: TestFile/Invalid_modification_time (0.00s)
=== RUN   TestFile/Valid_exact_size
    --- PASS: TestFile/Valid_exact_size (0.00s)
=== RUN   TestFile/Invalid_exact_size
    --- PASS: TestFile/Invalid_exact_size (0.00s)
=== RUN   TestFile/Valid_size_less_than
    --- PASS: TestFile/Valid_size_less_than (0.00s)
=== RUN   TestFile/Invalid_size_less_than
    --- PASS: TestFile/Invalid_size_less_than (0.00s)
=== RUN   TestFile/Valid_size_greater_than
    --- PASS: TestFile/Valid_size_greater_than (0.00s)
=== RUN   TestFile/Invalid_size_greater_than
    --- PASS: TestFile/Invalid_size_greater_than (0.00s)
=== RUN   TestFile/Valid_base_name_length
    --- PASS: TestFile/Valid_base_name_length (0.00s)
=== RUN   TestFile/Invalid_base_name_length
    --- PASS: TestFile/Invalid_base_name_length (0.00s)
=== RUN   TestFile/Valid_read-only
    --- PASS: TestFile/Valid_read-only (0.00s)
=== RUN   TestFile/Valid_write_required
    --- PASS: TestFile/Valid_write_required (0.00s)
=== RUN   TestFile/Valid_write-only
    --- PASS: TestFile/Valid_write-only (0.00s)
=== RUN   TestFile/Valid_file_mode
    --- PASS: TestFile/Valid_file_mode (0.00s)
=== RUN   TestFile/Invalid_file_mode
    --- PASS: TestFile/Invalid_file_mode (0.00s)
=== RUN   TestFile/Valid_symlink
    --- PASS: TestFile/Valid_symlink (0.00s)
=== RUN   TestFile/Symlink_with_valid_base_dir
    --- PASS: TestFile/Symlink_with_valid_base_dir (0.00s)
=== RUN   TestFile/Multiple_valid_conditions
    --- PASS: TestFile/Multiple_valid_conditions (0.00s)
=== RUN   TestFile/Multiple_conditions_with_one_invalid
    --- PASS: TestFile/Multiple_conditions_with_one_invalid (0.00s)
PASS
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/checkfs/file
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkFile
BenchmarkFile/BasicChecks
BenchmarkFile/BasicChecks-16         	 1345123	       882.3 ns/op
BenchmarkFile/ExtensiveChecks
BenchmarkFile/ExtensiveChecks-16     	 1100227	      1053 ns/op
PASS
```

### Directory

#### Unit Test

```log
=== RUN   TestDirectory
--- PASS: TestDirectory (0.00s)
=== RUN   TestDirectory/Valid_existing_directory
    --- PASS: TestDirectory/Valid_existing_directory (0.00s)
=== RUN   TestDirectory/Non-existent_directory_with_Exists=false
    --- PASS: TestDirectory/Non-existent_directory_with_Exists=false (0.00s)
=== RUN   TestDirectory/Non-existent_directory_with_Exists=true
    --- PASS: TestDirectory/Non-existent_directory_with_Exists=true (0.00s)
=== RUN   TestDirectory/Non-directory_path
    --- PASS: TestDirectory/Non-directory_path (0.00s)
=== RUN   TestDirectory/Will_create_in_existing_parent
    --- PASS: TestDirectory/Will_create_in_existing_parent (0.00s)
=== RUN   TestDirectory/Will_create_with_existing_target
    --- PASS: TestDirectory/Will_create_with_existing_target (0.00s)
=== RUN   TestDirectory/Will_create_without_existence_check
    --- PASS: TestDirectory/Will_create_without_existence_check (0.00s)
=== RUN   TestDirectory/Will_create_and_require_existence
    --- PASS: TestDirectory/Will_create_and_require_existence (0.00s)
=== RUN   TestDirectory/Will_create_without_existence_check#01
    --- PASS: TestDirectory/Will_create_without_existence_check#01 (0.00s)
=== RUN   TestDirectory/Valid_base_directory
    --- PASS: TestDirectory/Valid_base_directory (0.00s)
=== RUN   TestDirectory/Invalid_base_directory
    --- PASS: TestDirectory/Invalid_base_directory (0.00s)
=== RUN   TestDirectory/Valid_prefix
    --- PASS: TestDirectory/Valid_prefix (0.00s)
=== RUN   TestDirectory/Invalid_prefix
    --- PASS: TestDirectory/Invalid_prefix (0.00s)
=== RUN   TestDirectory/Valid_creation_time
    --- PASS: TestDirectory/Valid_creation_time (0.00s)
=== RUN   TestDirectory/Invalid_creation_time
    --- PASS: TestDirectory/Invalid_creation_time (0.00s)
=== RUN   TestDirectory/Valid_modification_time
    --- PASS: TestDirectory/Valid_modification_time (0.00s)
=== RUN   TestDirectory/Invalid_modification_time
    --- PASS: TestDirectory/Invalid_modification_time (0.00s)
=== RUN   TestDirectory/Read-only_directory_check
    --- PASS: TestDirectory/Read-only_directory_check (0.00s)
=== RUN   TestDirectory/Write_permission_check
    --- PASS: TestDirectory/Write_permission_check (0.00s)
=== RUN   TestDirectory/Invalid_write_permission
    --- PASS: TestDirectory/Invalid_write_permission (0.00s)
=== RUN   TestDirectory/Valid_owner
    --- PASS: TestDirectory/Valid_owner (0.00s)
=== RUN   TestDirectory/Invalid_owner
    --- PASS: TestDirectory/Invalid_owner (0.00s)
=== RUN   TestDirectory/Valid_group
    --- PASS: TestDirectory/Valid_group (0.00s)
=== RUN   TestDirectory/Invalid_group
    --- PASS: TestDirectory/Invalid_group (0.00s)
=== RUN   TestDirectory/Multiple_valid_conditions
    --- PASS: TestDirectory/Multiple_valid_conditions (0.00s)
=== RUN   TestDirectory/Multiple_conditions_with_one_invalid
    --- PASS: TestDirectory/Multiple_conditions_with_one_invalid (0.00s)
PASS
```

#### Benchmark Test

```log
goos: linux
goarch: amd64
pkg: github.com/andreimerlescu/checkfs/directory
cpu: Intel(R) Xeon(R) W-3245 CPU @ 3.20GHz
BenchmarkDirectory
BenchmarkDirectory/BasicChecks
BenchmarkDirectory/BasicChecks-16         	 1414285	       853.2 ns/op
BenchmarkDirectory/ExtensiveChecks
BenchmarkDirectory/ExtensiveChecks-16     	  493394	      2062 ns/op
PASS
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
