package main

import (
	"fmt"
)

type LinkNode struct {
	Value int
	Next  *LinkNode
}

func merge(a, b *LinkNode) *LinkNode {
	head := new(LinkNode)
	n := head
	for ; a != nil && b != nil; n = n.Next {
		if a.Value < b.Value {
			n.Next = a
			a = a.Next
		} else {
			n.Next = b
			b = b.Next
		}
	}
	if a != nil {
		n.Next = a
	}
	if b != nil {
		n.Next = b
	}
	return head.Next
}

func NodeSort(head *LinkNode) *LinkNode {
	if head == nil || head.Next == nil {
		return head
	}
	perSlow, slow, fast := head, head, head
	for fast != nil && fast.Next != nil {
		perSlow = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	perSlow.Next = nil
	head = NodeSort(head)
	slow = NodeSort(slow)

	return merge(head, slow)
}

func main() {
	head := new(LinkNode)
	var n int
	fmt.Scanf("%d", &n)
	Node1 := new(LinkNode)
	for i := 0; i < n; i++ {
		k := 0
		fmt.Scanf("%d", &k)
		if i == 0 {
			head.Value = k
		} else if i == 1 {
			Node1.Value = k
			head.Next = Node1
		} else {
			Node2 := new(LinkNode)
			Node2.Value = k
			Node1.Next = Node2
			Node1 = Node2
		}
	}
	Node := NodeSort(head)
	readNode := Node
	for readNode != nil {
		fmt.Printf("%d ", readNode.Value)
		readNode = readNode.Next
	}
}
