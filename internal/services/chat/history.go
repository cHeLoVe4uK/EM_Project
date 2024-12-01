package chat

import (
	"context"
	"sync"
	"sync/atomic"
)

type History struct {
	msgs [100]MessageDTO

	msgCount uint8
	current  uint8

	newMsgs   atomic.Uint32
	stashChan chan struct{}

	mu sync.RWMutex
}

func (r *Room) newHistory(ctx context.Context, chatID string) (*History, error) {

	msgs, err := r.msgService.GetChatMessages(ctx, chatID)
	if err != nil {
		return nil, err
	}

	h := &History{
		msgs:    [100]MessageDTO{},
		mu:      sync.RWMutex{},
		newMsgs: atomic.Uint32{},

		stashChan: r.saveMsgsChan,
	}

	h.AddBatch(FromMessageBatch(msgs))
	h.MarkReaded()

	return h, nil
}

func (h *History) Add(msg MessageDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.current++
	h.newMsgs.Add(1)

	if h.msgCount < 100 {
		h.msgCount++
	}

	if h.current == 100 {
		h.current = 0
	}

	h.msgs[h.current] = msg
}

func (h *History) AddBatch(msgs []MessageDTO) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.msgCount += uint8(len(msgs))
	if h.msgCount > 100 {
		h.msgCount = 100
	}

	h.newMsgs.Add(uint32(len(msgs)))
	if h.newMsgs.Load() > 80 {
		h.stashChan <- struct{}{}
	}

	for _, msg := range msgs {
		h.current++

		if h.current == 100 {
			h.current = 0
		}

		h.msgs[h.current] = msg
	}

}

func (h *History) Read() []MessageDTO {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]MessageDTO, h.msgCount)

	msgsToRead := h.msgCount

	for i := h.current; msgsToRead > 0; i-- {
		if i == 0 {
			i = 99
		}

		out[msgsToRead-1] = h.msgs[i]
		msgsToRead--
	}

	return out
}

func (h *History) ReadNew() []MessageDTO {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]MessageDTO, h.newMsgs.Load())

	msgsToRead := h.newMsgs.Load()

	for i := h.current; msgsToRead > 0; i-- {
		if i == 0 {
			i = 100
		}

		out[msgsToRead-1] = h.msgs[i]
		msgsToRead--
	}

	return out
}

func (h *History) MarkReaded() {
	h.newMsgs.Swap(0)
}
