package operativac2

type Props struct {
	actorFunc   func() Actor
	mailboxSize int
	naziv       string
}

func NewProps(naziv string, actorFunc func() Actor) *Props {
	return &Props{naziv: naziv, actorFunc: actorFunc, mailboxSize: 10}
}

func (p *Props) WithMailboxSize(size int) *Props {
	p.mailboxSize = size
	return p
}
