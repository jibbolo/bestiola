package main

import "testing"

func TestGame(t *testing.T) {
	g, _ := NewUser("Giacomo")
	r, _ := NewUser("Roberto")
	c, _ := NewUser("Carlo")
	f, _ := NewUser("Fabio")
	m, _ := NewMatch(g, r, c, f)
	if m.Pool != 0 {
		t.Fail()
	}
	m.NewHand(g.Win(1), c.Win(1), r.Win(1))
	if m.Pool != 40 {
		t.Fail()
	}
	m.NewHand(g.Win(3), c.Win(0), r.Win(0))
	if m.Pool != 80 || g.Amount(m) != 30 {
		t.Fatalf("pool:%v(80) amount:%v(30)", m.Pool, g.Amount(m))
	}
	m.NewHand(f.Win(1), c.Win(1), r.Win(1))
	if m.Pool != 80 || g.Amount(m) != 30 {
		t.Fatalf("pool:%v(80) amount:%v(30)", m.Pool, g.Amount(m))
	}
	m.NewHand(g.Win(3), c.Win(0), r.Win(0), f.Win(0))
	if m.Pool != 240 || g.Amount(m) != 110 {
		t.Fatalf("pool:%v(240) amount:%(110)v", m.Pool, g.Amount(m))
	}
	m.NewHand(g.Win(0), r.Win(3))
	if m.Pool != 240 || g.Amount(m) != -130 || r.Amount(m) != 110 {
		t.Fatalf("pool:%v(240) g-amount:%v(-130) r-amount:%v(110)", m.Pool, g.Amount(m), r.Amount(m))
	}
	m.NewHand(g.Win(1), r.Win(2), c.Win(0))
	if m.Pool != 240 || g.Amount(m) != -50 || r.Amount(m) != 270 {
		t.Fatalf("pool:%v(240) g-amount:%v(-50) r-amount:%v(270)", m.Pool, g.Amount(m), r.Amount(m))
	}
}
