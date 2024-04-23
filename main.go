/*
@author Thomas Hutton
This is a simple shell that utilizes my virtual file system for several commands. The shell can
use mkdir, mv, cp, cat (with redirect), rm, and more. The virtual filesystem can be navigated using
'cd >> directoryname'
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"project1/filesystem"
	"strings"
	"time"
)

func main() {
	filesystem.InitializeDisk()
	//got info about bufio and strings from here https://tutorialedge.net/golang/reading-console-input-golang/
	//create scanner
	scanner := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to shell, please enter commands")
	var workingdirectories []int
	currentworkingdirectory := 1
	workingdirectories = append(workingdirectories, currentworkingdirectory)
	//for loop, each iteration simulates one line of shell
	for {
		//create userinput using the scanner.ReadString, delimiter newline
		userinput, _ := scanner.ReadString('\n')
		list := strings.Fields(userinput)
		switch list[0] {
		//first case exit, exits the shell
		case "exit":
			os.Exit(0)
		//this removes a file from the virtual file system in current directory
		case "rm":
			filesystem.Unlink(list[1], currentworkingdirectory)
		//This uses the less command to page through a virtual file's string
		case "more":
			if filesystem.Read(list[1], currentworkingdirectory) == "" {
				fmt.Println("Could not find file or file empty")
			} else {
				cmd := exec.Command("less")
				cmd.Stdin = strings.NewReader(filesystem.Read(list[1], currentworkingdirectory))
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		//this prints a file in the real file system, virtual file system, or redirects
		//data from the main file system into the virtual one
		case "cat":
			if list[2] == ">>" {
				out, err := exec.Command("cat", list[1]).Output()
				if err != nil {
					log.Fatal(err)
				}
				newstring := string(out)
				filesystem.Open("open", list[3], currentworkingdirectory)
				filesystem.Write(list[3], currentworkingdirectory, newstring)
			} else if list[1] == ">>" {
				filesystem.Read(list[2], currentworkingdirectory)
			} else {
				command := exec.Command("ls", list[1:]...)
				out, err := command.Output()
				if err != nil {
					fmt.Println("Need valid args")
				} else {
					fmt.Println(string(out))
				}
			}
		//This function can cd into real system files or using >> can cd through the virtual file system
		case "cd":
			if len(list) > 1 {
				if list[1] == ">>" {
					if len(list) < 3 {
						currentworkingdirectory = 1
						fmt.Println("In directory root.dir")
					} else if list[2] == ".." {
						if currentworkingdirectory != 1 {
							workingdirectories = workingdirectories[:len(workingdirectories)-1]
							currentworkingdirectory = workingdirectories[len(workingdirectories)-1]
							inodes := filesystem.ReadInodesFromDisk()
							datablocks := inodes[currentworkingdirectory].Datablocks
							workingdirectory := filesystem.ReadFolder(datablocks[0], datablocks[1], datablocks[2], datablocks[3])
							fmt.Println("In directory ", workingdirectory.Filename)
						} else {
							fmt.Println("Already in top directory")
						}
					} else {
						found := false
						inodes := filesystem.ReadInodesFromDisk()
						workinginode := inodes[currentworkingdirectory]
						datablocks := workinginode.Datablocks
						workingfolder := filesystem.ReadFolder(datablocks[0], datablocks[1], datablocks[2], datablocks[3])
						for i := range workingfolder.Filenames {
							if workingfolder.Filenames[i] == list[2] {
								currentworkingdirectory = workingfolder.Files[i]
								workingdirectories = append(workingdirectories, currentworkingdirectory)
								found = true
								break
							}
						}
						if found {
							fmt.Println("In directory ", list[2])
						} else {
							fmt.Println("Could not find directory", list[2])
						}
					}
				} else {
					//cd plus a directory was typed
					os.Chdir(list[1])
				}
			} else {
				//Got this from here https://stackoverflow.com/questions/46028707/how-to-change-the-current-directory-in-go
				home, _ := os.UserHomeDir()
				err := os.Chdir(home)
				if err != nil {
					fmt.Println(err)
				}
			}
		//case whoami prints name
		case "whoami":
			fmt.Println("thutton2 Thomas Hutton")
		//case ls uses exec.Command to issue the ls command
		case "ls":
			// This command I understood from here https://stackoverflow.com/questions/22781788/how-could-i-pass-a-dynamic-set-of-arguments-to-gos-command-exec-command
			command := exec.Command("ls", list[1:]...)
			out, err := command.Output()
			if err != nil {
				fmt.Println("Need valid args")
			} else {
				fmt.Println(string(out))
			}
		//case wc uses exec.Command to issue wc command
		case "wc":
			command := exec.Command("wc", list[1:]...)
			out, err := command.Output()
			if err != nil {
				fmt.Println("Need valid args")
			} else {
				fmt.Println(string(out))
			}
		//case mkdir creates a new directory in the current directory in the virtual file system
		case "mkdir":
			var newdirectory filesystem.Directory
			var datablocks [4]int
			found := false
			superblock := filesystem.ReadSuperblock()
			newdirectory.Filename = list[1]
			inodes := filesystem.ReadInodesFromDisk()
			datablocks = inodes[currentworkingdirectory].Datablocks
			currentdirectory := filesystem.ReadFolder(datablocks[0], datablocks[1], datablocks[2], datablocks[3])
			for i := range currentdirectory.Filenames {
				if currentdirectory.Filenames[i] == list[1] {
					found = true
					break
				}
			}
			if !found {
				inodebitmap := filesystem.BytesToBools(filesystem.VirtualDisk[superblock.Inodebitmapoffset][:filesystem.EndInodeBitmap])
				blockBitmap := filesystem.BytesToBools(filesystem.VirtualDisk[superblock.Blockbitmapoffset][:filesystem.EndBlockBitmap])
				for i := range inodebitmap {
					if inodebitmap[i] == false {
						inodebitmap[i] = true
						newdirectory.Inode = i
						inodes[i].Filecreated = time.Now()
						inodes[i].Filemodified = time.Now()
						inodes[i].IsDirectory = true
						inodes[i].IsValid = true
						for j := range blockBitmap {
							if blockBitmap[j] == false {
								blockBitmap[j] = true
								inodes[i].Datablocks[0] = j + superblock.Datablocksoffset
								break
							}
						}
						break
					}
				}
				currentdirectory.Filenames = append(currentdirectory.Filenames, list[1])
				currentdirectory.Files = append(currentdirectory.Files, newdirectory.Inode)
				filesystem.AddWorkingDirectoryToDisk(currentdirectory, datablocks)
				filesystem.AddWorkingDirectoryToDisk(newdirectory, inodes[newdirectory.Inode].Datablocks)
				filesystem.AddBlockBitmapToDisk(blockBitmap)
				filesystem.AddInodeBitmapToDisk(inodebitmap)
				filesystem.WriteInodesToDisk(inodes)
				fmt.Println("Created new folder: ", newdirectory.Filename, " in folder: ", currentdirectory.Filename)
			} else {
				fmt.Println("Directory already exists")
			}
		//case cp copies a file in the working directory to a named child directory
		case "cp":
			filefound := false
			folderfound := false
			var fileinode int
			var folderinode int
			inodes := filesystem.ReadInodesFromDisk()
			workinginode := inodes[currentworkingdirectory]
			folderdatablocks := workinginode.Datablocks
			workingfolder := filesystem.ReadFolder(folderdatablocks[0], folderdatablocks[1], folderdatablocks[2], folderdatablocks[3])
			for i := range workingfolder.Filenames {
				if list[1] == workingfolder.Filenames[i] {
					filefound = true
					fileinode = workingfolder.Files[i]
					break
				}
			}
			for i := range workingfolder.Filenames {
				if list[2] == workingfolder.Filenames[i] {
					folderfound = true
					folderinode = workingfolder.Files[i]
					break
				}
			}
			if filefound && folderfound {
				newfolderinode := inodes[folderinode]
				oldfileinode := inodes[fileinode]
				foundfile := filesystem.DecodeDirectoryEntryFromDisk(oldfileinode)
				foundfolder := filesystem.ReadFolder(newfolderinode.Datablocks[0], newfolderinode.Datablocks[1], newfolderinode.Datablocks[2], newfolderinode.Datablocks[3])
				filesystem.Open("open", foundfile.Filename, newfolderinode.Inodenumber)
				filesystem.Write(foundfile.Filename, newfolderinode.Inodenumber, foundfile.Fileinfo)
				fmt.Println("Copied file: ", foundfile.Filename, " into directory: ", foundfolder.Filename)
			} else {
				fmt.Println("Could not locate file or folder")
			}
		//case mv moves a file in the current directory to a child directory
		case "mv":
			filefound := false
			folderfound := false
			var fileinodenumber int
			var folderinodenumber int
			inodes := filesystem.ReadInodesFromDisk()
			workinginode := inodes[currentworkingdirectory]
			folderdatablocks := workinginode.Datablocks
			workingfolder := filesystem.ReadFolder(folderdatablocks[0], folderdatablocks[1], folderdatablocks[2], folderdatablocks[3])
			for i := range workingfolder.Filenames {
				if list[1] == workingfolder.Filenames[i] {
					filefound = true
					fileinodenumber = workingfolder.Files[i]
					break
				}
			}
			for i := range workingfolder.Filenames {
				if list[2] == workingfolder.Filenames[i] {
					folderfound = true
					folderinodenumber = workingfolder.Files[i]
					break
				}
			}
			if filefound && folderfound {
				newfolderinode := inodes[folderinodenumber]
				newfolder := filesystem.ReadFolder(newfolderinode.Datablocks[0], newfolderinode.Datablocks[1], newfolderinode.Datablocks[2], newfolderinode.Datablocks[3])
				for i := range workingfolder.Filenames {
					if list[1] == workingfolder.Filenames[i] {
						workingfolder.Filenames[i] = ""
						workingfolder.Files[i] = 0
						break
					}
				}
				newfolder.Filenames = append(newfolder.Filenames, list[1])
				newfolder.Files = append(newfolder.Files, fileinodenumber)
				filesystem.AddWorkingDirectoryToDisk(workingfolder, workinginode.Datablocks)
				filesystem.AddWorkingDirectoryToDisk(newfolder, newfolderinode.Datablocks)
				fmt.Println("File: ", list[1], "moved to directory: ", newfolder.Filename)
			} else {
				fmt.Println("Could not locate file or folder")
			}
		//default returns "invalid command" string
		default:
			fmt.Println("Invalid Command")
		}

	}

}
