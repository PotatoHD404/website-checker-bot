package main

func main() {
	w := NewWebsite("test", "https://kaf22.ru", "")
	changed, err := w.CheckChanged()
	if err != nil {
		return
	}
	if changed {
		println("changed")
	} else {
		println("not changed")
	}
}
