package node

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"net"
	"ruthenium/src/chain"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
)

type Neighbor struct {
	ip         string
	port       uint16
	readWriter *bufio.ReadWriter
	mutex      sync.Mutex
}

func NewNeighbor(ip string, port uint16) *Neighbor {
	neighbor := new(Neighbor)
	neighbor.ip = ip
	neighbor.port = port
	return neighbor
}

func (neighbor *Neighbor) createReadWriter(targetPort uint16) *bufio.ReadWriter {
	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(neighbor.port)
	if err != nil {
		log.Fatal(err)
	}

	//ha.SetStreamHandler("/p2p/1.0.0", handleStream)

	// The following code extracts target's peer ID from the
	// given multiaddress
	target, err := makeBasicHost(targetPort)
	if err != nil {
		log.Fatal(err)
	}
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", target.ID().Pretty()))
	addr := target.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	ipfsaddr, err := multiaddr.NewMultiaddr(fullAddr.String())
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, time.Second)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		log.Fatalln(err)
	}

	// Create a buffered stream so that read and writes are non blocking.
	return bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(port uint16) (host.Host, error) {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	//if randseed == 0 {
	r = rand.Reader
	//} else {
	//	r = mrand.New(mrand.NewSource(randseed))
	//}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		libp2p.Identity(priv),
	}
	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)

	return basicHost, nil
}

func (neighbor *Neighbor) IpAndPort() string {
	return fmt.Sprintf("%s:%d", neighbor.ip, neighbor.port)
}

func (neighbor *Neighbor) ReadBlocks() []*chain.Block {
	var blocks []*chain.Block
	for {
		str, err := neighbor.readWriter.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" || str == "\n" {
			return []*chain.Block{}
		}

		if err := json.Unmarshal([]byte(str), &blocks); err != nil {
			log.Fatal(err)
		}
	}
}

func (neighbor *Neighbor) WriteBlocks(blocks []*chain.Block) {
	for {
		bytes, err := json.Marshal(blocks)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(blocks)
		neighbor.mutex.Lock()
		i, err := neighbor.readWriter.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		if err != nil || i == 0 {
			log.Println("ERROR: Failed to write blocks")
		}
		flushError := neighbor.readWriter.Flush()
		if flushError != nil {
			log.Println("ERROR: Failed to flush read writer")
		}
		neighbor.mutex.Unlock()
	}
}

func (neighbor *Neighbor) isFound() bool {
	target := fmt.Sprintf("%s:%d", neighbor.ip, neighbor.port)

	_, err := net.DialTimeout("tcp", target, time.Millisecond)
	if err != nil {
		fmt.Printf("%s not found, err:%v\n", target, err)
		return false
	}
	return true
}

func (neighbor *Neighbor) SetReadWriter(targetPort uint16) {
	neighbor.readWriter = neighbor.createReadWriter(targetPort)
}
