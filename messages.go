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
	op     *ebiten.DrawImageOptions
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

func (m *MessageQueue) DrawMessages(layer *ebiten.Image, x, y int) {

	if m.redraw || m.canvas == nil {

		canvas := ebiten.NewImage(300, 100)

		for i := 0; i < len(messages.queue); i++ {
			// TODO calculate longest string, justify each message and translate
			// 	the canvas accordingly (instead of using static coordinates)
			cy := 0 + (i * 14)
			msg := messages.queue[len(messages.queue)-1-i].ToString()
			bounds := text.BoundString(fontDetail, msg)
			textImg := ebiten.NewImage(bounds.Dx(), bounds.Dy()+16)
			text.Draw(textImg, msg, fontDetail, 0, bounds.Dy(), color.White)
			textImgOps := &ebiten.DrawImageOptions{}
			alpha := 1 - (0.3 * float64(i))
			textImgOps.ColorM.Scale(1, 1, 1, alpha)
			// TODO make it so we can draw bottom to top also via a bool or something
			// 	and/or justify right to left
			textImgOps.GeoM.Translate(0, float64(cy))
			canvas.DrawImage(textImg, textImgOps)
		}

		m.canvas = canvas
	}

	m.redraw = false
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(float64(x), float64(y))
	layer.DrawImage(m.canvas, ops)
}
