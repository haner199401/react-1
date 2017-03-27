package banana

import (
	"bytes"
	"errors"
	"path"
	"strings"
)

var _ Node = new(App)

type App struct {
	*rootNode
	prefix string

	_TaggingScreen *Tagging
	_Model         *BytesBufferPLeaf
}

func newApp(r *rootNode, prefix string) *App {
	prefix = path.Join(prefix, "App")

	res := &App{
		rootNode: r,
		prefix:   prefix,
	}
	res._TaggingScreen = newTagging(r, prefix)
	res._Model = newBytesBufferPLeaf(r, prefix)
	return res
}

func (n *App) Subscribe(cb func()) *Sub {
	return n.rootNode.subscribe(n.prefix, cb)
}
func (n *App) TaggingScreen() *Tagging {
	return n._TaggingScreen
}
func (n *App) Model() *BytesBufferPLeaf {
	return n._Model
}

var _ Node = new(Tagging)

type Tagging struct {
	*rootNode
	prefix string

	_Name *StringLeaf
}

func newTagging(r *rootNode, prefix string) *Tagging {
	prefix = path.Join(prefix, "Tagging")

	res := &Tagging{
		rootNode: r,
		prefix:   prefix,
	}
	res._Name = newStringLeaf(r, prefix)
	return res
}

func (n *Tagging) Subscribe(cb func()) *Sub {
	return n.rootNode.subscribe(n.prefix, cb)
}
func (n *Tagging) Name() *StringLeaf {
	return n._Name
}

type BytesBufferPLeaf struct {
	*rootNode
	prefix string
}

var _ Node = new(BytesBufferPLeaf)

func newBytesBufferPLeaf(r *rootNode, prefix string) *BytesBufferPLeaf {
	prefix = path.Join(prefix, "BytesBufferPLeaf")

	return &BytesBufferPLeaf{
		rootNode: r,
		prefix:   prefix,
	}
}

func (m *BytesBufferPLeaf) Get() *bytes.Buffer {
	var res *bytes.Buffer
	if v, ok := m.rootNode.get(m.prefix); ok {
		return v.(*bytes.Buffer)
	}
	return res
}

func (m *BytesBufferPLeaf) Set(v *bytes.Buffer) {
	m.rootNode.set(m.prefix, v)
}

func (m *BytesBufferPLeaf) Subscribe(cb func()) *Sub {
	return m.rootNode.subscribe(m.prefix, cb)
}

type StringLeaf struct {
	*rootNode
	prefix string
}

var _ Node = new(StringLeaf)

func newStringLeaf(r *rootNode, prefix string) *StringLeaf {
	prefix = path.Join(prefix, "StringLeaf")

	return &StringLeaf{
		rootNode: r,
		prefix:   prefix,
	}
}

func (m *StringLeaf) Get() string {
	var res string
	if v, ok := m.rootNode.get(m.prefix); ok {
		return v.(string)
	}
	return res
}

func (m *StringLeaf) Set(v string) {
	m.rootNode.set(m.prefix, v)
}

func (m *StringLeaf) Subscribe(cb func()) *Sub {
	return m.rootNode.subscribe(m.prefix, cb)
}
func NewRoot() *App {
	r := &rootNode{
		store: make(map[string]interface{}),
		cbs:   make(map[string]map[*Sub]struct{}),
		subs:  make(map[*Sub]struct{}),
	}

	return newApp(r, "")
}

type Node interface {
	Subscribe(cb func()) *Sub
}

type Sub struct {
	*rootNode
	prefix string
	cb     func()
}

func (s *Sub) Clear() {
	s.rootNode.unsubscribe(s)
}

var NoSuchSubErr = errors.New("No such sub")

type rootNode struct {
	store map[string]interface{}
	cbs   map[string]map[*Sub]struct{}
	subs  map[*Sub]struct{}
}

func (r *rootNode) subscribe(prefix string, cb func()) *Sub {

	res := &Sub{
		cb:       cb,
		prefix:   prefix,
		rootNode: r,
	}

	l, ok := r.cbs[prefix]
	if !ok {
		l = make(map[*Sub]struct{})
		r.cbs[prefix] = l
	}

	l[res] = struct{}{}
	r.subs[res] = struct{}{}

	return res
}

func (r *rootNode) unsubscribe(s *Sub) {
	if _, ok := r.subs[s]; !ok {
		panic(NoSuchSubErr)
	}

	l, ok := r.cbs[s.prefix]
	if !ok {
		panic("Real problems...")
	}

	delete(l, s)
	delete(r.subs, s)
}

func (r *rootNode) get(k string) (interface{}, bool) {
	v, ok := r.store[k]
	return v, ok
}

func (r rootNode) set(k string, v interface{}) {
	if curr, ok := r.store[k]; ok && v == curr {
		return
	}

	r.store[k] = v

	parts := strings.Split(k, "/")

	var subs []*Sub

	var kk string

	for _, p := range parts {
		kk = path.Join(kk, p)

		if ll, ok := r.cbs[kk]; ok {
			for k := range ll {
				subs = append(subs, k)
			}
		}

	}

	for _, s := range subs {
		s.cb()
	}
}