package inject

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImpl_Get(t *testing.T) {
	ctx := impl{}
	assert.NoError(t, ctx.Register(&BeanA{1}))
	assert.NoError(t, ctx.Register(&BeanB{2}))
	assert.NoError(t, ctx.Register(&BeanC{3}))

	var ra = BeanA{}
	var rb = BeanB{}
	var rc = BeanC{}
	var rd = BeanD{}
	var rfa InterfaceA
	var rfb InterfaceB
	assert.NoError(t, ctx.Get(&ra))
	assert.NoError(t, ctx.Get(&rb))
	assert.NoError(t, ctx.Get(&rc))
	assert.NoError(t, ctx.Get(&rfa))
	assert.NoError(t, ctx.Get(&rfb))
	assert.Equal(t, 1, ra.a)
	assert.Equal(t, 2, rb.b)
	assert.Equal(t, 3, rc.c)
	assert.Equal(t, 1, rfa.fa())
	assert.Equal(t, 2, rfb.fb())

	assert.NoError(t, ctx.Inject(&rd))
	assert.Equal(t, 1, rd.A.a)
	assert.Equal(t, 2, rd.B.b)
	assert.Nil(t, rd.C)
	assert.Equal(t, 1, rd.Fa.fa())
	assert.Equal(t, 2, rd.Fb.fb())
	assert.Equal(t, 3, rd.Fc.fa())
	assert.Equal(t, 3, rd.Fc.fb())
}

type (
	BeanA struct {
		a int
	}

	BeanB struct {
		b int
	}

	BeanC struct {
		c int
	}

	BeanD struct {
		A  BeanA  `inject:""`
		B  *BeanB `inject:""`
		C  *BeanC
		Fa InterfaceA `inject:""`
		Fb InterfaceB `inject:""`
		Fc interface {
			InterfaceA
			InterfaceB
		} `inject:""`
	}

	InterfaceA interface {
		fa() int
	}

	InterfaceB interface {
		fb() int
	}
)

func (a BeanA) fa() int {
	return a.a
}

func (a *BeanB) fb() int {
	return a.b
}

func (a BeanC) fa() int {
	return a.c
}

func (a BeanC) fb() int {
	return a.c
}
