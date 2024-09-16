package protocol

var nullBulkBytes = []byte("$-1\r\n")

// NullBulkReply is empty string
type NullBulkReply struct{}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var emptyMultiBulkBytes = []byte("*0\r\n")

// EmptyMultiBulkReply means is an empty list
type EmptyMultiBulkReply struct{}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}
func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}
