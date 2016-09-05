package main

import (
  "fmt"
  "log"
  "os"
  "io"
  "net/http"
  "database/sql"
  _"github.com/lib/pq"
  "encoding/json"
  "os/exec"
  "bufio"
  "strings"
)

func logging(message string)  {
  fmt.Printf(message)
}

func checkError(error error) {
  if error != nil {
    log.Fatal(error)
  }
}

func getURLsAndIds() ([]string, []int) {
  db, err := sql.Open("postgres", "user=" + os.Getenv("DATABASE_USERNAME") + " password=" + os.Getenv("DATABASE_PASSWORD") + " dbname=" + os.Getenv("DATABASE_NAME") + " sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  queryStmt, err := db.Prepare("SELECT id, instagram_object FROM portfolios WHERE instagram_id NOTNULL AND instagram_object NOTNULL") // AND portfolio_image_file_name ISNULL
  if err != nil {
    log.Fatal(err)
  }
  defer queryStmt.Close()

  var rowsArrayIds []int
  var rowsArray []string
  rows, err := queryStmt.Query()
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var id int
    var instagram_object []byte
    if err := rows.Scan(&id, &instagram_object); err != nil {
      log.Fatal(err)
    }
    var message map[string]interface{}
    if err := json.Unmarshal([]byte(instagram_object), &message); err != nil {
      log.Fatal(err)
    }
    code := message["code"].(string)
    images := message["images"].(map[string]interface{})
    standard_resolution := images["standard_resolution"].(map[string]interface{})
    standard_resolution_url := standard_resolution["url"].(string)
    rowsArray = append(rowsArray, standard_resolution_url)
    download_image(standard_resolution_url, code, id)
    updateDb(code, id)
  }
  return rowsArray, rowsArrayIds
}

func updateDb(code string, id int) {
  filename := code + ".jpg"
  db, err := sql.Open("postgres", "user=" + os.Getenv("DATABASE_USERNAME") + " password=" + os.Getenv("DATABASE_PASSWORD") + " dbname=" + os.Getenv("DATABASE_NAME") + " sslmode=disable")
  if err != nil { log.Fatal(err) }
  defer db.Close()

  queryStmt, err := db.Prepare("UPDATE portfolios SET portfolio_image_file_name = $1 WHERE id = $2")
  if err != nil {
    log.Fatal(err)
  }
  defer queryStmt.Close()

  update, err := queryStmt.Query(filename, id)
  if err != nil {
    log.Fatal(err)
  }
  defer update.Close()
}

func createDirectory(code string, size string) (string) {
  directory := os.Getenv("FILE_PATH") + code + "." + "jpg/" + size + "/"
  directory_err := os.MkdirAll(directory, 0755)
  if directory_err != nil {
    log.Fatal(directory_err)
  }
  return directory
}

func download_image(url string, code string, id int) {
  response, image_err := http.Get(url)
  if image_err != nil {
    log.Fatal(image_err)
  }

  defer response.Body.Close()

  original_directory := createDirectory(code, "original")

  download_string := original_directory + code + ".jpg"
  file, download_string_err := os.Create(download_string)
  if download_string_err != nil {
    log.Fatal(download_string_err)
  }

  _, err := io.Copy(file, response.Body)
  if err != nil {
    log.Fatal(err)
  }
  file.Close()
  updateDb(code, id)

  large_directory := createDirectory(code, "large")
  large_directory_string := large_directory + code + ".jpg"
  
  cmd := exec.Command(os.Getenv("IMAGEMAGICK_PATH"), download_string, "-resize", "1000x1000", "-gravity", "center", "-extent", "1000x1000", large_directory_string)
  cmdErr := cmd.Run()
  if cmdErr != nil {
    log.Fatal(cmdErr)
  }

  medium_directory := createDirectory(code, "medium")
  medium_directory_string := medium_directory + code + ".jpg"
  
  cmd = exec.Command(os.Getenv("IMAGEMAGICK_PATH"), download_string, "-resize", "600x600", "-gravity", "center", "-extent", "600x600", medium_directory_string)
  cmdErr = cmd.Run()
  if cmdErr != nil {
    log.Fatal(cmdErr)
  }

  thumb_directory := createDirectory(code, "thumb")
  thumb_directory_string := thumb_directory + code + ".jpg"
  
  cmd = exec.Command(os.Getenv("IMAGEMAGICK_PATH"), download_string, "-resize", "200x200", "-gravity", "center", "-extent", "200x200", thumb_directory_string)
  cmdErr = cmd.Run()
  if cmdErr != nil {
    log.Fatal(cmdErr)
  }
  
  
}

func func_name() {
  
}

func main() {
  envFile, fileErr := os.Open(os.Getenv("ENV_VARS_FILE"))
  checkError(fileErr)

  fileScanner := bufio.NewScanner(envFile)
  fileLine := 1
  fileLineText := ""
  for fileScanner.Scan() {
    fileLineText = fileScanner.Text()

    pair := strings.Split(fileLineText, "=")
    os.Setenv(pair[0], pair[1])
    fileLine++
  }

  if scannerErr := fileScanner.Err(); scannerErr != nil {
      checkError(scannerErr)
  }
  getURLsAndIds()
}
