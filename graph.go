package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/cloudspannerecosystem/memefish/ast"
)

const (
	rootGraphName = "G"
)

var (
	interleaveEdgeAttrs = map[string]string{}
	foreignKeyEdgeAttrs = map[string]string{}
	tableNodeAttrs      = map[string]string{}
	subgraphAttrs       = map[string]string{}
	groupNodeAttrs      = map[string]string{}
	groupEdgeAttrs      = map[string]string{}
)

func init() {
	interleaveEdgeAttrs["fontsize"] = "7"
	interleaveEdgeAttrs["dir"] = "both"
	interleaveEdgeAttrs["arrowsize"] = "0.9"
	interleaveEdgeAttrs["penwidth"] = "1.0"
	interleaveEdgeAttrs["labelangle"] = "32"
	interleaveEdgeAttrs["labeldistance"] = "1.0"
	interleaveEdgeAttrs["arrowhead"] = "none"
	foreignKeyEdgeAttrs["fontsize"] = "10"
	foreignKeyEdgeAttrs["arrowsize"] = "0.9"
	foreignKeyEdgeAttrs["penwidth"] = "1.0"
	foreignKeyEdgeAttrs["labelangle"] = "32"
	foreignKeyEdgeAttrs["labeldistance"] = "1.0"
	foreignKeyEdgeAttrs["arrowtail"] = "diamond"
	foreignKeyEdgeAttrs["style"] = "dotted"
	foreignKeyEdgeAttrs["dir"] = "back"
	tableNodeAttrs["shape"] = "\"Mrecord\""
	tableNodeAttrs["fontsize"] = "10"
	tableNodeAttrs["margin"] = "\"0.07,0.05\""
	tableNodeAttrs["penwidth"] = "1.0"
	groupNodeAttrs["label"] = "\"\""
	groupNodeAttrs["shape"] = "none"
	groupNodeAttrs["style"] = "\"\""
	groupEdgeAttrs["arrowhead"] = "none"
	groupEdgeAttrs["color"] = "white"
	subgraphAttrs["rank"] = "same"
}

type Graph struct {
	gvg *gographviz.Graph
}

func NewGraph() (*Graph, error) {
	gvg := gographviz.NewGraph()
	if err := gvg.SetName(rootGraphName); err != nil {
		return nil, err
	}
	if err := gvg.SetDir(true); err != nil {
		return nil, err
	}
	if err := gvg.AddAttr(rootGraphName, "fontsize", "12"); err != nil {
		return nil, err
	}

	if err := gvg.AddAttr(rootGraphName, "dpi", "150"); err != nil {
		return nil, err
	}
	return &Graph{
		gvg: gvg,
	}, nil
}

func (g *Graph) String() string {
	return g.gvg.String()
}

func (g *Graph) ApplyTables(tables []*ast.CreateTable) error {
	groupSize := groupSize(len(tables))
	if err := g.AddGroups(groupSize); err != nil {
		return err
	}
	if err := g.AddTables(groupSize, tables); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (g *Graph) AddGroups(size int) error {
	for i := 1; i <= size; i++ {
		if err := g.AddGroupNode(groupName(i)); err != nil {
			return err
		}
		if i > 1 {
			if err := g.AddGroupEdge(groupName(i-1), groupName(i)); err != nil {
				return err
			}
		}
		if err := g.AddSubGraph(groupName(i)); err != nil {
			return err
		}
		if err := g.addNode(groupName(i), groupName(i), map[string]string{}); err != nil {
			return err
		}

	}
	return nil
}

func (g *Graph) AddTables(groupSize int, tables []*ast.CreateTable) error {
	for i, t := range tables {
		// Handle interleave
		if t.Cluster != nil {
			// Get the parent table name from the Path
			parentName := getPathName(t.Cluster.TableName)
			tableName := getPathName(t.Name)
			if err := g.AddInterleaveEdge(parentName, tableName); err != nil {
				return err
			}
		}

		// Handle constraints
		for _, c := range t.TableConstraints {
			opt := make(map[string]string)
			if c.Name != nil {
				opt["label"] = c.Name.Name
			}

			// Handle foreign key constraints
			if fk, ok := c.Constraint.(*ast.ForeignKey); ok {
				refTableName := getPathName(fk.ReferenceTable)
				tableName := getPathName(t.Name)
				if err := g.AddForeignKeyEdge(refTableName, tableName, opt); err != nil {
					return err
				}
			}
		}

		// Handle columns
		colStr := ""
		for _, c := range t.Columns {
			pkMark := ""
			// Check if column is part of primary key
			if t.PrimaryKeys != nil {
				for _, pk := range t.PrimaryKeys {
					if pk.Name.Name == c.Name.Name {
						pkMark = "* "
						break
					}
				}
			}

			// Get type string
			typeStr := getSchemaTypeString(c.Type)
			typeStr = strings.Replace(typeStr, ">", "\\>", -1)
			typeStr = strings.Replace(typeStr, "<", "\\<", -1)

			colStr = colStr + fmt.Sprintf("%s%s (%s)\\l", pkMark, c.Name.Name, typeStr)
		}

		opt := make(map[string]string)
		tableName := getPathName(t.Name)
		opt["label"] = fmt.Sprintf("\"{%s|%s}\"", tableName, colStr)
		if err := g.AddTableNode(groupName(int(math.Mod(float64(i), float64(groupSize)))+1), tableName, opt); err != nil {
			return err
		}
	}
	return nil
}

// Helper function to get the name from a Path
func getPathName(path *ast.Path) string {
	if len(path.Idents) > 0 {
		return path.Idents[len(path.Idents)-1].Name
	}
	return ""
}

// Helper function to convert memefish SchemaType to string
func getSchemaTypeString(t ast.SchemaType) string {
	switch typ := t.(type) {
	case *ast.ScalarSchemaType:
		return string(typ.Name)
	case *ast.SizedSchemaType:
		if typ.Max {
			return fmt.Sprintf("%s(MAX)", typ.Name)
		}
		// Handle the Size field based on its actual type
		switch size := typ.Size.(type) {
		case *ast.IntLiteral:
			return fmt.Sprintf("%s(%s)", typ.Name, size.Value)
		default:
			return fmt.Sprintf("%s(%s)", typ.Name, "?")
		}
	case *ast.ArraySchemaType:
		return fmt.Sprintf("ARRAY<%s>", getSchemaTypeString(typ.Item))
	default:
		return fmt.Sprintf("%T", t)
	}
}

func (g *Graph) AddInterleaveEdge(parent, table string) error {
	return g.addEdge(parent, table, interleaveEdgeAttrs)
}

func (g *Graph) AddGroupEdge(src, dst string) error {
	return g.addEdge(src, dst, groupEdgeAttrs)
}

func (g *Graph) AddForeignKeyEdge(parent, table string, opt map[string]string) error {
	return g.addEdge(parent, table, merge(foreignKeyEdgeAttrs, opt))
}

func (g *Graph) addEdge(src, dst string, attr map[string]string) error {
	return g.gvg.AddEdge(src, dst, true, attr)
}

func (g *Graph) AddTableNode(groupName string, table string, opt map[string]string) error {
	return g.addNode(groupName, table, merge(tableNodeAttrs, opt))
}

func (g *Graph) AddGroupNode(groupName string) error {
	return g.addNode(rootGraphName, groupName, groupNodeAttrs)
}

func (g *Graph) addNode(parentGraph, name string, attr map[string]string) error {
	return g.gvg.AddNode(parentGraph, name, attr)
}

func (g *Graph) AddSubGraph(name string) error {
	return g.gvg.AddSubGraph(rootGraphName, name, subgraphAttrs)
}

func merge(m1, m2 map[string]string) map[string]string {
	merged := map[string]string{}

	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}

func groupName(i int) string { return fmt.Sprintf("%d", i) }
func groupSize(l int) int    { return int(math.Ceil(math.Sqrt(float64(l)))) }
