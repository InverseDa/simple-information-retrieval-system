package src

type Pair struct {
	first  int
	second float64
}

// 大小为K的最大堆
type Heap struct {
	data []Pair
	k    int
}

func InitKHeap(k int) *Heap {
	return &Heap{
		data: []Pair{},
		k:    k,
	}
}

func (h *Heap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) Less(i, j int) bool {
	return h.data[i].second > h.data[j].second
}

func (h *Heap) Push(x interface{}) {
	h.data = append(h.data, x.(Pair))
	h.ShiftUp()
	if h.Len() > h.k {
		h.Pop()
	}
}

func (h *Heap) Pop() interface{} {
	n := h.Len()
	h.swap(0, n-1)
	x := h.data[n-1]
	h.data = h.data[:n-1]
	h.ShiftDown(0)
	return x
}

func (h *Heap) ShiftUp() {
	k := h.Len() - 1
	for k > 0 && h.Less(k/2, k) {
		h.swap(k/2, k)
		k /= 2
	}
}

func (h *Heap) ShiftDown(i int) {
	n := h.Len()
	for i*2+1 < n {
		j := i*2 + 1
		if j+1 < n && h.Less(j, j+1) {
			j++
		}
		if !h.Less(i, j) {
			break
		}
		h.swap(i, j)
		i = j
	}
}
