# Pbq 
In-Memory message queuing over HTTP API using a concurrent, thread-safe pubsub architecture

# Installation
```sh
> go get -v github.com/Manan007224/Pbq
> cd $GOPATH/src/github.com/Manan007224/Pbq
> go build -o Pbq
> ./Pbq
```
Pbq by default runs on `localhost:3000`

## Examples - Client 
**Subscriber**

```go 

c := client.New()
tick := time.NewTicker(time.Second)
for _ = range tick.C {
  if err := c.Publish("foo", []byte(`bar`)); err != nil {
    log.Println(err)
    break
  }
}
```


**Publisher**

```go

c := client.New()
ch, err := c.Subscribe("foo")
if err != nil {
  log.Println(err)
  return
}
defer c.Unsubscribe(ch)

for {
  select {
  case e := <-ch:
    log.Println(string(e))
    fmt.Println(string(e))
  }
}
```  
