package main

type TraceHopResponse struct {
	Hop    int    `json:"hop"`
	IPAddr string `json:"ip"`
	RTT    int64  `json:"rtt"`
}

func (hop *TraceHop) toTraceHopResponse(hopIndex int) *TraceHopResponse {
	if hop == nil {
		return &TraceHopResponse{Hop: hopIndex + 1, IPAddr: "", RTT: -1}
	} else {
		return &TraceHopResponse{
			Hop:    hopIndex + 1,
			IPAddr: hop.IPAddr.String(),
			RTT:    hop.RTT.Milliseconds(),
		}
	}
}
