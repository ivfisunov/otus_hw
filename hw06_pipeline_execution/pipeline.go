package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, fun := range stages {
		newIn := make(Bi)
		prepareStage(done, newIn, out)
		out = fun(newIn)
	}
	return out
}

func prepareStage(done <-chan interface{}, newIn Bi, out Out) {
	go func() {
		defer close(newIn)
		for {
			select {
			case <-done:
				return
			case v, ok := <-out:
				if !ok {
					return
				}
				newIn <- v
			}
		}
	}()
}
