package utils

type IdList []uint

func (cl IdList) Len() int           { return len(cl) }
func (cl IdList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl IdList) Less(i, j int) bool { return cl[i] < cl[j] }
