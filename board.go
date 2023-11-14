package main

const (
	SIZE = 8
	IBB  = 0x0000000810000000
	IWB  = 0x0000001008000000
)

type Board struct {
	pb, ob uint
}

func NewBoard() *Board {
	return &Board{IBB, IWB}
}

func (b *Board) Blank() uint {
	return ^(b.pb | b.ob)
}

func (b *Board) checkRLine(mask, n uint) uint {
	res := mask & (b.pb >> n)
	for i := 0; i < 5; i++ {
		res |= mask & (res >> n)
	}
	return res >> n
}

func (b *Board) checkLLine(mask, n uint) uint {
	res := mask & (b.pb << n)
	for i := 0; i < 5; i++ {
		res |= mask & (res << n)
	}

	return res << n
}

func (b *Board) LegalBoard() uint {
	lrmask := 0x7e7e7e7e7e7e7e7e & b.ob //　左右の番人
	udmask := 0x00FFFFFFFFFFFF00 & b.ob // 上下の番人
	mask := 0x007e7e7e7e7e7e00 & b.ob   // 全辺の番人

	legal := b.checkRLine(lrmask, 1) // 右に連続する場所を求める
	legal |= b.checkLLine(lrmask, 1) // 左に連続する場所を求める
	legal |= b.checkRLine(udmask, 8) // 下に連続する場所を求める
	legal |= b.checkLLine(udmask, 8) // 上に連続する場所を求める
	legal |= b.checkRLine(mask, 7)   // 左下に連続する場所を求める
	legal |= b.checkLLine(mask, 9)   // 左上に連続する場所を求める
	legal |= b.checkRLine(mask, 9)   // 右下に連続する場所を求める
	legal |= b.checkLLine(mask, 7)   // 右上に連続する場所を求める

	return legal & b.Blank() // 空きますとの積を取る
}

func transfer(pos uint, k int) uint {
	switch k {
	case 0: //上
		return (pos << 8) & 0xffffffffffffff00
	case 1: //右上
		return (pos << 7) & 0x7f7f7f7f7f7f7f00
	case 2: //右
		return (pos >> 1) & 0x7f7f7f7f7f7f7f7f
	case 3: //右下
		return (pos >> 9) & 0x007f7f7f7f7f7f7f
	case 4: //下
		return (pos >> 8) & 0x00ffffffffffffff
	case 5: //左下
		return (pos >> 7) & 0x00fefefefefefefe
	case 6: //左
		return (pos << 1) & 0xfefefefefefefefe
	case 7: //左上
		return (pos << 9) & 0xfefefefefefefe00
	default:
		return 0
	}
}

func (b *Board) Reverse(pos uint) uint {
	var rev uint
	for i := 0; i < 8; i++ {
		var rev_ uint
		mask := transfer(pos, i)
		for mask != 0 && mask&b.ob != 0 {
			rev_ |= mask
			mask = transfer(mask, i)
		}
		if mask&b.pb != 0 {
			rev |= rev_
		}
	}
	return rev
}

func (b *Board) Put(pos, rev uint) {
	b.pb ^= rev | pos
	b.ob ^= rev
}
