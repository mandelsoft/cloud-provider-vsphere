/*
 * Copyright 2020 The Kubernetes Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 */

package loadbalancer

import (
	"testing"

	"github.com/vmware/go-vmware-nsxt/common"
)

func _checkTags(t *testing.T, msg string, tags Tags, tag ...common.Tag) {
	if len(tags) != len(tag) {
		t.Errorf("%s: length mismatch: expected %d entries, but found %d", msg, len(tag), len(tags))
	}
	for _, _t := range tag {
		if f, ok := tags[_t.Scope]; ok {
			if f != _t {
				t.Errorf("%s: tag %q mismatch: expected %v, but found %v", msg, _t.Scope, _t, f)
			}
		} else {
			t.Errorf("%s: tag with scope %q missing", msg, _t.Scope)
		}
	}
}

func _checkNormTags(t *testing.T, msg string, tags []common.Tag, tag ...common.Tag) {
	if len(tags) != len(tag) {
		t.Errorf("%s: length mismatch: expected %d entries, but found %d", msg, len(tag), len(tags))
		return
	}
	for i, _t := range tag {
		if tags[i] != _t {
			t.Errorf("%s: entry %d: tag %q mismatch: expected %v, but found %v", msg, i, _t.Scope, _t, tags[i])
		}
	}
}

func newTag(scope, tag string) common.Tag {
	return common.Tag{Scope: scope, Tag: tag}
}

func TestTagAdd(t *testing.T) {
	tags := Tags{}

	t1 := newTag("t1", "v1")
	t1a := newTag("t1", "v1a")
	t2 := newTag("t2", "v2")
	t3 := newTag("t3", "v3")

	n := tags.Add(t1, t2)

	_checkTags(t, "original tags still empty after add", tags)
	_checkTags(t, "simple add", n, t1, t2)

	tags = n

	n = tags.Add(t1a)
	_checkTags(t, "replacing keeps original unchanged", tags, t1, t2)
	_checkTags(t, "replace tag", n, t1a, t2)

	n = tags.Add(t3)
	_checkTags(t, "adding keeps original unchanged", tags, t1, t2)
	_checkTags(t, "add tag", n, t1, t2, t3)

	norm := n.Normalize()
	_checkNormTags(t, "Normalize tags", norm, t1, t2, t3)

	norm = Tags{}.Add(t3).Add(t2, t1).Normalize()
	_checkNormTags(t, "Normalize tags with other add order", norm, t1, t2, t3)
}
