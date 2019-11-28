/*
Copyright 2019 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package typed

import (
	"sigs.k8s.io/structured-merge-diff/fieldpath"
	"sigs.k8s.io/structured-merge-diff/schema"
	"sigs.k8s.io/structured-merge-diff/value"
)

type removingWalker struct {
	value    *value.Value
	schema   *schema.Schema
	toRemove *fieldpath.Set
}

func removeItemsWithSchema(value *value.Value, toRemove *fieldpath.Set, schema *schema.Schema, typeRef schema.TypeRef) {
	w := &removingWalker{
		value:    value,
		schema:   schema,
		toRemove: toRemove,
	}
	resolveSchema(schema, typeRef, value, w)
}

// doLeaf should be called on leaves before descending into children, if there
// will be a descent. It modifies w.inLeaf.
func (w *removingWalker) doLeaf() ValidationErrors { return nil }

func (w *removingWalker) doScalar(t *schema.Scalar) ValidationErrors { return nil }

func (w *removingWalker) doList(t *schema.List) (errs ValidationErrors) {
	l := value.ValueList(*w.value)

	// If list is null, empty, or atomic just return
	if l == nil || len(l) == 0 || t.ElementRelationship == schema.Atomic {
		return nil
	}

	var newItems []interface{}
	for i, item := range l {
		// Ignore error because we have already validated this list
		pe, _ := listItemToPathElement(t, i, item)
		path, _ := fieldpath.MakePath(pe)
		if w.toRemove.Has(path) {
			continue
		}
		if subset := w.toRemove.WithPrefix(pe); !subset.Empty() {
			val := value.Value(item)
			removeItemsWithSchema(&val, subset, w.schema, t.ElementType)
			item = val

		}
		newItems = append(newItems, item)
	}
	if len(newItems) > 0 {
		*w.value = newItems
	} else {
		*w.value = nil
	}
	return nil
}

func (w *removingWalker) doMap(t *schema.Map) ValidationErrors {
	m := value.ValueMap(*w.value)

	// If map is null, empty, or atomic just return
	if len(m) == 0 || t.ElementRelationship == schema.Atomic {
		return nil
	}

	fieldTypes := map[string]schema.TypeRef{}
	for _, structField := range t.Fields {
		fieldTypes[structField.Name] = structField.Type
	}

	newMap := map[string]interface{}{}
	for key, val := range m {
		k := key
		pe := fieldpath.PathElement{FieldName: &k}
		path, _ := fieldpath.MakePath(pe)
		fieldType := t.ElementType
		if ft, ok := fieldTypes[key]; ok {
			fieldType = ft
		} else {
			if w.toRemove.Has(path) {
				continue
			}
		}
		if subset := w.toRemove.WithPrefix(pe); !subset.Empty() {
			v := value.Value(val)
			removeItemsWithSchema(&v, subset, w.schema, fieldType)
			val = v
		}
		newMap[key] = val
	}
	if len(newMap) > 0 {
		*w.value = newMap
	} else {
		*w.value = nil
	}
	return nil
}
