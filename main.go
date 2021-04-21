package main

import (
  "log"
  "os"
  "path/filepath"
  "regexp"
  "strings"
)

func main() {

  if len(os.Args) < 1 {
    os.Exit(0)
  }

  searchDir := os.Args[1]

  processAll(searchDir)

}

func processAll(searchDir string) {
  err := filepath.Walk(searchDir, processFiles)
  if err != nil {
    log.Println(err)
  }
}

func processFiles(path string, f os.FileInfo, err error) error {

  if path == "." || f.IsDir() {
    return err
  }

  newPath := getNewPath(path)
  if newPath != path {
    rename(path, newPath)
  }

  return err

}

func rename(prevName, newName string) {
  err := os.Rename(prevName, newName)
  if err != nil {
    log.Println(err)
  }
}

func getNewPath(path string) string {

  // get the dir
  dir, fileName := filepath.Split(path)

  // get the ext
  ext := filepath.Ext(path)

  // get the base
  baseName := getBasename(fileName)

  newName := format(baseName) + ext
  newPath := filepath.Join(dir, newName)

  return newPath

}

func getBasename(fileName string) string {
  return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func format(s string) string {

  camel := regexp.MustCompile(`([a-z])([A-Z])`)
  spaces := regexp.MustCompile(`(\s+-\s+|\s+)`)
  quotes := regexp.MustCompile(`('|")`)
  parentheses := regexp.MustCompile(`(\(|\)|{|})`)
  accentsA := regexp.MustCompile(`(ã|á|à)`)
  accentsE := regexp.MustCompile(`(é|è)`)
  accentsI := regexp.MustCompile(`(í|ì)`)
  accentsO := regexp.MustCompile(`(ó|ò)`)
  accentsU := regexp.MustCompile(`(ú|ù)`)
  accentsN := regexp.MustCompile(`(ñ)`)
  accentsC := regexp.MustCompile(`(ç|ć)`)
  other := regexp.MustCompile(`(,|;|:|<|>|\?|!|@|#|\$|%|\^|\*|\+|=|~)`)
  multidashes := regexp.MustCompile(`(--+|__+)`)
  trailingDashes := regexp.MustCompile(`(^-|-$)`)

  s = strings.ToValidUTF8(s, "")
  s = strings.TrimSpace(s)
  s = quotes.ReplaceAllString(s, "")
  s = accentsA.ReplaceAllString(s, "a")
  s = accentsE.ReplaceAllString(s, "e")
  s = accentsI.ReplaceAllString(s, "i")
  s = accentsO.ReplaceAllString(s, "o")
  s = accentsU.ReplaceAllString(s, "u")
  s = accentsN.ReplaceAllString(s, "n")
  s = accentsC.ReplaceAllString(s, "c")
  s = parentheses.ReplaceAllString(s, "-")
  s = other.ReplaceAllString(s, "-")
  s = trailingDashes.ReplaceAllString(s, "")
  s = camel.ReplaceAllString(s, "$1-$2")
  s = spaces.ReplaceAllString(s, "-")
  s = strings.ReplaceAll(s, "_", "-")
  s = strings.ReplaceAll(s, "-.", ".")
  s = strings.ReplaceAll(s, ".-", "-")
  s = strings.ReplaceAll(s, "..", ".")
  s = multidashes.ReplaceAllString(s, "-")
  s = strings.ToLower(s)

  return s

}
