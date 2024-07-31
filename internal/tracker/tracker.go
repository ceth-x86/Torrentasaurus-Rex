package tracker

type BencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
