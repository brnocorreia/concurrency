package entity

type Player struct {
	Id     int
	Power  int
	Points int
}

func NewPlayer(id int, power int) *Player {
	return &Player{
		Id:     id,
		Power:  power,
		Points: 0,
	}
}

func (p *Player) GetDamage() int {
	return p.Power
}

func (p *Player) GetPoints() int {
	return p.Points
}

func (p *Player) AddPoint() {
	p.Points++
	// fmt.Printf("O player %d tem %d pontos\n", p.Id, p.Points)
}
