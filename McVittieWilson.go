package main

import (
	"math"
)

func McVittieWilson(r map[int]*Resident, p map[int]*Program) {
	for resident := range r {
		offer(r[resident])
	}
}

func offer(r *Resident) {

}

func evaluate(r *Resident, p *Program) {

	isIn := false

	rRank := math.MaxInt

	for i := 0; i < len(p.rol); i++ {
		if r.residentID == p.rol[i] {
			isIn = true
			rRank = i
		}
	}

	if !isIn {
		offer(r)
	} else if p.nPositions > 0 {
		r.matchedProgram = p.programID
		p.selectedResidents.Push(resRank{rRank, r})
	} else if p.selectedResidents.Peek().(resRank).rank > rRank {
		ejectedRes := p.selectedResidents.Pop().(resRank).res
		p.selectedResidents.Push(resRank{rRank, r})
		r.matchedProgram = p.programID
		offer(ejectedRes)
	} else {
		offer(r)
	}
}
