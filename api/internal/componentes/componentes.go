package componentes

import (
	"quepc/api/internal/componentes/discos"
	"quepc/api/internal/componentes/gabinetes"
	"quepc/api/internal/componentes/gpus"
	"quepc/api/internal/componentes/motherboards"
	"quepc/api/internal/componentes/procesadores"
	"quepc/api/internal/componentes/psus"
	"quepc/api/internal/componentes/rams"
)

type Componentes struct {
	Discos       *discos.Discos
	Gabinetes    *gabinetes.Gabinetes
	GPUs         *gpus.GPUs
	Motherboards *motherboards.Motherboards
	Procesadores *procesadores.Procesadores
	PSUs         *psus.PSUs
	RAMs         *rams.RAMs
}

func New(
	ds *discos.Discos,
	gs *gabinetes.Gabinetes,
	gp *gpus.GPUs,
	mbs *motherboards.Motherboards,
	prs *procesadores.Procesadores,
	ps *psus.PSUs,
	rms *rams.RAMs) *Componentes {
	return &Componentes{
		Discos:       ds,
		Gabinetes:    gs,
		GPUs:         gp,
		Motherboards: mbs,
		Procesadores: prs,
		PSUs:         ps,
		RAMs:         rms,
	}
}
