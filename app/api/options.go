package api

import "golang-exercise/writer"

func Writer(w *writer.Writer) func(*handler) {
	return func(h *handler) {
		h.writer = w
	}
}
