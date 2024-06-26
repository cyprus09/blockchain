package main

func main() {
	bc := NewBlockChain()

	bc.db.Close()
}
