package entity

import (
	"bytelyon-functions/pkg/service/s3"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var pathRegexp = regexp.MustCompile("([A-Z]+)")

type builder interface {
	Context(ctx context.Context) Entity
	Bucket(string) Entity
	Path(interface{}) Entity
	ID(interface{}) Entity
	Value(interface{}) Entity
	Type(interface{}) Entity
}

type Entity interface {
	builder
	Save() error
	Find() error
	Exists() bool
	Page(int32, string) error
}

type Model struct {
	ctx    context.Context
	bucket string
	path   string
	id     string
	data   []byte
	value  interface{}
}

func (m *Model) Save() error {
	return s3.NewClient(m.ctx).Put(m.ctx, m.bucket, m.path+m.id, m.data)
}

func (m *Model) Find() error {
	out, err := s3.NewClient(m.ctx).Get(m.ctx, m.bucket, m.path+m.id)
	if err != nil {
		return err
	}
	return json.Unmarshal(out, &m.value)
}

func (m *Model) Exists() bool {
	err := m.Find()
	return err != nil && strings.Contains(err.Error(), "NoSuchKey")
}

func (m *Model) Page(size int32, after string) error {
	keys, err := s3.NewClient(m.ctx).KeysAfter(m.ctx, size, m.bucket, m.path, after)

	var vals []interface{}
	var v interface{}
	var out []byte
	for _, key := range keys {
		out, err = s3.NewClient(m.ctx).Get(m.ctx, m.bucket, key)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = json.Unmarshal(out, &v)
		vals = append(vals, v)
	}
	b, _ := json.Marshal(vals)
	return json.Unmarshal(b, &m.value)
}

func (m *Model) Context(ctx context.Context) Entity {
	if ctx == nil {
		ctx = context.Background()
	}
	m.ctx = ctx
	return m
}

func (m *Model) Bucket(name string) Entity {
	m.bucket = name
	return m
}

func (m *Model) Value(v interface{}) Entity {
	m.data, _ = json.Marshal(&v)
	if m.path == "" && m.id == "" {
		return m.Type(v).Path(nil).ID(nil)
	} else if m.path == "" {
		return m.Type(v).Path(nil)
	} else if m.id == "" {
		return m.Type(v).ID(nil)
	} else {
		return m.Type(v)
	}
}

func (m *Model) Type(v interface{}) Entity {
	m.value = v
	return m
}

func (m *Model) Path(v interface{}) Entity {
	if v == nil && m.value == nil {
		return m
	}
	if v == nil {
		v = m.value
	}
	val := reflect.TypeOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	m.path = val.Name()
	m.path = pathRegexp.ReplaceAllString(m.path, `/$1`)
	m.path = strings.ToLower(m.path)
	if strings.HasPrefix(m.path, "/") {
		m.path = m.path[1:]
	}
	if !strings.HasSuffix(m.path, "/") {
		m.path += "/"
	}
	return m
}

func (m *Model) ID(v interface{}) Entity {

	// if the given id is null, use an "ID" field from the model value
	if v == nil {
		val := reflect.ValueOf(m.value)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		val = val.FieldByName("ID")
		if !val.IsValid() {
			panic("invalid model")
		}
		v = val
	}

	// define the string value of this id
	m.id = fmt.Sprintf("%s", v)

	// check if a filetype needs to be appended
	if !strings.HasSuffix(m.id, ".json") {
		m.id += ".json"
	}

	return m
}

func New(c ...context.Context) Entity {
	var ctx context.Context
	if c != nil && len(c) > 0 {
		ctx = c[0]
	}
	return new(Model).Context(ctx).Bucket("bytelyon")
}
