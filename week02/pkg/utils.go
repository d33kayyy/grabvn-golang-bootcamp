package pkg

func Reduce(m1, m2 map[string]int) map[string]int {
	// first merge m2 to m1
	for k := range m1 {
		v, ok := m2[k]
		if ok {
			m1[k] += v
		}
	}
	// add missing word to m1
	for k, v := range m2 {
		_, ok := m1[k]
		if !ok {
			m1[k] = v
		}
	}
	return m1
}

//func main() {
//	m1 := map[string]int{"a": 3, "b": 2}
//	m2 := map[string]int{"a": 1, "c": 3}
//	r := Reduce(m2, m1)
//	fmt.Println(r)
//}
