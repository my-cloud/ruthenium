package libp2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p/core/network"
	"io"
	"log"
	mrand "math/rand"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

type listen struct {
	protocol string
	ipVersion string
	ipAddress string
	port int
}

type Connection struct {
	stream 			network.Stream
	buffer 			*bufio.ReadWriter
	c  chan string
}
type Host struct {
	ctx        		context.Context
	cancel     		context.CancelFunc
	ha         		host.Host
	listening       listen
	insecure   		bool
	randomSeed 		int64
	source 			string
    connection 		Connection
}

// StartHost is temporary and meant to test breaking changes from a shortcut calling
func StartHost() {
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	targetF := flag.String("d", "", "target peer to dial")
	insecureF := flag.Bool("insecure", false, "use an unencrypted connection")
	seedF := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	h, err := NewHost(*listenF, *insecureF, *seedF)
	if err != nil {
	  log.Fatal(err)
	}
	defer h.cancel()
	h.connection.c = make(chan string)
	if *targetF == "" {
      h.source = "Server"
      //h.NewServer()

      	  str := <- h.connection.c
		  log.Printf("%v",str )
		  <-h.ctx.Done()


	} else {
		h.source = "Sender"
		h.NewClient( *targetF)
	}
	log.Println("end of starthost func")
}


func NewHost(port int, insecure bool, randomSeed int64) (Host, error) {
	var h Host
	var err error
	h.listening= listen{"tcp","ip4","0.0.0.0",port}
    h.insecure = insecure
    h.randomSeed = randomSeed
	h.ctx, h.cancel = context.WithCancel(context.Background())

	// LibP2P code uses golog. They log with different string IDs (i.e. "swarm")
	golog.SetAllLoggers(golog.LevelInfo) // Change to INFO for extra info

	// Make a host that listens on the given multiaddress
	h.ha, err = h.makeBasicHost()

	fullAddr := h.getHostAddress()
	log.Printf("%sAddress %s\n", h.source, fullAddr)

	// Set a stream handler on host A. /ruthenium/1.0.0 is
	// a user-defined protocol name.
	h.ha.SetStreamHandler("/ruthenium/1.0.0", h.handleStream)
	//if h.source == "Server"{
	//// Run until canceled.
	//<-h.ctx.Done()
	//}
	return h, err
}

func NewServer(h Host){
	defer h.cancel()
		<-h.ctx.Done()
}

func (h Host) NewClient(target string){
	h.runSender(target)
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It won't encrypt the connection if insecure is true.
func (h Host) makeBasicHost() (host.Host, error) {
	var r io.Reader
	if h.randomSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(h.randomSeed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/%s/%s/%s/%d",
			h.listening.ipVersion,
			h.listening.ipAddress,
			h.listening.protocol,
			h.listening.port,
		)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	if h.insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	return libp2p.New(opts...)
}

func (h Host) getHostAddress() string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", h.ha.ID().String()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := h.ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

