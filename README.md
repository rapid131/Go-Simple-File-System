# Thomas Hutton
## Project3 "putting it all together" readme
This is a working shell that uses my virtual filesystem to perform basic file maintenance tasks. All the functions will work on a file which exists in the current working directory, how to change the working directory and make new ones are explained in the "cd" and "mkdir" sections below.
## "rm" command
The rm command will take a filename and remove the file from the filesystem such as `rm hello.txt`
## "more" command
The more command will take a filename and present the file contents in a pager. The pages can be filed through using spacebar and 'q' will exit the pager and return to the shell. command should appear as `more hello.txt`
## "cat" command
the cat command is a special command that can work with the real and virtual filesystem. To redirect text from a .txt file in the real file system to a new file in the virtual file system, the user must type `cat longfile.txt >> hello.txt` where longfile.txt is a file in the real file system and hello.txt is the file the user wishes to create in the working directory of the virtual file system. The user can also put the redirect directly after "cat" such as `cat >> hello.txt` to print out the text from a virtual file.
## "cd" command
cd is another special command that works with the virtual and real file system. To use it as normally with the real file system, simply type `cd directory` or `cd ..` to navigate the real directories. For the virtual file system, one must use the redirect operator such as `cd >> folder.dir`. The user may also use `cd >> ..` and `cd >>` to traverse the virtual file system in a similar way to the real directory system.
## "whoami" command
This command simply returns the string I have provided it, in this case "thutton2 Thomas Hutton"
## "ls" and "wc" commands
These commands continue to only use the real file system and do not interact with the virtual file system.
## "mkdir" command
the mkdir command will create a new folder in the virtual file system inside the current working directory. use should appear as `mkdir folder.dir`
## "cp" command
The cp command works with the virtual file system. This command will take a file name in the current working directory and a folder inside the current working directory and copy the file to the child directory chosen. Usage should appear as `cp hello.txt folder.dir`
## "mv" command
The mv command works with the virtual file system. This command will take a file name in the current working directory and a folder inside the current working directory and move the file to the child directory chosen. Usage should appear as `mv hello.txt folder.dir`

## Example testing
Below is an example of testing you could do with my program. 
```
cat longfile.txt >> hello.txt
mkdir folder.dir
mv hello.txt folder.dir
cd >> folder.dir
more hello.txt
rm hello.txt
cd >> ..
```
you can also pull the inodes, print them, and write them back in using the following calls:
```
 inodes := filesystem.ReadInodesFromDisk()
 fmt.Println(inodes)
filesystem.WriteInodesToDisk(inodes)
```
