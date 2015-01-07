package idgen

// ts(22) | tag(5) | wid(5) | seq(12)
func DecodeId(id int64) (ts int64, tag int64, wid int64, seq int64) {
	ts = (id >> uint64(TimestampLeftShift)) + twepoch
	seq = id & SequenceMask
	tag = (id >> TagIdShift) & ((1 << TagIdBits) - 1)
	wid = (id >> WorkerIdShift) & ((1 << WorkerIdBits) - 1)
	return
}
