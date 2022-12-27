package main

import "fmt"

type LinkNode struct {
	Value int
	Next  *LinkNode
}

func intersection(node1 **LinkNode, node2 **LinkNode) *LinkNode {
	head := new(LinkNode)
	writeNode := head
	A := new(*LinkNode)
	A = node1
	B := new(*LinkNode)
	B = node2
	m := make(map[int]bool)
	Aa := new(LinkNode)
	Aa = *A
	Bb := new(LinkNode)
	Bb = *B
	fmt.Println(Bb.Next.Next.Next.Value)
	fmt.Println(Aa.Next.Next.Next.Value)
	for Aa != nil && Bb != nil {
		if Aa.Value == Bb.Value {
			if !m[Aa.Value] {
				m[Aa.Value] = true
				writeNode1 := new(LinkNode)
				writeNode1.Value = Aa.Value
				Aa = Aa.Next
				Bb = Bb.Next
				writeNode.Next = writeNode1
				writeNode = writeNode.Next
			} else {
				fmt.Println(11211)
				Aa = Aa.Next
				Bb = Bb.Next
			}
		} else if Aa.Next != nil && Bb.Next != nil {
			fmt.Println(2222)
			Bb = Bb.Next
		} else if Bb.Next == nil {
			fmt.Println(3333)
			Aa = Aa.Next
			Bb = *node2
		} else {
			Bb = Bb.Next
		}
	}
	fmt.Println(head.Next)
	return head.Next
}

func main() {
	node1 := new(LinkNode)
	node2 := new(LinkNode)
	a, b := new(LinkNode), new(LinkNode)
	a = node1
	b = node2
	a1, b1 := new(LinkNode), new(LinkNode)
	a.Value = 1
	a1.Value = 3
	a.Next = a1
	a = a.Next
	a1 = new(LinkNode)
	a1.Value = 5
	a.Next = a1
	a = a.Next
	a1 = new(LinkNode)
	a1.Value = 7
	a.Next = a1
	b.Value = 2
	b1.Value = 1
	b.Next = b1
	b = b.Next
	b1 = new(LinkNode)
	b1.Value = 8
	b.Next = b1
	b = b.Next
	b1 = new(LinkNode)
	b1.Value = 7
	b.Next = b1
	//for node1 != nil {
	//	fmt.Printf("%d ", node1.Value)
	//	node1 = node1.Next
	//}
	//fmt.Println()
	//for node2 != nil {
	//	fmt.Printf("%d ", node2.Value)
	//	node2 = node2.Next
	//}
	readNode := intersection(&node1, &node2)
	fmt.Printf("Set1∩Set2 : ")
	for readNode != nil {
		fmt.Printf("%d ", readNode.Value)
		readNode = readNode.Next
	}
	fmt.Println("Set1∪Set2 : ")
	fmt.Printf("Set1-Set2 : ")
	fmt.Printf("Set2-Set1 : ")
}
