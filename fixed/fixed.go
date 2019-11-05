package fixed

type Int20_12 int32

func From(i int) Int20_12      { return Int20_12(i << 12) }
func FromI32(i int32) Int20_12 { return Int20_12(i << 12) }
func FromI64(i int64) Int20_12 { return Int20_12(i << 12) }
func FromU8(i uint8) Int20_12  { return Int20_12(i) << 12 }

func (x Int20_12) ToI32() int32 { return int32(x >> 12) }
func (x Int20_12) ToU8() uint8  { return uint8(x >> 12) }

func (x Int20_12) Add(y Int20_12) Int20_12 { return x + y }
func (x Int20_12) Sub(y Int20_12) Int20_12 { return x - y }
func (x Int20_12) Mul(y Int20_12) Int20_12 { return Int20_12((x * y) >> 12) }
func (x Int20_12) Div(y Int20_12) Int20_12 { return Int20_12((x << 12) / y) }

func (x Int20_12) AddI(y int) Int20_12 { return x.Add(From(y)) }
func (x Int20_12) DivI(y int) Int20_12 { return x.Div(From(y)) }
func (x Int20_12) SubI(y int) Int20_12 { return x.Sub(From(y)) }
func (x Int20_12) MulI(y int) Int20_12 { return x.Mul(From(y)) }
