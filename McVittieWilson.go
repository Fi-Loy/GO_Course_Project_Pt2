package main

import (
	"container/heap"
	"math"
	"sync"
)

func McVittieWilson(r map[int]*Resident, p map[string]*Program, concurrentFlag bool) {

	//initialize heaps
	for pid := range p {
		p[pid].selectedResidents = &residentHeap{}
		heap.Init(p[pid].selectedResidents)
	}

	var wg sync.WaitGroup

	for residentID := range r {
		wg.Add(1)
		if concurrentFlag {
			go offer(residentID, r, p, concurrentFlag, &wg)
		} else {
			offer(residentID, r, p, concurrentFlag, &wg)
		}
	}

	wg.Wait()
}

func offer(
	rid int,
	residents map[int]*Resident,
	programs map[string]*Program,
	concurrentFlag bool,
	wg *sync.WaitGroup) {

	defer wg.Done()

	resident := residents[rid]
	resident.lock.Lock()

	if len(resident.rol) == 0 {
		resident.matchedProgram = ""
	} else {
		pid := resident.rol[0]
		resident.rol = resident.rol[1:]
		resident.lock.Unlock()
		evaluate(rid, pid, residents, programs, concurrentFlag, wg)
	}
}

func evaluate(rid int, pid string,
	residents map[int]*Resident,
	programs map[string]*Program,
	concurrentFlag bool,
	wg *sync.WaitGroup) {

	isIn := false
	rRank := math.MaxInt

	p := programs[pid]
	p.lock.Lock()
	r := residents[rid]
	r.lock.Lock()

	for i := 0; i < len(p.rol); i++ {
		if rid == p.rol[i] {
			isIn = true
			rRank = i
		}
	}

	if !isIn {
		p.lock.Unlock()
		r.lock.Unlock()
		wg.Add(1)
		if concurrentFlag {
			go offer(rid, residents, programs, concurrentFlag, wg)
		} else {
			offer(rid, residents, programs, concurrentFlag, wg)
		}
	} else if p.nPositions > 0 {
		r.matchedProgram = pid
		heap.Push(p.selectedResidents, resRank{rRank, rid})
		p.nPositions--
		p.lock.Unlock()
		r.lock.Unlock()
	} else if p.selectedResidents.Peek().(resRank).rank > rRank {
		ejectedResID := heap.Pop(p.selectedResidents).(resRank).rid
		heap.Push(p.selectedResidents, resRank{rRank, rid})
		r.matchedProgram = pid
		p.lock.Unlock()
		r.lock.Unlock()
		wg.Add(1)
		if concurrentFlag {
			go offer(ejectedResID, residents, programs, concurrentFlag, wg)
		} else {
			offer(ejectedResID, residents, programs, concurrentFlag, wg)
		}
	} else {
		p.lock.Unlock()
		r.lock.Unlock()
		wg.Add(1)
		if concurrentFlag {
			go offer(rid, residents, programs, concurrentFlag, wg)
		} else {
			offer(rid, residents, programs, concurrentFlag, wg)
		}
	}
}
