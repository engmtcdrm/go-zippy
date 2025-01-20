package main

import "example.com/m/zippy"

func main() {
	err := zippy.Unzip("test/test.zip", "test2")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test2", "test3/test.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test/test1", "test3/testfile.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test/test1", "test3/testfile2.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test", "test3/testdir.zip")
	if err != nil {
		panic(err)
	}
}
