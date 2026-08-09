package main

import "sort"

// CCAATree almacena la relación jerárquica CCAA -> Provincias
type CCAATree map[string][]string

// BuildCCAATree procesa el slice de municipios en RAM y devuelve la jerarquía de regiones ordenada
func BuildCCAATree(municipis []Municipio) (CCAATree, []string) {
	tree := make(CCAATree)
	provSeen := make(map[string]bool)

	for _, m := range municipis {
		if m.CCAA == "" || m.Provincia == "" {
			continue
		}
		key := m.CCAA + "|" + m.Provincia
		if !provSeen[key] {
			provSeen[key] = true
			tree[m.CCAA] = append(tree[m.CCAA], m.Provincia)
		}
	}

	var ccaaKeys []string
	for ccaa, provs := range tree {
		sort.Strings(provs)
		ccaaKeys = append(ccaaKeys, ccaa)
	}
	sort.Strings(ccaaKeys)

	return tree, ccaaKeys
}
