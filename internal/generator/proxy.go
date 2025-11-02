package generator

import (
	"fmt"
	"log"
	"strings"

	"SubGen/internal/config"

	"gopkg.in/yaml.v3"
)

// ReplaceProxyProvidersAndEncodeBase64 parses the base YAML and replaces the
// top-level "proxy-providers" section with a generated mapping for the given
// subscriptions. It preserves the merge alias (<<: *p) via YAML node operations,
// removes comments, and returns the resulting YAML mapping encoded as base64.
func ReplaceProxyProvidersAndEncodeBase64(base string, subs []config.Subscription) (string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(base), &doc); err != nil {
		return "", err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		// ensure top-level mapping exists
		return "", fmt.Errorf("invalid base YAML: top-level mapping not found")
	}
	top := doc.Content[0]

	// Build providers mapping node
	providers := &yaml.Node{Kind: yaml.MappingNode}
	uses := &yaml.Node{Kind: yaml.SequenceNode}
	for _, s := range subs {
		// key: subscription name
		key := &yaml.Node{Kind: yaml.ScalarNode, Value: s.Name}

		// value: mapping with merge alias and url (and optional override)
		val := &yaml.Node{Kind: yaml.MappingNode}

		// merge: <<: *p (alias node)
		mergeKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "<<"}
		alias := &yaml.Node{Kind: yaml.AliasNode, Value: "p"}
		val.Content = append(val.Content, mergeKey, alias)

		// url: "..."
		urlKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "url"}
		urlVal := &yaml.Node{Kind: yaml.ScalarNode, Value: s.URL}
		val.Content = append(val.Content, urlKey, urlVal)

		if strings.TrimSpace(s.AdditionalPrefix) != "" {
			// override:
			ovKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "override"}
			ovVal := &yaml.Node{Kind: yaml.MappingNode}
			apKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "additional-prefix"}
			apVal := &yaml.Node{Kind: yaml.ScalarNode, Value: s.AdditionalPrefix}
			ovVal.Content = append(ovVal.Content, apKey, apVal)
			val.Content = append(val.Content, ovKey, ovVal)
		}

		if strings.TrimSpace(s.Path) != "" {
			pathKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "path"}
			pathVal := &yaml.Node{Kind: yaml.ScalarNode, Value: s.Path}
			val.Content = append(val.Content, pathKey, pathVal)
		}

		providers.Content = append(providers.Content, key, val)
		uses.Content = append(uses.Content, key)
	}

	// Replace or append the top-level key
	replaced := false
	for i := 0; i+1 < len(top.Content); i += 2 {
		k := top.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == "proxy-providers" {
			top.Content[i+1] = providers
			replaced = true
			break
		}
	}
	if !replaced {
		top.Content = append(top.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "proxy-providers"}, providers)
	}

	// Replace only nested key path: use.use -> sequence of provider names
	updatedNestedUse := false
	for i := 0; i+1 < len(top.Content); i += 2 {
		k := top.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == "use" {
			v := top.Content[i+1]
			if v.Kind == yaml.MappingNode {
				innerReplaced := false
				for j := 0; j+1 < len(v.Content); j += 2 {
					ik := v.Content[j]
					if ik.Kind == yaml.ScalarNode && ik.Value == "use" {
						v.Content[j+1] = uses
						innerReplaced = true
						break
					}
				}
				if !innerReplaced {
					v.Content = append(v.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "use"}, uses)
				}
			}
			// If not a mapping, do not replace the whole 'use' node per requirement
			updatedNestedUse = true
			break
		}
	}
	if !updatedNestedUse {
		log.Printf("Top-level 'use' not found, nested 'use' not updated")
	}

	// Remove all comments recursively
	removeComments(top)

	// Marshal and encode
	out, err := yaml.Marshal(top)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func removeComments(n *yaml.Node) {
	if n == nil {
		return
	}
	n.HeadComment = ""
	n.LineComment = ""
	n.FootComment = ""
	for _, c := range n.Content {
		removeComments(c)
	}
}
