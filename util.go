package main

func chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}
