package media

type ProgressEvent struct {
	Method      string  `json:"method"`
	Pathname    string  `json:"pathName"`
	TrackIdx    int     `json:"trackIdx"`
	CurrentTime float64 `json:"currentTime"`
}

type MediaProgress struct {
	Track map[string]int
	Time  map[string]map[int]float64
}

func (p *MediaProgress) HandleEvent(event ProgressEvent) ProgressEvent {
	switch event.Method {
	case "saveTrack":
		return p.saveTrack(event.Pathname, event.TrackIdx)
	case "getTrack":
		return p.getTrack(event.Pathname)
	case "saveTime":
		return p.saveTime(event.Pathname, event.TrackIdx, event.CurrentTime)
	case "getTime":
		return p.getTime(event.Pathname, event.TrackIdx)
	}
	return ProgressEvent{}
}

func (p *MediaProgress) saveTrack(pathname string, trackIdx int) ProgressEvent {
	p.Track[pathname] = trackIdx
	return ProgressEvent{Pathname: pathname, TrackIdx: trackIdx, Method: "saveTrack"}
}

func (p *MediaProgress) getTrack(pathname string) ProgressEvent {
	trackIdx := p.Track[pathname]
	return ProgressEvent{Pathname: pathname, TrackIdx: trackIdx, Method: "getTrack"}
}

func (p *MediaProgress) getTime(pathname string, trackIdx int) ProgressEvent {
	time := p.Time[pathname][trackIdx]
	return ProgressEvent{Pathname: pathname, TrackIdx: trackIdx, CurrentTime: time, Method: "getTime"}
}

func (p *MediaProgress) saveTime(pathname string, trackIdx int, currentTime float64) ProgressEvent {

	if p.Time[pathname] == nil {
		p.Time[pathname] = make(map[int]float64)
	}
	p.Time[pathname][trackIdx] = currentTime
	return ProgressEvent{Pathname: pathname, TrackIdx: trackIdx, CurrentTime: currentTime, Method: "saveTime"}
}
