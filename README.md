# pbq 
In-Memory message broker over HTTP API using a concurrent, thread-safe pubsub architecture.

# Installation
```sh
> go get -v github.com/Manan007224/pbq
> cd $GOPATH/src/github.com/Manan007224/pbq
> go build -o pbq
> ./pbq
```
Pbq by default runs on `localhost:3000`

## Usage

Run server on `localhost:3000`
```shell
./pqb
```
### Run Client

**Subscribe**
```shell
./pbq --client --topic=foo --subscribe 
```

**Publish**
```
echo "A completely arbitrary message" | ./pbq --client --topic=foo --publish
```

## Examples - Client (Library)
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


## Author

Manan Maniyar:[ E-mail](mailto:maniyarmanan1996@gmail.com), [@Manan007224](https://www.github.com/Manan007224)
