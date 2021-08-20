# Learn gracefull shutdown

Gracefull shutdown akan menunggu service memproses hingga waktu jeda yang ditentukan.
Misalnya rata rata proses dari service yang berjalan adalah 20 detik. 
Maka kita bisa membuat timeout dengan bantuan `context.WithTimeout(serverCtx, 30*time.Second)`
Jadi jika masih ada proses akan ditunggu hingga 30 detik dan setelah itu dimatikan.

## Logika

Define context `serverCtx, serverStopCtx := context.WithCancel(context.Background())`
Lalu define channel `sig` yang harus didefine sebelum `signal.Notify`
`sig := make(chan os.Signal, 1)`

Sistem akan membaca sinyal interupsi dengan bantuan func `signal.Notify`. 
signal termination yang biasanya didengar:

- [SIGHUP](https://en.wikipedia.org/wiki/SIGHUP)
- [SIGINT](https://dsa.cs.tsinghua.edu.cn/oj/static/unix_signal.html#:~:text=The%20SIGINT%20signal%20is%20sent,break%22%20key%20can%20be%20used.&text=The%20SIGKILL%20signal%20is%20sent,to%20terminate%20immediately%20(kill).)   = Terminate dengan ctrl + c
- [SIGTERM](https://en.wikipedia.org/wiki/SIGTERM)
- [SIGQUIT](https://en.wikipedia.org/wiki/SIGQUIT)

Bisa baca disini https://dsa.cs.tsinghua.edu.cn/oj/static/unix_signal.html.


jika ada interupsi ke signal, maka sig akan masuk kedalam goroutine.
set timeout 30 detik sebelum shutdown. Ketika shutdown, maka tidak akan ada lagi request yang masuk.
Hanya akan menunggu process selesai selama 30 detik. Lebih dari itu akan langsung dimatikan.


```
    go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
```


## Cases

### Gracefull Shutdown Using Cobra CLI

Terkadang dengan memanfaatkan cobra cli kita design code kita seperti ini.

```bash
- cmd
	- root.go
	- rest.go
- handlers
	- server.go
- main.go
```

Kita tidak meletakan server di file main.go seperti sebelumnya.