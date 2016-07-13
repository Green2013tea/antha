package ast

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/graph"
)

const (
	AllDeps  = iota // Follow all AST edges
	DataDeps        // Follow only consumer-producer edges
)

var (
	_ Command = &Incubate{}
	_ Command = &Mix{}
)

// Input to code generation. An abstract syntax tree generated via execution of
// an Antha protocol.
//
// The basic design philosophy is to capture the semantics of the Antha
// language while reducing the cases for code generation. A secondary goal is
// to ease the creation of the AST at runtime (e.g., online, incremental
// generation of nodes).
//
// Conveniently, a tree naturally expresses the single-use (i.e., destructive
// update) aspect of physical things, so the code generation keeps this
// representation longer than a traditional compiler flow would.
type Node interface {
	graph.Node
	NodeString() string
}

type Command interface {
	Node
	Requests() []Request // Requirements for device selection
	Output() interface{} // Output from compilation
	SetOutput(interface{})
}

// Use of a liquid component
type UseComp struct {
	From  []Node
	Value *wtype.LHComponent
}

func (a *UseComp) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Value interface{}
	}{
		Value: a.Value,
	})
}

// Incubate expression
type Incubate struct {
	From []Node
	Reqs []Request
	Time wunit.Time
	Temp wunit.Temperature
	Out  interface{}
}

func (a *Incubate) Requests() []Request {
	return a.Reqs
}

func (a *Incubate) Output() interface{} {
	return a.Out
}

func (a *Incubate) SetOutput(x interface{}) {
	a.Out = x
}

func (a *Incubate) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Requests interface{}
	}{
		Requests: a.Requests,
	})
}

// Mix expression
type Mix struct {
	From []Node
	Reqs []Request
	Inst *wtype.LHInstruction // Data for planner
	Out  interface{}
}

func (a *Mix) Requests() []Request {
	return a.Reqs
}

func (a *Mix) Output() interface{} {
	return a.Out
}

func (a *Mix) SetOutput(x interface{}) {
	a.Out = x
}

func (a *Mix) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Requests interface{}
	}{
		Requests: a.Requests,
	})
}

// Unordered collection of expressions
type Bundle struct {
	From []Node
}

func (a *Bundle) NodeString() string {
	return ""
}

// Low-level move instruction
type Move struct {
	From     []*UseComp
	FromLocs []Location
	ToLoc    Location
	Out      interface{}
}

func (a *Move) Requests() []Request {
	return nil
}

func (a *Move) Output() interface{} {
	return a.Out
}

func (a *Move) SetOutput(x interface{}) {
	a.Out = x
}

func (a *Move) NodeString() string {
	return ""
}

// View AST as a graph
type Graph struct {
	Nodes     []Node
	whichDeps int
}

func (a *Graph) NumNodes() int {
	return len(a.Nodes)
}

func (a *Graph) Node(i int) graph.Node {
	return a.Nodes[i]
}

// Return subset of nodes that match the predicate
func matching(pred func(Node) bool, nodes ...Node) (r []Node) {
	for _, n := range nodes {
		if !pred(n) {
			continue
		}
		r = append(r, n)
	}
	return
}

func notNil(n Node) bool {
	return n != nil
}

func setOut(n Node, i, deps int, x Node) {
	switch n := n.(type) {
	case *UseComp:
		n.From[i] = x
	case *Bundle:
		n.From[i] = x
	case *Mix:
		n.From[i] = x
	case *Incubate:
		n.From[i] = x
	case *Move:
		n.From[i] = x.(*UseComp)
	default:
		panic(fmt.Sprintf("ast.setOut: unknown node type %T", n))
	}
}

func getOut(n Node, i, deps int) Node {
	switch n := n.(type) {
	case *UseComp:
		return n.From[i]
	case *Bundle:
		return n.From[i]
	case *Mix:
		return n.From[i]
	case *Incubate:
		return n.From[i]
	case *Move:
		return n.From[i]
	default:
		panic(fmt.Sprintf("ast.getOut: unknown node type %T", n))
	}
}

func numOuts(n Node, deps int) int {
	switch n := n.(type) {
	case *UseComp:
		return len(n.From)
		return 1
	case *Bundle:
		return len(n.From)
	case *Mix:
		return len(n.From)
	case *Incubate:
		return len(n.From)
	case *Move:
		return len(n.From)
	default:
		panic(fmt.Sprintf("ast.numOuts: unknown node type %T", n))
	}
}

func (a *Graph) NumOuts(n graph.Node) int {
	return numOuts(n.(Node), a.whichDeps)
}

func (a *Graph) Out(n graph.Node, i int) graph.Node {
	return getOut(n.(Node), i, a.whichDeps)
}

func (a *Graph) SetOut(n Node, i int, x Node) {
	setOut(n.(Node), a.whichDeps, i, x)
}

type ToGraphOpt struct {
	Roots     []Node // Roots of program
	WhichDeps int    // Edges to follow when building graph
}

// Create a graph from a list of roots. Incude any referenced ast nodes in the
// resulting graph.
func ToGraph(opt ToGraphOpt) *Graph {
	g := &Graph{
		whichDeps: opt.WhichDeps,
	}

	seen := make(map[graph.Node]bool)
	for _, root := range opt.Roots {
		// Traverse doesn't use Graph.NumNodes() or Graph.Node(int), so we can pass
		// in our partially constructed graph to extract the reachable nodes in the
		// AST
		results, _ := graph.Visit(graph.VisitOpt{
			Graph: g,
			Root:  root,
			Visitor: func(n graph.Node) error {
				if seen[n] {
					return graph.NextNode
				}
				return nil
			},
		})

		for _, k := range results.Seen.Range() {
			if seen[k] {
				continue
			}
			g.Nodes = append(g.Nodes, k.(Node))
			seen[k] = true
		}
	}

	return g
}

// Construct the data dependencies between a set of commands.
func Deps(roots []Node) graph.Graph {
	g := ToGraph(ToGraphOpt{Roots: roots, WhichDeps: DataDeps})
	root := make(map[graph.Node]bool)
	for _, r := range roots {
		root[r] = true
	}
	return graph.Eliminate(graph.EliminateOpt{
		Graph: g,
		In: func(n graph.Node) bool {
			return root[n]
		},
	})
}
