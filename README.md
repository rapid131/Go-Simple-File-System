# Thomas Hutton
## Project2 "file system" readme
This is a simple file system that uses a byte array known as var VirtualDisk [6010][1024]byte. All files and directories are encoded using gob encoding and pushed onto the byte array. 
## Superblock
The superblock inhabits VirtualDisk[0][:] and holds values for the inode offset (3) blockbitmapoffset (2) inodebitmapoffset (1) and datablockoffset (9)
## Inodes
The system begins with 120 inodes that are roughly 45 bytes in size each and inhabit the block space VirtualDisk[3][:] to VirtualDisk[8][:] taking up 5 blocks. The inodes follow this format:
type Inode struct {
	IsValid      bool
	IsDirectory  bool
	Datablocks   [4]int
	Filecreated  time.Time
	Filemodified time.Time
	Inodenumber  int
}
## Allocation Bitmap
The block bitmap is located at VirtualDisk[2][:] and the Inode bitmap is located at VirtualDisk[1][:]
## Root Directory and Directory Structure
The root directory is located at VirtualDisk[9][:] and my directory structure follows the format: 
type Directory struct {
	Filename  string
	Inode     int
	Files     []int
	Filenames []string
}
Directory.Files is an integer array that holds the inodes for flies, and Directory.Filenames is a string array that holds the names of those files. Files and their inode numbers are pushed onto the arrays in the same index so they are always tied to each other. Directory.Inode and Directory.Filename both belong to the directory itself.
## Open Function
The open function takes filename as a string, mode as a string, and parent directory integer as arguments in the form "Open(mode string, filename string, searchnode int)"
