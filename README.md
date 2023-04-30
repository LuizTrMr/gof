# Gof
Under the dowhateveryouwantwithit license

## Usage
```
-exclude string
  Folders/files to ignore while searching, separated by a comma
-go
  Use go routines to search files
-path string
  Folder/file to search for the search term (default: ".")
-st string
  Term to be searched
```
> The `-go` flag runs a go routine for each file being searched. If there are small files in your folder, using this option results in a slower execution time. In the other hand, if your files are very big, this option can be better.

---
## **Easy search**: You can just type `gof searchterm` and use all other defaults, like this:
```bash
gof "your searchterm"
```
> But, if some other flag is used, you must pass in the `-st` flag for the searchterm

---
## Examples
```bash
gof -st "your searchterm" -exclude ".git,./folder or file"
```
```bash
gof -st "your searchterm" -go
```
```bash
gof -st "your searchterm" -path "./filename"
```
OR
```bash
gof -st "your searchterm" -path "./foldername"
```