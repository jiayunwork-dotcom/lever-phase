package phasediag

import "fmt"

// Node 表示相界折线上的一个结点：成分（质量分数）与温度（开尔文）。
type Node struct {
	Composition float64 `json:"composition"`
	Temperature float64 `json:"temperature"`
}

// NewNode 构造一个相界结点。
func NewNode(comp, temp float64) Node {
	return Node{Composition: comp, Temperature: temp}
}

// Comp 返回结点成分（质量分数）。
func (n Node) Comp() float64 { return n.Composition }

// Temp 返回结点温度（开尔文）。
func (n Node) Temp() float64 { return n.Temperature }

// String 输出结点的人类可读描述，用于错误文案与调试。
func (n Node) String() string {
	return fmt.Sprintf("(c=%.6f, T=%.3f K)", n.Composition, n.Temperature)
}
