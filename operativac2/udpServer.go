package operativac2

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type udpServer struct {
	Address *net.UDPAddr
	sluzba  *Sluzba
	nastavi bool
	conn    *net.UDPConn
}

func NoviServer(port string, s *Sluzba) *udpServer {
	connIP, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("err")
	}
	defer connIP.Close()
	address := connIP.LocalAddr().(*net.UDPAddr)
	add := strings.Split(address.String(), ":")
	novaA := add[0]
	novaA += ":"
	novaA += port
	udpAddr, err := net.ResolveUDPAddr("udp", novaA)
	if err != nil {
		fmt.Printf("Error resolving address: %v\n", err)
		return nil
	}
	return &udpServer{
		Address: udpAddr,
		sluzba:  s,
		nastavi: true,
	}
}
func (u *udpServer) Pokreni() {

	// Create UDP listener
	conn, err := net.ListenUDP("udp", u.Address)
	if err != nil {
		fmt.Printf("Error listening: %v\n", err)
		return
	}
	fmt.Println("Adresa udp servera : ", u.Address)
	u.conn = conn
	//u.sluzba.conn = conn
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading UDP message: %v\n", err)
			continue
		}
		msg := buffer[:n]
		env, err := Deserijalizuj(msg)
		if err != nil {
			//greska pri prijemu, obracdjeno u deserijalizuj metodi
		} else {
			fmt.Printf("Received message from %s: %s\n", clientAddr, msg)
			u.sluzba.DodajSluzbu(env.SenderIp)
			u.sluzba.Send(env.ReceiverId, env)
		}
		if !u.nastavi {
			break
		} else {
			continue
		}
	}
}

func (u *udpServer) Zaustavi() { u.nastavi = false }

func (u *udpServer) SendRemote(env Envelope, adr *net.UDPAddr) {
	u.conn.WriteToUDP(Serijalizuj(env), adr)
}

func Serijalizuj(env Envelope) []byte {
	fmt.Printf("PRE SERIJALIZACIJE : %s, %s, %s, %s, %d\n", env.Message, env.SenderIp, env.SenderId, env.ReceiverId, env.CntFail)
	serializedData, err := json.Marshal(env)
	if err != nil {
		fmt.Printf("Greska pri serijalizaciji !\n")
		return nil
	}
	fmt.Println("Serijalizovani obj: ", string(serializedData))
	return serializedData
}

func Deserijalizuj(serializedDate []byte) (Envelope, error) {
	var env Envelope
	err := json.Unmarshal(serializedDate, &env)
	if err != nil {
		fmt.Printf("Greska pri serijalizaciji! ")
	} else {
		if v, ok := env.Message.(interface{}); ok {
			fmt.Println("Primljena Poruka:", v)

		}
	}
	return env, err
}
