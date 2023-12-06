package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/bus"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/domain/lru"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/repo"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/transport"
	"github.com/adrian-lin-1-0-0/scaling-websocket-demo/usecases"
	"golang.org/x/net/websocket"
)

func GetLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}

func main() {
	ips, err := GetLocalIPs()
	if err != nil {
		panic("Cant get local ip")
	}

	host := ips[0].String()
	messageServicePrefix := "/message"
	serviceName := "/chat/"
	cacheSize := 1 << 20

	port := os.Getenv("SERVER_PORT")
	etcd := os.Getenv("ETCD_ENDPOINT") //localhost:2379

	if port == "" {
		port = ":8000"
	} else {
		port = ":" + port
	}

	if etcd == "" {
		etcd = "localhost:2379"
	}
	self := host + port

	etcdEndpoints := []string{etcd}

	messageBus := bus.NewMessageBus(1)

	ctx := context.Background()
	discovery, err := repo.NewDiscovery(ctx, etcdEndpoints, self)
	if err != nil {
		panic(err)
	}

	cache := lru.New(cacheSize)

	lookupRepo := repo.NewLookupRepo(serviceName, domain.Endpoint(self), discovery)

	userRepo := repo.NewUserRepo(lookupRepo, cache)
	lookupPeer := transport.NewLookupPeer(userRepo)
	messagePeer := transport.NewMessagePeer(messageServicePrefix)

	lookupUsecase := usecases.NewLookupUsecase(userRepo)
	messageUsecase := usecases.NewMessageUsecase(messageBus)
	userUsecase := usecases.NewUserUsecase(
		domain.Endpoint(self),
		userRepo,
		messagePeer,
		messageBus,
		lookupPeer,
	)

	chatService := transport.NewChatService(userUsecase)
	lookupService := transport.NewLookupService(lookupUsecase, cache)
	messageService := transport.NewMessageService(messageUsecase)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocket.Handler(chatService.Handler))

	http.HandleFunc("/lookup", lookupService.LookupHandler)
	http.HandleFunc("/register", lookupService.RegisterHandler)

	http.HandleFunc(messageServicePrefix, messageService.Handler)

	log.Default().Println("Server started at ", self)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
