// +build !riscv64

package mathutil

import "utilware/sqlite/libc/bigfft"

func (f *float) sqr() {
	f.n = bigfft.Mul(f.n, f.n)
	f.fracBits *= 2
	f.normalize()
}
