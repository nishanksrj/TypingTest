package main

import (
  "log"
  "net/http"
  "encoding/json"

  "github.com/gorilla/websocket"

  badger "github.com/dgraph-io/badger"

)


// connected clients
var clients = make(map[*websocket.Conn]bool)

// connected dashboards
var dashboards = make(map[*websocket.Conn]bool)


// broadcast channel
var broadcast = make(chan Message)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// Database
var db *badger.DB


// Message object
type Message struct {
    Roll string `json:"roll"`
    Message string `json:"message"`
    Speed int `json:"speed"`
    Words int `json:"words"`
    Chars int `json:"chars"`
}




func main(){

  var err error
  db, err = badger.Open(badger.DefaultOptions("/tmp/badger22"))
  if err != nil {
	  log.Fatal(err)
  }
  defer db.Close()


  fs := http.FileServer(http.Dir("./public"))
  http.Handle("/", fs)

  // handle socket connection request
  http.HandleFunc("/ws", handleConnections)

  // handle socket connection for dashbaord
  http.HandleFunc("/dash", handleDashConnections)

  // Start listening for incoming chat messages
  go handleMessages()

  // Start the server on localhost port 8000 and log and errors
  log.Println("Hola! http server started on :8000")
  err = http.ListenAndServe(":8000", nil)
  if err != nil{
    log.Fatal("ListenAndServe: ", err)
  }

}


func handleConnections(w http.ResponseWriter, r *http.Request){
  // Upgrade initial GET request to a websocket
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil{
    log.Fatal(err)
  }

  defer ws.Close()

  // add the new client to clients list(map)
  clients[ws] = true

  // listen for data
  for {
    var msg Message

    err := ws.ReadJSON(&msg)

    if err!= nil{
      log.Printf("Error: %v", err)
      delete(clients, ws)
      break
    }


    // update the database with key being the roll number and value being the
    // message recieved
    err = db.Update(func(txn *badger.Txn) error {

      // convert the data into bytes
      byteMsg,err := json.Marshal(msg)

      if err!=nil{
        log.Printf("%s",err)
      }

      e := badger.NewEntry([]byte(msg.Roll), byteMsg)
      err = txn.SetEntry(e)
      return err
    })

    if err!=nil{
      log.Printf("Error: %v", err)
    }

    // broadcast the message into channel
    broadcast <- msg
  }
}

func handleDashConnections(w http.ResponseWriter, r *http.Request){
  // Upgrade initial GET request to a websocket
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil{
    log.Fatal(err)
  }

  defer ws.Close()

  // add the new dashboard client into dashboards list(map)
  dashboards[ws] = true

  // extract the previous data by iterating over all Key-value pairs
  err = db.View(func(txn *badger.Txn) error {

    opts := badger.DefaultIteratorOptions
    opts.PrefetchSize = 20    // maximum limit of 20 studets to be displayed 

    it := txn.NewIterator(opts)
    defer it.Close()

    // iterate over all key value pairs
    for it.Rewind(); it.Valid(); it.Next() {
      item := it.Item()
      err := item.Value(func(v []byte) error {

        // unmarshal the data and send it to the dashboard client
        var msg Message
        err := json.Unmarshal(v, &msg)

        if err!= nil{
          log.Printf("Error: %v", err)
        }

        err = ws.WriteJSON(msg)

        if err != nil{
          log.Printf("Error: %v", err)
          ws.Close()
          delete(dashboards, ws)
        }
        return nil
      })
      if err != nil {
        return err
      }
    }
    return nil
  })

  // listen for incoming data (ping request usually for this websocket)
  for {
    var msg Message
    err := ws.ReadJSON(&msg)

    if err!=nil{
      log.Printf("Error: %v", err)
      delete(dashboards, ws)
      break
    }
  }

}


func handleMessages(){

  // send the data stored in broadcast to all the dashboard client
  for{
    msg := <-broadcast
    for dash := range dashboards{
      err := dash.WriteJSON(msg)
      if err != nil{
        log.Printf("Error: %v", err)
        dash.Close()
        delete(dashboards, dash)
      }
    }
  }
}
