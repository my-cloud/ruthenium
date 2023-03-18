package libp2p

import (
	"context"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"time"
)

func (h Host) runSender(targetPeer string) {
	// Turn the targetPeer into a multiaddr.
	maddr, err := ma.NewMultiaddr(targetPeer)
	if err != nil {
		log.Println(err)
		return
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return
	}

	// Adding the peer ID and a targetAddr to the peerstore so LibP2P knows how to contact it
	h.ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

for {
	log.Println("sender opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /ruthenium/1.0.0 protocol
	s, err := h.ha.NewStream(context.Background(), info.ID, "/ruthenium/1.0.0")
	if err != nil {
		log.Println(err)
		return
	}
	// Sending json in byte format
	_, err = s.Write([]byte("{\"Address\":\"0x9C69000000000000000000000000000000CB\"}\n"))
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(1 * time.Second)

	out, err := io.ReadAll(s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("read reply: %q\n", out)

	time.Sleep(1 * time.Second)
}
}