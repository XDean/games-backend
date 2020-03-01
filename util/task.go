package util

func FirstError(errors ...error) error {
	for _, e := range errors {
		if e != nil {
			return e
		}
	}
	return nil
}

func DoUntilError(tasks ...func() error) error {
	for _, t := range tasks {
		err := t()
		if err != nil {
			return err
		}
	}
	return nil
}
