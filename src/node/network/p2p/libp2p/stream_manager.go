package libp2p

import (
	"bufio"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
)


// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}

func (h Host) handleStreaming (s network.Stream) {
    h.debugConnection(s)
	//log.Printf("%s received new stream", h.source)
	if err := doEcho(s); err != nil {
		log.Println(err)
		s.Reset()
	} else {
		s.Close()
	}
}

func (h Host) handleStream(s network.Stream) {
	h.debugConnection(s)
    h.connection.stream = s
	// Create a buffer stream for non blocking read and write.
	h.connection.buffer = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go h.readData()
	//go h.writeData()

	// stream 's' will stay open until you close it (or the other side closes it).
}

func (h Host) readData() error {
		str, err := h.connection.buffer.ReadString('\n')
	    log.Println("readData")
		if err != nil {
			return err
		}

			log.Printf("\x1b[32m%s\x1b[0m ", str)
	_, err = h.connection.buffer.Write([]byte("{\"Amount\":\"900\"}\n"))

	log.Println("writeData")
	h.connection.buffer.Flush()
	if err != nil {
		log.Println(err)
		h.connection.stream.Reset()
	} else {
		h.connection.stream.Close()
	}
	h.connection.c <- str
	return nil
}

func (h Host) writeData() {
	//stdReader := bufio.NewReader(os.Stdin)
	//
	//for {
	//	log.Print("> ")
	//	sendData, err := stdReader.ReadString('\n')
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	rw.WriteString(fmt.Sprintf("%s\n", sendData))
	//	rw.Flush()
	//}
	//rw.Write([]byte("{\"Amount\":\"800\"}\n"))
	//rw.Flush()
	_, err := h.connection.buffer.Write([]byte("{\"Amount\":\"900\"}\n"))
	log.Println("writeData")
	h.connection.buffer.Flush()
	if err != nil {
		log.Println(err)
		h.connection.stream.Reset()
	} else {
		h.connection.stream.Close()
	}
}

func (h Host) debugConnection (s network.Stream)  {
	log.Printf("new stream ID %s on protocol %s", s.ID(),s.Protocol())
	log.Printf("hosts %v",h.ha.Addrs())
	log.Printf("Peers %v", h.ha.Network().Peers())
	for _, peerId := range h.ha.Network().Peers() {
		log.Printf("peer address %v", h.ha.Peerstore().Addrs(peerId))
	}
}