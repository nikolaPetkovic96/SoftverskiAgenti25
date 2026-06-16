// package main

// import (
// 	op "at24/operativac2"
// 	"fmt"
// 	"strconv"
// 	"time"
// )

// type PrinterActor struct{}

// func (p *PrinterActor) Receive(envelope op.Envelope) {
// 	switch msg := envelope.Message.(type) {
// 	case string:
// 		fmt.Println("PrinterActor received:", msg)
// 	default:
// 		fmt.Println("PrinterActor received unknown message")
// 	}
// }

// type CreateChildren struct {
// 	brojDece int
// 	op       *op.Operativac
// }
// type WorkerActor struct {
// 	brojac int
// }
// type ChildActor struct {
// 	brojac int
// }

// // Receive implements operativac2.Actor.
// func (c *ChildActor) Receive(envelope op.Envelope) {
// 	fmt.Printf("child primio poruku pod brojem: %d", c.brojac)
// 	c.brojac++
// }

// func (w *WorkerActor) Receive(envelope op.Envelope) {
// 	w.brojac++
// 	switch msg := envelope.Message.(type) {
// 	case string:
// 		fmt.Println("WorkerActor processed:", msg+"!")
// 	case *CreateChildren:
// 		fmt.Printf("Kreiranje dece")
// 		chProps := op.NewProps("work", func() op.Actor {
// 			return &ChildActor{brojac: 5}
// 		})

// 		for i := 0; i < msg.brojDece; i++ {

// 			msg.op.SpawnChild(chProps, strconv.Itoa(i))
// 		}
// 		time.Sleep(1 * time.Second)
// 		msg.op.SendToChildren(struct{}{})
// 	default:
// 		fmt.Printf("WorkerActor received unknown message. BROJAC: %d\n", w.brojac)
// 	}
// }

// func main() {
// 	printerProps := op.NewProps("print", func() op.Actor {
// 		return &PrinterActor{}
// 	})
// 	workerProps := op.NewProps("work", func() op.Actor {
// 		return &WorkerActor{brojac: 0}
// 	})

// 	sl := op.NovaSluzba()
// 	sl.Spawn(printerProps, "p")
// 	w := sl.Spawn(workerProps, "w")

// 	sl.Send("p", *op.NewEnvelope("HELLO P", "unknown", "p", ""))
// 	sl.Send("w", *op.NewEnvelope("HELLO W", "unknown", "W", ""))

// 	sl.Stop("p")

// 	sl.Send("p", *op.NewEnvelope("35", "unknown", "p", ""))
// 	sl.Send("w", *op.NewEnvelope(23, "unknown", "2", ""))

// 	sl.Send("w", *op.NewEnvelope(&CreateChildren{brojDece: 5, op: w}, "unknown", "p", ""))
// 	time.Sleep(1 * time.Second)

// 	sl.UgasiSluzbu()
// 	time.Sleep(1 * time.Second)
// }
