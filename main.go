package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "os/exec"

  "github.com/pbnjay/gosns"
)

var args []string
var binary string

// RunCommand runs the provided command-line argument
func RunCommand(msg *gosns.Message) {
  if msg == nil {
    log.Println("Topic Subscription Confirmed.")
    return
  }

  log.Println("----")
  log.Printf("Running command: '%s...'\n", binary)
  out, err := exec.Command(binary, args...).CombinedOutput()
  if err != nil {
    log.Println(err)
  }
  log.Printf("%s", out)
  log.Println("----")
}

func checkArgs(command []string, topicArn string) {
  if topicArn == "REQUIRED" {
    log.Fatal("No SNS topic ARN specified")
  }

  var lookErr error
  binary, lookErr = exec.LookPath(command[0])
  if lookErr != nil {
    log.Fatal(lookErr)
  }
  args = command[1:len(command)]
}

func main() {
  var endpoint = flag.String("e", "/", "web endpoint")
  var port = flag.Int("p", 8080, "port on which to listen")
  var topicArn = flag.String("a", "REQUIRED", "SNS topic ARN")
  flag.Parse()
  command := flag.Args()
  checkArgs(command, *topicArn)

  snsServer := &gosns.Server{}
  snsServer.Logger = log.New(os.Stderr, "GOSNS ", log.LstdFlags)
  snsServer.AddTopic(*topicArn, *endpoint, RunCommand)
  log.Fatal(snsServer.ListenAndServe(fmt.Sprintf(":%d", *port)))
}
