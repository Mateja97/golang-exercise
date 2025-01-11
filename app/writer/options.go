package writer

func Brokers(v ...string) func(*Writer) {
	return func(w *Writer) {
		w.brokers = v
	}
}

func DestinationTopic(t string) func(writer *Writer) {
	return func(w *Writer) {
		w.destinationTopic = t
	}
}
