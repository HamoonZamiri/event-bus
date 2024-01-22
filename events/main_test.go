package main

import (
	"reflect"
	"testing"
)

func TestFindEvent(t *testing.T) {
	body := []byte(`{"type": "comment_created", "data": {"id": "1", "post_id": "1", "content": "hello"}}`)
	expected := "comment_created"
	reflect.DeepEqual(findEventType(body), expected)
}

func TestFindEvent2(t *testing.T) {
	body := []byte(`{"type": "post_created", "data": {"id": "1", "title": "hello"}}`)
	expected := "post_created"
	reflect.DeepEqual(findEventType(body), expected)
}

func TestFindEventDomain(t *testing.T) {
	event := "comment_created"
	expected := "comment"
	reflect.DeepEqual(FindEventDomain(event), expected)
}

func TestFindEventDomain2(t *testing.T) {
	event := "post_created"
	expected := "post"
	reflect.DeepEqual(FindEventDomain(event), expected)
}
