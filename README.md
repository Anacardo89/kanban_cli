# Kanboards - a Kanban CLI app  
![image](demo.gif)  

Kanboards is a kanban CLI app designed to be customizable and portable.  

### Install:  
For rn, just clone the repo and `go get` the dependencies, i'm working on getting a vendor directory going  

### Usage:  
Commands are written on screen, apart from that:  
- the first time you run kanboards a directory with the same name will be created in your home directory  
- import/export from/to YAML with `i` and `e` when at main screen  
- YAML file is placed in the kanboards directory  
- you can only import into an empty DB, if there are aleady projects loaded `i` will do nothing  
