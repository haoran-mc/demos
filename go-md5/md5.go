package md5

import (
	"fmt"
	"math"
)

const (
	A uint32 = 0x67452301
	B uint32 = 0xEFCDAB89
	C uint32 = 0x98BADCFE
	D uint32 = 0x10325476
)

var (
	k [64]uint32
	// 循环左移的长度
	r [64]int = [64]int{7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21}
)

func F(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func G(x, y, z uint32) uint32 {
	return (x & z) | (y & ^z)
}

func H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func I(x, y, z uint32) uint32 {
	return y ^ (x | ^z)
}

func leftrotate(x uint32, n int) uint32 {
	return (x>>(32-n) | x<<n)
}

// 2 ^ 32 * abs(sin(i + 1))
func fill() {
	for i := 0; i < 64; i++ {
		k[i] = uint32(math.Abs(math.Sin(float64(i+1))) * float64(1<<32))
	}
}

// 消息填充到 512 bit
func intensify(msg []byte) []byte {
	bytes := uint64(len(msg)) // byte
	bitLen := bytes * 8       // bit

	// 1. 填充一个 1，然后填充若干个 0
	msg = append(msg, 0x80)
	bytes++

	// 2. 一直填充 0 到 512*N + 448
	for bytes%64 < 56 {
		msg = append(msg, 0)
		bytes++
	}

	// 3. 最后的 64 bit 用来放原始数据的长度，小端序（计算机电路先处理低位字节效率比较高）
	// 12345678 -> 78563412
	//           低地址 -> 高地址
	// 78 是低位，放在低地址
	for i := 0; i < 8; i++ {
		b := bitLen >> (i * 8) & 0x0000000ff
		msg = append(msg, byte(b))
	}
	return msg
}

func MainLoop(msg []byte) string {
	n := len(msg)
	var w [16]uint32
	atemp, btemp, ctemp, dtemp := A, B, C, D
	for i := 0; i < n; i += 64 {
		// 512 bit block -> n * 32 bit block
		for j := i; j < i+64; j += 4 {
			w[j/4] = uint32(msg[j]) + uint32(msg[j+1])<<8 + uint32(msg[j+2])<<16 + uint32(msg[j+3])<<24
		}
		a, b, c, d := atemp, btemp, ctemp, dtemp
		for j := 0; j < 64; j++ {
			var f uint32
			var g int
			switch {
			case j < 16:
				f, g = F(b, c, d), j
			case j < 32:
				f, g = G(b, c, d), (5*j+1)%16
			case j < 48:
				f, g = H(b, c, d), (3*j+5)%16
			default:
				f, g = I(b, c, d), (7*j)%16
			}
			temp := d
			d = c
			c = b
			b = leftrotate((a+f+k[j]+w[g]), r[j]) + b
			a = temp
		}
		atemp, btemp, ctemp, dtemp = atemp+a, btemp+b, ctemp+c, dtemp+d
	}
	// 格式化输出
	format := func(x uint32) []byte {
		res := make([]byte, 0)
		for i := 0; i < 4; i++ {
			mid := x >> (8 * i) & 0x000ff
			res = append(res, byte(mid))
		}
		return res
	}
	ra, rb, rc, rd := format(atemp), format(btemp), format(ctemp), format(dtemp)
	res := fmt.Sprintf("%x%x%x%x", ra, rb, rc, rd)
	return res
}
