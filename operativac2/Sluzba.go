package operativac2

import (
	"fmt"
	"net"
	"sync"
)

type Sluzba struct {
	Operativci    map[string]*Operativac
	Kinds         map[string]*Props
	mu            sync.Mutex
	wg            sync.WaitGroup
	Conn          *udpServer
	PoznateSluzbe map[string]*net.UDPAddr
}

func NovaSluzba() *Sluzba {
	return &Sluzba{
		Operativci:    make(map[string]*Operativac),
		PoznateSluzbe: make(map[string]*net.UDPAddr),
		Conn:          nil,
		Kinds:         make(map[string]*Props),
	}
}

func (sl *Sluzba) Spawn(props *Props, id string) *Operativac {
	sl.mu.Lock()
	_, exProps := sl.Kinds[props.naziv]
	if !exProps {
		sl.Kinds[props.naziv] = props
	}
	sl.mu.Unlock()
	sl.mu.Lock()
	//defer sl.mu.Unlock()
	_, exists := sl.Operativci[id]
	if !exists {
		actor := props.actorFunc()
		op := &Operativac{
			mailbox:     make(chan Envelope, props.mailboxSize),
			actor:       actor,
			stopSignal:  make(chan struct{}),
			Sluzba:      sl,
			parent:      nil,
			obustavljen: false,
			penzionisan: false,
			info:        Info{nazivSluzbe: "", id: id, cntFail: 0, kind: props.naziv, parent: ""},
		}
		sl.Operativci[id] = op
		sl.mu.Unlock()
		//go sl.AktivirajOperativca(op)
		op.Start()
		return op
	} else {
		sl.mu.Unlock()
		fmt.Println("Posotji operativac sa zadatim id")
		return nil
	}
}

// func (sl *Sluzba) AktivirajOperativca(op *Operativac) {
// 	for {
// 		select {
// 		case msg := <-op.mailbox:
// 			op.actor.Receive(msg)
// 		case <-op.stopSignal:
// 			return
// 		}
// 	}
//}

func (sl *Sluzba) Send(id string, msg Envelope) {
	sl.mu.Lock()
	op, exists := sl.Operativci[id]
	sl.mu.Unlock()
	if exists { //TODO razradi ako je puno sanduce
		op.mailbox <- msg
	} else {
		fmt.Printf("Ne postoji operativac sa zadatim id = %s , msg:[%v]\n", id, msg.Message)
	}
}

func (sl *Sluzba) Stop(id string) {
	sl.mu.Lock()
	op, exists := sl.Operativci[id]
	sl.mu.Unlock()
	//if exists {
	//	fmt.Printf("exists, %s\n", id)
	//}

	if exists && !op.obustavljen {
		//fmt.Printf("Stop signal : %s\n", id)
		fmt.Printf(" Poslata poruka za gasenje, %s\n", op.info.id)
		op.stopSignal <- struct{}{} // Signal actor to stop
		//delete(sl.operativci, id)
	}
}

func (sl *Sluzba) Remove(id string) {
	sl.mu.Lock()
	op, exists := sl.Operativci[id]

	if exists && op.penzionisan {
		delete(sl.Operativci, id)
	}
	sl.mu.Unlock()
	//sl.wg.Done()
}

func (sl *Sluzba) UgasiSluzbu() {
	sl.mu.Lock()
	for _, op := range sl.Operativci {
		//sl.wg.Add(1)
		if op.info.parent == "" {
			go func() {
				//defer sl.wg.Done()
				sl.Stop(op.info.id)
			}()
		}
	}
	sl.mu.Unlock()
	sl.wg.Wait()
	//fmt.Printf("Sluzba < %s> je uspesno UGASENA\n", sl.naziv)
	fmt.Printf("Sluzba je uspesno UGASENA\n")
}

func (sl *Sluzba) DodajPoznateSluzbu(poznate []string) {
	for _, s := range poznate {
		addr, err := net.ResolveUDPAddr("udp", s)
		if err != nil {
			fmt.Printf("NEISPRAVAN FORMAT ADRESE : %s \n", s)
			continue
		}
		fmt.Printf("Dodavanje adrese :%s\n", s)
		sl.PoznateSluzbe[s] = addr
	}
}

func (sl *Sluzba) DodajSluzbu(poznate string) {

	addr, err := net.ResolveUDPAddr("udp", poznate)
	if err != nil {
		fmt.Printf("NEISPRAVAN FORMAT ADRESE : %s \n", poznate)
		return
	}
	fmt.Printf("Dodavanje adrese :%s\n", poznate)
	sl.mu.Lock()
	sl.PoznateSluzbe[poznate] = addr
	sl.mu.Unlock()

}

func (sl *Sluzba) PosaljiDrugojSluzbi(adr string, env Envelope) {
	sl.Conn.SendRemote(env, sl.PoznateSluzbe[adr])
}
