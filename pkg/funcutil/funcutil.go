package funcutil



func continueIfNoError(f func() (interface{}, error)) (res interface{}, err error) {
	if res, err = f(); err != nil {
		return res, err
	}
	return nil, err
}



