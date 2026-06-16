package main

import (
	op "at24/operativac2"
	"fmt"
	"time"
)

type PrinterActor struct {
}

func (p *PrinterActor) Receive(ctx op.Context) {
	switch msg := ctx.Env.Message.(type) {

	case string:
		fmt.Println("PrinterActor received:", msg)
		time.Sleep(1 * time.Second)
		ctx.Pok.Sluzba.Conn.SendRemote(op.Envelope{Message: "Zahveljujem na pozdravu!", SenderId: ctx.Pok.GetId(), SenderIp: ctx.Pok.Sluzba.Conn.Address.String(), ReceiverId: ctx.Env.SenderId},
			ctx.Pok.Sluzba.PoznateSluzbe[ctx.Env.SenderIp])
	default:
		fmt.Println("PrinterActor received unknown message")
	}
}

func main() {
	printerProps := op.NewProps("print", func() op.Actor {
		return &PrinterActor{}
	})

	sl := op.NovaSluzba()
	sl.Spawn(printerProps, "p")

	sl.Send("p", *op.NewEnvelope("HELLO P", "unknown", "p", ""))
	//sl.Send("p", *op.NewEnvelope(p, "root", "p", sl.Conn.Address.String()))
	server := op.NoviServer("9090", sl)
	sl.Conn = server
	go server.Pokreni()

	sl2 := op.NovaSluzba()
	sl2.Spawn(printerProps, "p2")

	server2 := op.NoviServer("9091", sl2)
	sl2.Conn = server2
	go server2.Pokreni()
	var addr1 []string
	addr1 = append(addr1, server.Address.String())
	sl2.DodajPoznateSluzbu(addr1)
	time.Sleep(2 * time.Second)

	sl2.PosaljiDrugojSluzbi(server.Address.String(),
		*op.NewEnvelope("pozdrav sa server2 !", "p2", "p", server2.Address.String()),
	)
	time.Sleep(2 * time.Second)

	//sl.Stop("p")

	//sl.Send("p", *op.NewEnvelope("35", "unknown", "p"))
	//sl.Send("w", *op.NewEnvelope(23, "unknown", "2"))

	//sl.Send("w", *op.NewEnvelope(&CreateChildren{brojDece: 5, op: w}, "unknown", "p"))
	//time.Sleep(1 * time.Second)

	sl.UgasiSluzbu()
	time.Sleep(1 * time.Second)
}
