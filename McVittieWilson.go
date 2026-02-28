package main

import (
	"container/heap"
	"math"
)

func McVittieWilson(r map[int]*Resident, p map[string]*Program) {

	//initialize heaps
	for pid := range p {
		p[pid].selectedResidents = &residentHeap{}
		heap.Init(p[pid].selectedResidents)
	}

	for residentID := range r {
		offer(residentID, r, p)
	}
}

func offer(
	rid int,
	residents map[int]*Resident,
	programs map[string]*Program) {

	resident := residents[rid]

	if len(resident.rol) == 0 {
		resident.matchedProgram = ""
	} else {
		pid := resident.rol[0]
		resident.rol = resident.rol[1:]
		evaluate(rid, pid, residents, programs)
	}
}

func evaluate(rid int, pid string,
	residents map[int]*Resident,
	programs map[string]*Program) {

	isIn := false
	rRank := math.MaxInt

	p := programs[pid]
	r := residents[rid]

	for i := 0; i < len(p.rol); i++ {
		if rid == p.rol[i] {
			isIn = true
			rRank = i
		}
	}

	if !isIn {
		offer(rid, residents, programs)
	} else if p.nPositions > 0 {
		r.matchedProgram = pid
		heap.Push(p.selectedResidents, resRank{rRank, rid})
		p.nPositions--
	} else if p.selectedResidents.Peek().(resRank).rank > rRank {
		ejectedResID := heap.Pop(p.selectedResidents).(resRank).rid
		heap.Push(p.selectedResidents, resRank{rRank, rid})
		r.matchedProgram = pid
		offer(ejectedResID, residents, programs)
	} else {
		offer(rid, residents, programs)
	}
}
