package idgen

func DecodeId(id int64) (ts int64, tag int64, seq int64) {
	ts = (id >> uint64(TimestampLeftShift)) + twepoch
	seq = id & SequenceMask
	tag = (id >> SequenceBits) & ((1 << (TimestampLeftShift + WorkerIdShift + TagIdShift)) - 1)
	return
}
