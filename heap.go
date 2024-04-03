package mincache

func (h itemHeap) Len() int {
	return len(h)
}

func (h itemHeap) Less(i, j int) bool {
	return h[i].Expiration < h[j].Expiration
}

func (h itemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *itemHeap) Push(x any) {
	n := len(*h)
	item := x.(*cacheItem)
	item.Index = n
	*h = append(*h, item)
}

func (h *itemHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*h = old[0 : n-1]
	return item
}
