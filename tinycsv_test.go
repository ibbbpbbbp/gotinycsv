package gotinycsv

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_eachStructFieldRefs(t *testing.T) {
	// normal case
	{
		type teststruct struct {
			a int
			b string
		}

		slice := []teststruct{{a: 0, b: ""}, {a: 0, b: ""}, {a: 0, b: ""}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		assert.Equal(t, 3, len(refs))
		assert.Equal(t, 2, len(refs[0]))
		assert.Equal(t, 2, len(refs[1]))
		assert.Equal(t, 2, len(refs[2]))
		assert.Equal(t, unsafe.Pointer(&(slice[0].a)), unsafe.Pointer(refs[0][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[0].b)), unsafe.Pointer(refs[0][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].a)), unsafe.Pointer(refs[1][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].b)), unsafe.Pointer(refs[1][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].a)), unsafe.Pointer(refs[2][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].b)), unsafe.Pointer(refs[2][1].UnsafeAddr()))
	}
	// illegal case 1 (slice elements are not struct type)
	{
		slice := []int{0, 0, 0}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.EqualError(t, err, "elements of slice must be struct")
		assert.Nil(t, refs)
	}
	// illegal case 2 (slice elements are nesting)
	{
		type substruct struct {
			c int
			d string
		}

		type teststruct struct {
			a string
			b substruct
		}

		// cannot access b.c and b.d from refs
		slice := []teststruct{{"", substruct{0, ""}}, {"", substruct{0, ""}}, {"", substruct{0, ""}}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		assert.Equal(t, 3, len(refs))
		assert.Equal(t, 2, len(refs[0]))
		assert.Equal(t, 2, len(refs[1]))
		assert.Equal(t, 2, len(refs[2]))
		assert.Equal(t, unsafe.Pointer(&(slice[0].a)), unsafe.Pointer(refs[0][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[0].b)), unsafe.Pointer(refs[0][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].a)), unsafe.Pointer(refs[1][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].b)), unsafe.Pointer(refs[1][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].a)), unsafe.Pointer(refs[2][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].b)), unsafe.Pointer(refs[2][1].UnsafeAddr()))
	}
	// illegal case 3 (not supported type (slice) is exist in slice)
	{
		type teststruct struct {
			a string
			b []int
		}

		// cannnot access b[i] from refs
		slice := []teststruct{{"", []int{0, 0}}, {"", []int{0, 0}}, {"", []int{0, 0}}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		assert.Equal(t, 3, len(refs))
		assert.Equal(t, 2, len(refs[0]))
		assert.Equal(t, 2, len(refs[1]))
		assert.Equal(t, 2, len(refs[2]))
		assert.Equal(t, unsafe.Pointer(&(slice[0].a)), unsafe.Pointer(refs[0][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[0].b)), unsafe.Pointer(refs[0][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].a)), unsafe.Pointer(refs[1][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].b)), unsafe.Pointer(refs[1][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].a)), unsafe.Pointer(refs[2][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].b)), unsafe.Pointer(refs[2][1].UnsafeAddr()))
	}
	// illegal case 4 (not supported type (pointer) is exist in slice)
	{
		type teststruct struct {
			a string
			b *int
		}

		// cannnot access *b from refs
		// access itself is easy, but there is no guarantee that memory will be allocated as in this test, so it is not supported.
		slice := []teststruct{{"", nil}, {"", nil}, {"", nil}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		assert.Equal(t, 3, len(refs))
		assert.Equal(t, 2, len(refs[0]))
		assert.Equal(t, 2, len(refs[1]))
		assert.Equal(t, 2, len(refs[2]))
		assert.Equal(t, unsafe.Pointer(&(slice[0].a)), unsafe.Pointer(refs[0][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[0].b)), unsafe.Pointer(refs[0][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].a)), unsafe.Pointer(refs[1][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[1].b)), unsafe.Pointer(refs[1][1].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].a)), unsafe.Pointer(refs[2][0].UnsafeAddr()))
		assert.Equal(t, unsafe.Pointer(&(slice[2].b)), unsafe.Pointer(refs[2][1].UnsafeAddr()))
	}
}

func Test_setEntityViaRef(t *testing.T) {
	// normal case
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := []teststruct{{}, {}, {}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		setEntityViaRef(refs[0][0], "2006-01-02", "10")
		setEntityViaRef(refs[0][1], "2006-01-02", "aa")
		setEntityViaRef(refs[0][2], "2006-01-02", "2022-01-01")
		setEntityViaRef(refs[1][0], "2006-01-02", "20")
		setEntityViaRef(refs[1][1], "2006-01-02", "bb")
		setEntityViaRef(refs[1][2], "2006-01-02", "2022-01-02")
		setEntityViaRef(refs[2][0], "2006-01-02", "30")
		setEntityViaRef(refs[2][1], "2006-01-02", "cc")
		setEntityViaRef(refs[2][2], "2006-01-02", "2022-01-03")
		assert.Equal(t, 10, slice[0].a)
		assert.Equal(t, "aa", slice[0].b)
		assert.Equal(t, "2022-01-01 00:00:00 +0000 UTC", slice[0].c.String())
		assert.Equal(t, 20, slice[1].a)
		assert.Equal(t, "bb", slice[1].b)
		assert.Equal(t, "2022-01-02 00:00:00 +0000 UTC", slice[1].c.String())
		assert.Equal(t, 30, slice[2].a)
		assert.Equal(t, "cc", slice[2].b)
		assert.Equal(t, "2022-01-03 00:00:00 +0000 UTC", slice[2].c.String())
	}
	// illegal case 1 (passed value that cannot be converted)
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := []teststruct{{}, {}, {}}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		setEntityViaRef(refs[0][0], "2006-01-02", "not number") // cannot convert int
		setEntityViaRef(refs[0][1], "2006-01-02", "aa")
		setEntityViaRef(refs[0][2], "2006-01-02", "illegal time format") // cannot convert time.Time
		setEntityViaRef(refs[1][0], "2006-01-02", "not number")          // cannot convert int
		setEntityViaRef(refs[1][1], "2006-01-02", "bb")
		setEntityViaRef(refs[1][2], "2006-01-02", "illegal time format")
		setEntityViaRef(refs[2][0], "2006-01-02", "not number") // cannot convert int
		setEntityViaRef(refs[2][1], "2006-01-02", "cc")
		setEntityViaRef(refs[2][2], "2006-01-02", "illegal time format") // cannot convert time.Time
		assert.Equal(t, 0, slice[0].a)
		assert.Equal(t, "aa", slice[0].b)
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", slice[0].c.String())
		assert.Equal(t, 0, slice[1].a)
		assert.Equal(t, "bb", slice[1].b)
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", slice[1].c.String())
		assert.Equal(t, 0, slice[2].a)
		assert.Equal(t, "cc", slice[2].b)
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", slice[2].c.String())
	}
	// illegal case 2 (entities are not supported types)
	{
		type teststruct struct {
			a *int        // not supported type
			b []int       // not supported type
			c map[int]int // not supported type
		}

		slice := []teststruct{
			{new(int), make([]int, 0), make(map[int]int)},
			{new(int), make([]int, 0), make(map[int]int)},
			{new(int), make([]int, 0), make(map[int]int)},
		}
		ref := reflect.ValueOf(slice)
		refs, err := eachStructFieldRefs(ref)
		assert.NoError(t, err)
		assert.NotNil(t, refs)
		setEntityViaRef(refs[0][0], "2006-01-02", "10")
		setEntityViaRef(refs[0][1], "2006-01-02", "aa")
		setEntityViaRef(refs[0][2], "2006-01-02", "2022-01-01")
		setEntityViaRef(refs[1][0], "2006-01-02", "20")
		setEntityViaRef(refs[1][1], "2006-01-02", "bb")
		setEntityViaRef(refs[1][2], "2006-01-02", "2022-01-02")
		setEntityViaRef(refs[2][0], "2006-01-02", "30")
		setEntityViaRef(refs[2][1], "2006-01-02", "cc")
		setEntityViaRef(refs[2][2], "2006-01-02", "2022-01-03")
		assert.Zero(t, *slice[0].a)
		assert.Empty(t, slice[0].b)
		assert.Empty(t, slice[0].c)
		assert.Zero(t, *slice[1].a)
		assert.Empty(t, slice[1].b)
		assert.Empty(t, slice[1].c)
		assert.Zero(t, *slice[1].a)
		assert.Empty(t, slice[1].b)
		assert.Empty(t, slice[1].c)
	}
}

func Test_sliceRefPointer(t *testing.T) {
	// normal case
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := &[]*teststruct{}
		refp, err := sliceRefPointer(slice)
		assert.NoError(t, err)
		assert.NotNil(t, refp)
	}
	// illegal case 1 (did not pass pointer of slice)
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := []*teststruct{}
		refp, err := sliceRefPointer(slice)
		assert.EqualError(t, err, "failed to obtain a reference to i (did you forget &?)")
		assert.Nil(t, refp)
	}
	// illegal case 2 (pointer does not point to a slice)
	{
		var str string
		slice := &str
		refp, err := sliceRefPointer(slice)
		assert.EqualError(t, err, "i reference does not point to a slice")
		assert.Nil(t, refp)
	}
}

func Test_ensureSliceCapacity(t *testing.T) {
	// normal case 1
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := &[]teststruct{}
		refp, err := sliceRefPointer(slice)
		assert.NoError(t, err)
		assert.NotNil(t, refp)
		err = ensureSliceCapacity(*refp, 10)
		assert.Equal(t, 10, len(*slice))
	}
	// normal case 2
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := &[]*teststruct{}
		refp, err := sliceRefPointer(slice)
		assert.NoError(t, err)
		assert.NotNil(t, refp)
		err = ensureSliceCapacity(*refp, 10)
		assert.Equal(t, 10, len(*slice))
		// check that slice element (pointer) is also allocated.
		for _, v := range *slice {
			assert.NotNil(t, v)
		}
	}
	// normal case 3
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := make([]*teststruct, 100)
		refp, err := sliceRefPointer(&slice)
		assert.NoError(t, err)
		assert.NotNil(t, refp)
		err = ensureSliceCapacity(*refp, 10)
		// check that slice is not shrinking (100 -> 10).
		assert.Equal(t, 100, len(slice))
		// but, slice element (pointer) is allocated.
		for _, v := range slice {
			assert.NotNil(t, v)
		}
	}
	// normal case 4
	{
		type teststruct struct {
			a int
			b string
			c time.Time
		}

		slice := make([]*teststruct, 100)
		bforeaddresses := make([]unsafe.Pointer, 100)
		afteraddresses := make([]unsafe.Pointer, 100)
		for i := range slice {
			slice[i] = &teststruct{}
			bforeaddresses[i] = unsafe.Pointer(slice[i])
		}
		refp, err := sliceRefPointer(&slice)
		assert.NoError(t, err)
		assert.NotNil(t, refp)
		err = ensureSliceCapacity(*refp, 10)
		// check that slice is not shrinking (100 -> 10).
		assert.Equal(t, 100, len(slice))
		for i, v := range slice {
			afteraddresses[i] = unsafe.Pointer(v)
		}
		// Check that if the elements of the slice are allocated from the beginning, no pointer reallocation occurs
		for i, b := range bforeaddresses {
			a := afteraddresses[i]
			assert.Equal(t, b, a)
		}
	}
}

func Test_Load(t *testing.T) {
	// normal case 1 (with header)
	{
		csv := `ID,説明,誕生日
1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 1, 21, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, float64(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// normal case 2 (without header)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float32
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 0, 20, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, float32(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// illegal case 1 (collapse format)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 0, 20, &entries)

		assert.EqualError(t, err, "record on line 7: wrong number of fields")
		assert.Empty(t, entries)
	}
	// illegal case 2 (too large topmergin)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 20, 20, &entries)

		assert.EqualError(t, err, "topmergin is too large")
		assert.Empty(t, entries)
	}
	// illegal case 3 (too large rows)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float32
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 0, 19, &entries)

		assert.EqualError(t, err, "rows are too large")
		assert.Empty(t, entries)
	}
	// illegal case 4 (struct fields less than CSV fields)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id   float32
			desc string
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 0, 100, &entries)

		assert.EqualError(t, err, "number of fields in the defined structure may not match the number of fields in the CSV.")
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, float32(0), entries[i-1].id)
		}
		assert.Equal(t, "", entries[0].desc)
		assert.Equal(t, "", entries[1].desc)
		assert.Equal(t, "", entries[2].desc)
		assert.Equal(t, "", entries[3].desc)
		assert.Equal(t, "", entries[4].desc)
		assert.Equal(t, "", entries[5].desc)
		assert.Equal(t, "", entries[6].desc)
		assert.Equal(t, "", entries[7].desc)
		assert.Equal(t, "", entries[8].desc)
		assert.Equal(t, "", entries[9].desc)
		assert.Equal(t, "", entries[10].desc)
		assert.Equal(t, "", entries[11].desc)
		assert.Equal(t, "", entries[12].desc)
		assert.Equal(t, "", entries[13].desc)
		assert.Equal(t, "", entries[14].desc)
		assert.Equal(t, "", entries[15].desc)
		assert.Equal(t, "", entries[16].desc)
		assert.Equal(t, "", entries[17].desc)
		assert.Equal(t, "", entries[18].desc)
		assert.Equal(t, "", entries[19].desc)
	}
	// illegal case 5 (struct fields more than CSV fields)
	{
		csv := `1.0,あ,2011.12.12
2.0,い,2011.12.13
3.0,う,2011.12.14
4.0,え,2011.5.12
5.0,お,2011.3.22
6.0,か,2011.4.1
7.0,き,2000.12.1
8.0,く,2011.1.11
9.0,け,2011.2.10
10.0,こ,2011.3.15
11.0,さ,2011.7.21
12.0,し,2011.8.9
13.0,す,2011.10.15
14.0,せ,2011.11.30
15.0,そ,2011.9.3
16.0,た,2011.6.5
17.0,つ,2011.5.5
18.0,て,2011.4.3
19.0,と,2011.2.3
20.0,な,2011.10.3
`
		type csventry struct {
			id    float32
			desc  string
			birth time.Time
			dummy string
		}

		entries := []*csventry{}
		err := Load(strings.NewReader(csv), 0, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, float32(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
		assert.Equal(t, "", entries[0].dummy)
		assert.Equal(t, "", entries[1].dummy)
		assert.Equal(t, "", entries[2].dummy)
		assert.Equal(t, "", entries[3].dummy)
		assert.Equal(t, "", entries[4].dummy)
		assert.Equal(t, "", entries[5].dummy)
		assert.Equal(t, "", entries[6].dummy)
		assert.Equal(t, "", entries[7].dummy)
		assert.Equal(t, "", entries[8].dummy)
		assert.Equal(t, "", entries[9].dummy)
		assert.Equal(t, "", entries[10].dummy)
		assert.Equal(t, "", entries[11].dummy)
		assert.Equal(t, "", entries[12].dummy)
		assert.Equal(t, "", entries[13].dummy)
		assert.Equal(t, "", entries[14].dummy)
		assert.Equal(t, "", entries[15].dummy)
		assert.Equal(t, "", entries[16].dummy)
		assert.Equal(t, "", entries[17].dummy)
		assert.Equal(t, "", entries[18].dummy)
		assert.Equal(t, "", entries[19].dummy)
	}
	{
		// illegal case 6 (io.Reader is nil)
		entries := []*struct{}{}
		err := Load(nil, 0, 100, &entries)

		assert.EqualError(t, err, "reader is nil")
	}
}

func Test_LoadVertically(t *testing.T) {
	// normal case 1 (with top-left-header)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int8
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int8(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// normal case 2 (without top-header)
	{
		csv := `ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int16
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int16(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// normal case 3 (without top-left-header)
	{
		csv := `1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int32
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 0, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int32(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// normal case 4 (base64)
	{
		csv := `aGVhZGVyLCwsLCwsLCwsLCwsLCwsLCwsLCwKSUQsMSwyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTMsMTQsMTUsMTYsMTcsMTgsMTksMjAK6Kqs5piOLOOBgizjgYQs44GGLOOBiCzjgYos44GLLOOBjSzjgY8s44GRLOOBkyzjgZUs44GXLOOBmSzjgZss44GdLOOBnyzjgaQs44GmLOOBqCzjgaoK6KqV55Sf5pelLDIwMTEuMTIuMTIsMjAxMS4xMi4xMywyMDExLjEyLjE0LDIwMTEuNS4xMiwyMDExLjMuMjIsMjAxMS40LjEsMjAwMC4xMi4xLDIwMTEuMS4xMSwyMDExLjIuMTAsMjAxMS4zLjE1LDIwMTEuNy4yMSwyMDExLjguOSwyMDExLjEwLjE1LDIwMTEuMTEuMzAsMjAxMS45LjMsMjAxMS42LjUsMjAxMS41LjUsMjAxMS40LjMsMjAxMS4yLjMsMjAxMS4xMC4z`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(base64.NewDecoder(base64.StdEncoding, strings.NewReader(csv)), 1, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int64(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// normal case 5 (large leftmergin)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 20, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(entries))
		assert.Equal(t, int64(20), entries[0].id)
		assert.Equal(t, "な", entries[0].desc)
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[0].birth.String())
	}
	// illegal case 1 (did not pass slice pointer)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 0, 100, entries)

		assert.EqualError(t, err, "failed to obtain a reference to i (did you forget &?)")
	}
	// illegal case 2 (passed empty csv)
	{
		csv := ``
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 0, 100, &entries)

		assert.EqualError(t, err, "EOF")
	}
	// illegal case 2-1 (passed too short csv)
	{
		csv := `,,,`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 0, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 4, len(entries))
	}
	// illegal case 3 (passed too large csv)
	{
		csv := `,,,,,,
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 0, 0, 5, &entries)

		assert.EqualError(t, err, "columns are too large")
	}
	// illegal case 4 (collapse format)
	// The second line is missing one ",".
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 1, 100, &entries)

		assert.EqualError(t, err, "record on line 3: wrong number of fields")

		// "id" fields are filled.
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int64(i), entries[i-1].id)
		}
		// "desc" fields are empty because failed 2nd csv row read.
		assert.Equal(t, "", entries[0].desc)
		assert.Equal(t, "", entries[1].desc)
		assert.Equal(t, "", entries[2].desc)
		assert.Equal(t, "", entries[3].desc)
		assert.Equal(t, "", entries[4].desc)
		assert.Equal(t, "", entries[5].desc)
		assert.Equal(t, "", entries[6].desc)
		assert.Equal(t, "", entries[7].desc)
		assert.Equal(t, "", entries[8].desc)
		assert.Equal(t, "", entries[9].desc)
		assert.Equal(t, "", entries[10].desc)
		assert.Equal(t, "", entries[11].desc)
		assert.Equal(t, "", entries[12].desc)
		assert.Equal(t, "", entries[13].desc)
		assert.Equal(t, "", entries[14].desc)
		assert.Equal(t, "", entries[15].desc)
		assert.Equal(t, "", entries[16].desc)
		assert.Equal(t, "", entries[17].desc)
		assert.Equal(t, "", entries[18].desc)
		assert.Equal(t, "", entries[19].desc)
		// "birth" fields are empty because failed 2nd csv row read.
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// illegal case 5 (too large topmergin)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 100, 1, 100, &entries)

		assert.EqualError(t, err, "EOF")
	}
	// illegal case 6 (topmergin is incorrectly specified)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 3, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for _, e := range entries {
			assert.Zero(t, e.id)
		}
		assert.Equal(t, "", entries[0].desc)
		assert.Equal(t, "", entries[1].desc)
		assert.Equal(t, "", entries[2].desc)
		assert.Equal(t, "", entries[3].desc)
		assert.Equal(t, "", entries[4].desc)
		assert.Equal(t, "", entries[5].desc)
		assert.Equal(t, "", entries[6].desc)
		assert.Equal(t, "", entries[7].desc)
		assert.Equal(t, "", entries[8].desc)
		assert.Equal(t, "", entries[9].desc)
		assert.Equal(t, "", entries[10].desc)
		assert.Equal(t, "", entries[11].desc)
		assert.Equal(t, "", entries[12].desc)
		assert.Equal(t, "", entries[13].desc)
		assert.Equal(t, "", entries[14].desc)
		assert.Equal(t, "", entries[15].desc)
		assert.Equal(t, "", entries[16].desc)
		assert.Equal(t, "", entries[17].desc)
		assert.Equal(t, "", entries[18].desc)
		assert.Equal(t, "", entries[19].desc)
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", entries[19].birth.String())
	}
	// illegal case 7 (too large leftmergin)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 10000, 21, &entries)

		assert.EqualError(t, err, "leftmergin is too large")
	}
	// illegal case 8 (struct fields less than CSV fields)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id   int64
			desc string
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int64(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
	}
	// illegal case 9 (struct fields more than CSV fields)
	{
		csv := `header,,,,,,,,,,,,,,,,,,,,
ID,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
説明,あ,い,う,え,お,か,き,く,け,こ,さ,し,す,せ,そ,た,つ,て,と,な
誕生日,2011.12.12,2011.12.13,2011.12.14,2011.5.12,2011.3.22,2011.4.1,2000.12.1,2011.1.11,2011.2.10,2011.3.15,2011.7.21,2011.8.9,2011.10.15,2011.11.30,2011.9.3,2011.6.5,2011.5.5,2011.4.3,2011.2.3,2011.10.3
`
		type csventry struct {
			id    int64
			desc  string
			birth time.Time
			dummy string
		}

		entries := []*csventry{}
		err := LoadVertically(strings.NewReader(csv), 1, 1, 100, &entries)

		assert.NoError(t, err)
		assert.Equal(t, 20, len(entries))
		for i := 1; i <= 20; i++ {
			assert.Equal(t, int64(i), entries[i-1].id)
		}
		assert.Equal(t, "あ", entries[0].desc)
		assert.Equal(t, "い", entries[1].desc)
		assert.Equal(t, "う", entries[2].desc)
		assert.Equal(t, "え", entries[3].desc)
		assert.Equal(t, "お", entries[4].desc)
		assert.Equal(t, "か", entries[5].desc)
		assert.Equal(t, "き", entries[6].desc)
		assert.Equal(t, "く", entries[7].desc)
		assert.Equal(t, "け", entries[8].desc)
		assert.Equal(t, "こ", entries[9].desc)
		assert.Equal(t, "さ", entries[10].desc)
		assert.Equal(t, "し", entries[11].desc)
		assert.Equal(t, "す", entries[12].desc)
		assert.Equal(t, "せ", entries[13].desc)
		assert.Equal(t, "そ", entries[14].desc)
		assert.Equal(t, "た", entries[15].desc)
		assert.Equal(t, "つ", entries[16].desc)
		assert.Equal(t, "て", entries[17].desc)
		assert.Equal(t, "と", entries[18].desc)
		assert.Equal(t, "な", entries[19].desc)
		assert.Equal(t, "2011-12-12 00:00:00 +0000 UTC", entries[0].birth.String())
		assert.Equal(t, "2011-12-13 00:00:00 +0000 UTC", entries[1].birth.String())
		assert.Equal(t, "2011-12-14 00:00:00 +0000 UTC", entries[2].birth.String())
		assert.Equal(t, "2011-05-12 00:00:00 +0000 UTC", entries[3].birth.String())
		assert.Equal(t, "2011-03-22 00:00:00 +0000 UTC", entries[4].birth.String())
		assert.Equal(t, "2011-04-01 00:00:00 +0000 UTC", entries[5].birth.String())
		assert.Equal(t, "2000-12-01 00:00:00 +0000 UTC", entries[6].birth.String())
		assert.Equal(t, "2011-01-11 00:00:00 +0000 UTC", entries[7].birth.String())
		assert.Equal(t, "2011-02-10 00:00:00 +0000 UTC", entries[8].birth.String())
		assert.Equal(t, "2011-03-15 00:00:00 +0000 UTC", entries[9].birth.String())
		assert.Equal(t, "2011-07-21 00:00:00 +0000 UTC", entries[10].birth.String())
		assert.Equal(t, "2011-08-09 00:00:00 +0000 UTC", entries[11].birth.String())
		assert.Equal(t, "2011-10-15 00:00:00 +0000 UTC", entries[12].birth.String())
		assert.Equal(t, "2011-11-30 00:00:00 +0000 UTC", entries[13].birth.String())
		assert.Equal(t, "2011-09-03 00:00:00 +0000 UTC", entries[14].birth.String())
		assert.Equal(t, "2011-06-05 00:00:00 +0000 UTC", entries[15].birth.String())
		assert.Equal(t, "2011-05-05 00:00:00 +0000 UTC", entries[16].birth.String())
		assert.Equal(t, "2011-04-03 00:00:00 +0000 UTC", entries[17].birth.String())
		assert.Equal(t, "2011-02-03 00:00:00 +0000 UTC", entries[18].birth.String())
		assert.Equal(t, "2011-10-03 00:00:00 +0000 UTC", entries[19].birth.String())
		assert.Equal(t, "", entries[0].dummy)
		assert.Equal(t, "", entries[1].dummy)
		assert.Equal(t, "", entries[2].dummy)
		assert.Equal(t, "", entries[3].dummy)
		assert.Equal(t, "", entries[4].dummy)
		assert.Equal(t, "", entries[5].dummy)
		assert.Equal(t, "", entries[6].dummy)
		assert.Equal(t, "", entries[7].dummy)
		assert.Equal(t, "", entries[8].dummy)
		assert.Equal(t, "", entries[9].dummy)
		assert.Equal(t, "", entries[10].dummy)
		assert.Equal(t, "", entries[11].dummy)
		assert.Equal(t, "", entries[12].dummy)
		assert.Equal(t, "", entries[13].dummy)
		assert.Equal(t, "", entries[14].dummy)
		assert.Equal(t, "", entries[15].dummy)
		assert.Equal(t, "", entries[16].dummy)
		assert.Equal(t, "", entries[17].dummy)
		assert.Equal(t, "", entries[18].dummy)
		assert.Equal(t, "", entries[19].dummy)
	}
	{
		// illegal case 10 (io.Reader is nil)
		entries := []*struct{}{}
		err := Load(nil, 0, 100, &entries)

		assert.EqualError(t, err, "reader is nil")
	}
}
