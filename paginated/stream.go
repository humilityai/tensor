package paginated

// Stream is an interface for streaming data
// to a paginated tensor.
type Stream struct {
	*Tensor
	data     chan Array
	finalize chan bool
}

// Stream can be used to more easily write a large amount of data
// to an initialized paginated tensor.
func (p *Tensor) Stream(mutable []bool) (*Stream, error) {
	if len(mutable) != p.Dims() {
		return &Stream{}, ErrSize
	}

	dataChan := make(chan Array, 1)
	finalizeChan := make(chan bool, 1)

	go func(p *Tensor) {
		select {
		case data := <-dataChan:
			p.addData(data)
		case <-finalizeChan:
			close(dataChan)
			close(finalizeChan)
			break
		}
	}(p)

	return &Stream{
		Tensor:   p,
		data:     dataChan,
		finalize: finalizeChan,
	}, nil
}

// Send will block until the Stream is ready
// to accept and write new data to the paginated tensor.
func (s *Stream) Send(data Array) {
	s.data <- data
}

// Finalize will close all channels, thus killing the Stream
func (s *Stream) Finalize() {
	s.finalize <- true
}

func (p *Tensor) addData(data Array) error {
	page, ok := p.firstUnfilledPage()
	if !ok {
		return ErrFilled
	}

	var pageData Array
	d, ok := p.cache.Get(page.Id)
	if !ok {
		err := p.Swap(page.Id)
		if err != nil {
			return err
		}

		d, ok = p.cache.Get(page.Id)
		if !ok {
			return ErrCache
		}
	}
	pageData = d.(Array)

	// remaining page-capacity
	if page.Datasize > 0 {
		var err error
		pageData, err = pageData.Slice(page.Datasize-1, pageData.Len()-1)
		if err != nil {
			return err
		}
	}

	// slice of new data
	capacity := pageData.Len()
	if data.Len() < capacity {
		capacity = data.Len()
	}
	slice, err := data.Slice(0, capacity)
	if err != nil {
		return err
	}

	// copy data
	err = pageData.Copy(slice)
	if err != nil {
		return err
	}
	page.Datasize += slice.Len()

	// recurse: for all data in supplied array
	if slice.Len() < data.Len() {
		remaining, err := data.Slice(slice.Len(), data.Len()-1)
		if err != nil {
			return err
		}

		return p.addData(remaining)
	}

	return nil
}