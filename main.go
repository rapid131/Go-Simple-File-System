/*
@author Thomas Hutton
This is a simple shell that takes user input and can perform cd, ls, whoami,
wc, mkdir, cp, and mv commands from the OS. Can also type exit to exit the
shell. cd and whoami are run natively from this program while the rest are
run throuh the exec.Command function from os/exec
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
		//got this partially from here https://stackoverflow.com/questions/28705716/paging-output-from-go
		case "rm":
			if list[1] == ">>" {
				filesystem.Unlink(list[2], currentworkingdirectory)
			}
		case "more":
			if list[1] == ">>" {
				cmd := exec.Command("less")
				cmd.Stdin = strings.NewReader(filesystem.Read(list[2], currentworkingdirectory))
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
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
		case "cd":
			//only cd was typed
			if list[1] == ">>" && len(list) > 1 {
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
				if len(list) < 2 {
					//Got this from here https://stackoverflow.com/questions/46028707/how-to-change-the-current-directory-in-go
					home, _ := os.UserHomeDir()
					err := os.Chdir(home)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					//cd plus a directory was typed
					os.Chdir(list[1])
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
		//case mkdir uses exec.Command to issue mkdir command
		case "mkdir":
			if list[1] == ">>" {
				var newdirectory filesystem.Directory
				var datablocks [4]int
				found := false
				superblock := filesystem.ReadSuperblock()
				newdirectory.Filename = list[2]
				inodes := filesystem.ReadInodesFromDisk()
				datablocks = inodes[currentworkingdirectory].Datablocks
				currentdirectory := filesystem.ReadFolder(datablocks[0], datablocks[1], datablocks[2], datablocks[3])
				for i := range currentdirectory.Filenames {
					if currentdirectory.Filenames[i] == list[2] {
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
					currentdirectory.Filenames = append(currentdirectory.Filenames, list[2])
					currentdirectory.Files = append(currentdirectory.Files, newdirectory.Inode)
					filesystem.AddWorkingDirectoryToDisk(currentdirectory, datablocks)
					filesystem.AddWorkingDirectoryToDisk(newdirectory, inodes[newdirectory.Inode].Datablocks)
					filesystem.AddBlockBitmapToDisk(blockBitmap)
					filesystem.AddInodeBitmapToDisk(inodebitmap)
					filesystem.WriteInodesToDisk(inodes)
				} else {
					fmt.Println("Directory already exists")
				}
			} else {
				command := exec.Command("mkdir", list[1:]...)
				out, err := command.Output()
				if err != nil {
					fmt.Println("Need valid args")
				} else {
					fmt.Println(string(out))
				}
			}
		//case cp uses exec.Command to issue cp command
		case "cp":
			command := exec.Command("cp", list[1:]...)
			out, err := command.Output()
			if err != nil {
				fmt.Println("Need valid args")
			} else {
				fmt.Println(string(out))
			}
		//case mv uses exec.Command to issue mv command
		case "mv":
			command := exec.Command("mv", list[1:]...)
			out, err := command.Output()
			if err != nil {
				fmt.Println("Need valid args")
			} else {
				fmt.Println(string(out))
			}
		//default returns "invalid command" string
		default:
			fmt.Println("Invalid Command")
		}

	}

}
