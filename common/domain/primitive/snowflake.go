package primitive

func GetId() int64 {
	return node.Generate().Int64()
}
