package main

import "example.com/m/zippy"

func main() {
	err := zippy.Unzip("test/test.zip", "test2")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test2", "test3/test2-bak.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test/test1", "test3/testfile1.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test/test2", "test3/testfile2.zip")
	if err != nil {
		panic(err)
	}

	err = zippy.Zip("test", "test3/test-dir.zip")
	if err != nil {
		panic(err)
	}
}
