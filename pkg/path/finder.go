package path

import (
	"log"

	"github.com/andycondon/pathfinder/pkg/gps"
	"github.com/andycondon/pathfinder/pkg/ir"
	"github.com/andycondon/pathfinder/pkg/motor"
)

type Finder struct {
	IR    <-chan ir.Reading
	GPS   <-chan gps.Reading
	Drive chan<- motor.Command
	Done  chan struct{}
}

func (f *Finder) Find() {
	var (
		IR      ir.Reading
		GPS     gps.Reading
		forward = motor.Command{M: motor.Forward, S: motor.Slow}
		stop    = motor.Command{M: motor.Park}
		left    = motor.Command{M: motor.RotateLeft, S: motor.Medium}
	)
	for {
		select {
		case <-f.Done:
			return
		case GPS = <-f.GPS:
			log.Println("finder gps - " + GPS.String())
		case IR = <-f.IR:
			log.Println("finder ir - " + IR.String())
		}

		if IR.AllClear() {
			f.Drive <- forward
			continue
		}
		if !IR.F.IsNear() && (IR.L.IsFar() || IR.R.IsFar()) {
			f.Drive <- forward
			continue
		}
		if IR.F.IsNear() && !IR.L.IsNear() {
			f.Drive <- left
			continue
		}
		if IR.F.IsNear() {
			f.Drive <- stop
			continue
		}
	}
}
