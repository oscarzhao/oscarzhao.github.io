package josephus

import (
	"math/bits"
)

// JosephusBitMap 基于 bitmap 实现
// 为了简化代码，使用 []bool 替代 bitmap
func JosephusBitMap(n int) int {
	bitmap := make([]bool, n, n)
	toDel := false
	for left := n; left > 1; {
		for idx := 0; idx < n; idx++ {
			if bitmap[idx] {
				continue
			}
			if toDel {
				bitmap[idx] = true
				toDel = false
				left--
			} else {
				toDel = true
			}
		}
	}
	for i := 0; i < n; i++ {
		if !bitmap[i] {
			return i + 1
		}
	}
	return -1
}

// Node 定义一个链表节点
type Node struct {
	val  int
	next *Node
}

// JosephusLinklist 是基于环形链表的实现
func JosephusLinklist(n int) int {
	// 初始化环形链表
	head := &Node{val: 1}
	cur := head
	for i := 2; i <= n; i++ {
		cur.next = &Node{val: i}
		cur = cur.next
	}
	cur.next = head

	// 删除节点
	for cur := head; ; {
		next := cur.next
		cur.next = next.next
		cur = cur.next
		// 终止条件：只有一个节点
		if cur.next == cur {
			return cur.val
		}
	}

	// 返回结果
	return -1
}

// JosephusRecursion 是递归实现
func JosephusRecursion(n int) int {
	if n == 1 {
		return 1
	}
	if n%2 == 0 {
		return 2*JosephusRecursion(n/2) - 1
	} else {
		return 2*JosephusRecursion(n/2) + 1
	}
}

// JosephusBit 是二进制实现
func JosephusBit(n int) int {
	uintN := uint(n)
	leftMove := bits.UintSize - bits.LeadingZeros(uintN) - 1
	mask := (uint(1) << leftMove) - 1
	i := mask & uintN
	return int(i*2 + 1)
}
