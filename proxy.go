package main

import(
"log"
"os"
"bufio"
)

func getProxyFile(filename string) ([]string, error){
  file, err := os.Open(filename)

  if err != nil{
    log.Println("failed to open file")
  }
  defer file.Close()

  var proxyList []string

  scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        proxyList = append(proxyList, line)
    }

     if scanner.Err() != nil {
        log.Println("Reading proxy list error:", scanner.Err())
        return nil, scanner.Err()
    }


    return proxyList, nil

}

