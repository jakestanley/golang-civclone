package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type Message struct {
	content string
	dupes   int
}

func (m *Message) ToString() string {
	if m.dupes > 0 {
		return fmt.Sprintf("%s (%d)", m.content, m.dupes+1)
	}
	return m.content
}

type MessageQueue struct {
	// TODO duplicate message indicator
	queue  []*Message
	max    int
	redraw bool
	canvas *ebiten.Image
}

func CreateMessages() MessageQueue {
	q := make([]*Message, 0)
	return MessageQueue{
		queue:  q,
		max:    8,
		redraw: true,
	}
}

func (m *MessageQueue) AddMessage(content string) {

	// get the last message. if the content matches, increment the value
	if len(m.queue) > 0 {
		lm := m.queue[len(m.queue)-1]
		if lm.content == content {
			lm.dupes++
			return
		}
	}

	queue := append(m.queue, &Message{
		content: content,
		dupes:   0,
	})
	if len(queue) > m.max {
		// dequeue
		queue = queue[1:]
	}
	m.queue = queue
	m.redraw = true
}

func (m *MessageQueue) DrawMessages(screen *ebiten.Image) {

	x := 300
	y := 0

	if m.redraw || m.canvas == nil {

		canvas := ebiten.NewImage(600, 100)

		for i := 0; i < len(messages.queue); i++ {
			// TODO calculate longest string, justify each message and translate
			// 	the canvas accordingly (instead of using static coordinates)
			text.Draw(canvas, messages.queue[len(messages.queue)-1-i].ToString(), fontDetail, 20, 20, color.White)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y+(i*14)))
			alpha := 1 - (0.3 * float64(i))
			op.ColorM.Scale(1, 1, 1, alpha)
			screen.DrawImage(canvas, op)
		}

		m.canvas = canvas
	}

	m.redraw = false
	screen.DrawImage(m.canvas, &ebiten.DrawImageOptions{})
}
