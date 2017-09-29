package main

// set HOSTNAME, PREFIX, VAULT_ADDR and VAULT_TOKEN environment variable and others as needed.
// curl -XPOST http://localhost:8080/add -d "message=This is my secet message to you."

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "crypto/rc4"
  "crypto/rand"
  "crypto/sha256"
  "encoding/base32"

  "github.com/gorilla/mux"
  vault "github.com/hashicorp/vault/api"
)

func main() {

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/health", Health)
  router.HandleFunc("/get/{key}", GetSecret).Methods("GET")
  router.HandleFunc("/{add:add\\/?}", PutSecret).Methods("POST")

  log.Fatal(http.ListenAndServe(":8080", router))
}

func Health(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "live")
}

func PutSecret(w http.ResponseWriter, r *http.Request) {
  // Generate a random key,
  key := make([]byte, 20)
  rand.Read(key)
  // make it URL friendly
  sKey := base32.StdEncoding.EncodeToString(key)
  //url := "https://" + os.Getenv("HOSTNAME") + "/get/" + sKey
  url := os.Getenv("HOSTNAME") + "/get/" + sKey
  fmt.Fprintln(w, "{\"url\": \"" + url + "\"}")
  msg := r.FormValue("message")

  cipher := EncOrDec(key,msg)
  encodedMsg := base32.StdEncoding.EncodeToString([]byte(cipher))
  WriteToVault(key, encodedMsg)
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sKey := vars["key"]
  key, _ := base32.StdEncoding.DecodeString(sKey)
  encodedMgs := ReadFromVault(key)
  cipher, _ := base32.StdEncoding.DecodeString(encodedMgs)
  msg := EncOrDec(key,string(cipher))
  fmt.Fprintln(w, "{\"secret\": \"" + msg + "\"}")
}

func EncOrDec(key []byte, sMsg string) string {
    msg := []byte(sMsg)
    c, _ := rc4.NewCipher(key)
    c.XORKeyStream(msg, msg)
    return string(msg)
}

func WriteToVault(key []byte, secret string){
  // get the hash of the key
  sha_256 := sha256.New()
  sha_256.Write(key)
  keyHash := base32.StdEncoding.EncodeToString(sha_256.Sum(nil))
  path := os.Getenv("PREFIX") + "/onetime/" + keyHash
  client, _ := vault.NewClient(vault.DefaultConfig())
  LogicalClient := client.Logical()
  _, err := LogicalClient.Write(path,
    map[string]interface{}{
      "secret":  secret,
    })
  if (err != nil){
    fmt.Println(err)
  }
}

func ReadFromVault(key []byte) string {
  // get the hash of the key
  sha_256 := sha256.New()
  sha_256.Write(key)
  keyHash := base32.StdEncoding.EncodeToString(sha_256.Sum(nil))
  path := os.Getenv("PREFIX") + "/onetime/" + keyHash
  client, _ := vault.NewClient(vault.DefaultConfig())
  LogicalClient := client.Logical()
  secret, _ := LogicalClient.Read(path)
  if secret != nil {
    for _, secret := range secret.Data {
      LogicalClient.Delete(path)
      return secret.(string)
    }
  }
  return ""
}
