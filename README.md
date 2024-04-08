# Thomas Hutton
## Project2 "file system" readme
This is a simple file system that uses a byte array known as var VirtualDisk [6010][1024]byte. All files and directories are encoded using gob encoding and pushed onto the byte array. 
## Superblock
The superblock inhabits VirtualDisk[0][:] and holds values for the inode offset (3) blockbitmapoffset (2) inodebitmapoffset (1) and datablockoffset (9)
## Inodes
The system begins with 120 inodes that are roughly 45 bytes in size each and inhabit the block space VirtualDisk[3][:] to VirtualDisk[8][:] taking up 5 blocks. The inodes follow this format
```
type Inode struct {
	IsValid      bool,
	IsDirectory  bool,
	Datablocks   [4]int,
	Filecreated  time.Time,
	Filemodified time.Time,
	Inodenumber  int
}
```
The inodes are numbered from 0-119, the first inode is the null inode, the root directory is located at inode number 1 which is the second inode. Datablocks are 3 direct blocks and 1 indirect block. The blocks point directly to the real block number the data is located in, meaning the datablocks in the inode already have the datablockoffset added to them.
## Allocation Bitmap
The block bitmap is located at VirtualDisk[2][:] and the Inode bitmap is located at VirtualDisk[1][:]
## Root Directory and Directory Structure
The root directory is located at VirtualDisk[9][:] and my directory structure follows the format: 
```
"type Directory struct {
	Filename  string,
	Inode     int,
	Files     []int,
	Filenames []string,
}
```
Directory.Files is an integer array that holds the inodes for flies, and Directory.Filenames is a string array that holds the names of those files. Files and their inode numbers are pushed onto the arrays in the same index so they are always tied to each other. Directory.Inode and Directory.Filename both belong to the directory itself. Filenames are restricted in size by the Open system call. The root directory is located at the second inode which is inode number 1
## File Structure
Files in the system are based on the following struct 
```
type DirectoryEntry struct {
	Filename string
	Inode    int
	Fileinfo string
}
```
The filename string is limited in size by the Open function. The fileinfo string is the string which holds the data of the file. 
## Open Function
The Open function takes filename as a string, mode as a string, and parent directory integer as arguments in the form "Open(mode string, filename string, searchnode int)"
An example Open call might look like filesystem.Open("open","hello.txt",1) the root directory is located at inode number one so that is the integer entered.
The open function also supports "read", "write", and "append" modes. Write and append take user input to add a string to the file, while read prints the contents of a file.
## Write Function
The Write function will take filename as a string, fileinfo as a string, and an int for the node of the directory in the form "Write(filename string, searchnode int, fileinfo string)". An example call of this function might look like filesystem.Write("hello.txt", 1, bigstring)
## Read Function
The read function will take filename as a string and an int for the node of the directory in the form "Read(filename string, searchnode int)".
An example call of this function might look like filesystem.Read("hello.txt",1)
## Unlink Function
The unlink function will take filename as a string and an int for the node of the directory and delete the file from the disk. This call zeroes out the inode related to the file, deletes the file entries inside the directory struct, and updates the relevant bitmap portions to false. The function takes the form "Unlink(filename string, searchnode int)". An example call of this function might look like filesystem.Unlink("hello.txt",1)
## Example testing
Below is an example of testing you could do with my program. The Open call with mode "open" will create a file called "hello.txt" and the write call will write the "big" string to the files contents. The contents can then be read via the Read call
```
filesystem.InitializeDisk()
big := "Big string"
filesystem.Open("open", "hello.txt", 1)
filesystem.Write("hello.txt", 1, big)
filesystem.Read("hello.txt", 1)
filesystem.Unlink("hello.txt", 1)
filesystem.Read("hello.txt", 1)
```
you can also pull the inodes, print them, and write them back in using the following calls:
```
 inodes := filesystem.ReadInodesFromDisk()
 fmt.Println(inodes)
filesystem.WriteInodesToDisk(inodes)
```
