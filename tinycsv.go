package gotinycsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type options []string

func (o options) timeLayout() string {
	if len(o) != 0 {
		return o[0]
	}
	return "2006.1.2"
}

func eachStructFieldRefs(ref reflect.Value) ([][]reflect.Value, error) {
	elem0t := ref.Index(0).Type()
	if elem0t.Kind() == reflect.Ptr {
		elem0t = elem0t.Elem()
	}
	if elem0t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("elements of slice must be struct")
	}
	refs := make([][]reflect.Value, ref.Len())
	for i := 0; i < ref.Len(); i++ {
		refs[i] = make([]reflect.Value, elem0t.NumField())
		elem := ref.Index(i)
		elemp := reflect.NewAt(elem.Type(), unsafe.Pointer(elem.UnsafeAddr())).Elem()
		if elemp.Kind() == reflect.Ptr {
			elemp = elemp.Elem()
		}
		for j := 0; j < elem0t.NumField(); j++ {
			field := elemp.Field(j)
			refs[i][j] = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		}
	}
	return refs, nil
}

func setEntityViaRef(ref reflect.Value, timelayout, v string) error {
	if !ref.CanSet() {
		return fmt.Errorf("cannot set via reference: %s", v)
	}
	switch ref.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv, _ := strconv.Atoi(v)
		ref.SetInt(reflect.ValueOf(iv).Int())
	case reflect.Float32:
		fv, _ := strconv.ParseFloat(v, 32)
		ref.SetFloat(reflect.ValueOf(fv).Float())
	case reflect.Float64:
		fv, _ := strconv.ParseFloat(v, 64)
		ref.SetFloat(reflect.ValueOf(fv).Float())
	case reflect.String:
		ref.SetString(strings.TrimSpace(reflect.ValueOf(v).String()))
	case reflect.Struct:
		switch ref.Interface().(type) {
		case time.Time:
			t, _ := time.Parse(timelayout, v)
			ref.Set(reflect.ValueOf(t))
		}
	}
	return nil
}

func sliceRefPointer(i interface{}) (*reflect.Value, error) {
	ref := reflect.ValueOf(i)
	if ref.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("failed to obtain a reference to i (did you forget &?)")
	}
	refp := ref.Elem()
	if refp.Kind() != reflect.Slice {
		return nil, fmt.Errorf("i reference does not point to a slice")
	}
	return &refp, nil
}

func ensureSliceCapacity(ref reflect.Value, len int) error {
	if !ref.CanAddr() {
		return fmt.Errorf("failed allocate slice capacity and is not addressable (did you forget &?)")
	}
	if ref.Len() < len {
		ref.Set(reflect.MakeSlice(ref.Type(), len, len))
	}
	if ref.Index(0).Kind() == reflect.Ptr {
		for i := 0; i < ref.Len(); i++ {
			refi := ref.Index(i)
			if refi.IsNil() {
				refi.Set(reflect.New(refi.Type().Elem()))
			}
		}
	}
	return nil
}

// Load a CSV.
// "r" is CSV format reader.
// Skip the "topmergin" lines from the top line.
// An error occurs when the number of rows read reaches "maxrows".
// "out" is load destination. automatically ensures optimal capacity.
// The first element of "ops" is time-layout
func Load(r io.Reader, topmergin int, maxrows int, out interface{}, ops ...string) error {
	refp, err := sliceRefPointer(out)
	if err != nil {
		return err
	}

	cr := csv.NewReader(r)

	records := make([][]*string, 0, maxrows)

	rows := 0
	for ; ; rows++ {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if rows < topmergin {
			continue
		}
		if err != nil {
			return err
		}
		if rows >= maxrows {
			return fmt.Errorf("rows are too large")
		}
		r := make([]*string, len(record))
		for i := range record {
			r[i] = &record[i]
		}
		records = append(records, r)
	}
	if rows <= topmergin {
		return fmt.Errorf("topmergin is too large")
	}

	rows -= topmergin

	// create "out" for all rows
	if err = ensureSliceCapacity(*refp, rows); err != nil {
		return err
	}

	// create slice of references to struct field
	refs, err := eachStructFieldRefs(*refp)
	if err != nil {
		return err
	}

	timelayout := options(ops).timeLayout()

	if len(refs[0]) < len(records[0]) {
		return fmt.Errorf("number of fields in the defined structure may not match the number of fields in the CSV.")
	}

	for rows := 0; rows < len(records); rows++ {
		// sets csv record into "out" via references
		for cols, v := range records[rows] {
			if err = setEntityViaRef(refs[rows][cols], timelayout, *v); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load a CSV with fileds arranged vertically.
// "r" is CSV format reader.
// Skip the "topmergin" lines from the top line.
// Skip the "leftmergin" columns from left edge.
// An error occurs when the number of columns read reaches "maxcols".
// "out" is load destination. automatically ensures optimal capacity.
// The first element of "ops" is time-layout
func LoadVertically(r io.Reader, topmergin int, leftmergin int, maxcols int, out interface{}, ops ...string) error {
	refp, err := sliceRefPointer(out)
	if err != nil {
		return err
	}

	cr := csv.NewReader(r)

	// discard header
	discards := topmergin
	if topmergin == 0 {
		discards = 1
	}

	var record []string

	for i := 0; i < discards; i++ {
		record, err = cr.Read()
		if err == io.EOF {
			return err
		}
	}

	if leftmergin >= len(record) {
		return fmt.Errorf("leftmergin is too large")
	}

	if len(record[leftmergin:]) > maxcols {
		return fmt.Errorf("columns are too large")
	}

	// create "out" for all rows
	if err = ensureSliceCapacity(*refp, len(record[leftmergin:])); err != nil {
		return err
	}

	// create slice of references to struct field
	refs, err := eachStructFieldRefs(*refp)
	if err != nil {
		return err
	}

	timelayout := options(ops).timeLayout()

	rows := 0
	// if topmergin is 0, stored the first line at first.
	if topmergin == 0 {
		// sets csv record into "out" via references
		for cols, v := range record[leftmergin:] {
			if err = setEntityViaRef(refs[cols][rows], timelayout, v); err != nil {
				return err
			}
		}
		rows++
	}

	for ; rows < len(refs[0]); rows++ {
		record, err = cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// sets csv record into "out" via references
		for cols, v := range record[leftmergin:] {
			if err = setEntityViaRef(refs[cols][rows], timelayout, v); err != nil {
				return err
			}
		}
	}

	return nil
}
