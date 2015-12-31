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
	if m.Pool != 80 || g.Amount != 30 {
		t.Fatalf("pool:%v amount:%v", m.Pool, g.Amount)
	}
	m.NewHand(f.Win(1), c.Win(1), r.Win(1))
	if m.Pool != 80 || g.Amount != 30 {
		t.Fatalf("pool:%v amount:%v", m.Pool, g.Amount)
	}
	m.NewHand(g.Win(3), c.Win(0), r.Win(0), f.Win(0))
	if m.Pool != 240 || g.Amount != 110 {
		t.Fatalf("pool:%v amount:%v", m.Pool, g.Amount)
	}
	m.NewHand(g.Win(0), r.Win(3))
	if m.Pool != 240 || g.Amount != -130 || r.Amount != 110 {
		t.Fatalf("pool:%v g-amount:%v r-amount:%v", m.Pool, g.Amount, r.Amount)
	}
	m.NewHand(g.Win(1), r.Win(2), c.Win(0))
	if m.Pool != 240 || g.Amount != -50 || r.Amount != 270 {
		t.Fatalf("pool:%v g-amount:%v r-amount:%v", m.Pool, g.Amount, r.Amount)
	}
}
