package twitch

type Listener struct {
	Id       string
	Streamer string
	Callback func(*TwitchMessage)
}

func (ts *TwitchService) JoinStreamer(listener Listener) {
	ts.StreamsMutex.Lock()
	defer ts.StreamsMutex.Unlock()
	channels := ts.Streams[listener.Streamer]

	if len(channels) == 0 {
		ts.client.Join(listener.Streamer)
	}

	ts.Streams[listener.Streamer] = append(channels, listener)
}

func (ts *TwitchService) LeaveStreamer(listener Listener) {
	ts.StreamsMutex.Lock()
	defer ts.StreamsMutex.Unlock()
	listeners, isOK := ts.Streams[listener.Streamer]

	if !isOK {
		return
	}

	for i, l := range listeners {
		if listener.Id == l.Id {
			ts.Streams[listener.Streamer] = append(listeners[:i], listeners[i+1:]...)
		}
	}

	if len(ts.Streams[listener.Streamer]) == 0 {
		ts.client.Depart(listener.Streamer)
	}
}
