// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package profile

import (
	"fmt"
	"sync"
	"time"
)

type TimeRecord struct {
	count    int
	duration time.Duration
}

type Profiler struct {
	profMap map[string]*TimeRecord

	// TODO: Use simple synchronization mechanism for thread safety.
	// TODO: Need improvement
	mutex sync.Mutex
}

func (p *Profiler) Profile(key string, d time.Duration) {
	p.mutex.Lock()

	if r, ok := p.profMap[key]; ok {
		r.count += 1
		r.duration += d
	} else {
		p.profMap[key] = &TimeRecord{count: 1, duration: d}
	}

	p.mutex.Unlock()
}

func (p *Profiler) PrintProfileInfo() {
	fmt.Printf("%s", p.GetProfileInfoString())
}

func (p *Profiler) GetProfileInfoString() string {
	str := fmt.Sprintf("Key,Count,Time\n")
	for k, v := range p.profMap {
		str += fmt.Sprintf("%s,%d,%f\n", k, v.count, v.duration.Seconds())
	}

	return str
}

func NewProfiler() *Profiler {
	return &Profiler{
		profMap: make(map[string]*TimeRecord),
	}
}

var Prof = NewProfiler()
