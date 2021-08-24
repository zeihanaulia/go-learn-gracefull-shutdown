# Learn gracefull shutdown

Kenapa harus menegetahui praktik ini. Karena masuk kedalam [12 factor app](https://12factor.net/disposability) pada bagian disposable

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

## Build docker

```
docker build -t glgs-http-server  .
docker tag glgs-http-server zeihanaulia/glgs-http-server:1.0.0
```

## Run using kubernetes.

```
 kubectl apply -f http-server.yml

 -- RESTART
 kubectl rollout restart deployments/glgs-http-server

 -- Check logs
 kubectl logs -f deployments/glgs-http-server

 -- delete service
kubectl delete deployment glgs-http-server

-- check
kubectl get pods
kubectl get service
kubectl get deployments
```

## Referensi

- https://learnk8s.io/graceful-shutdown