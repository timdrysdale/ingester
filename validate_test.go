package ingester

/*
func TestValidate(t *testing.T) {

	gradexpath.SetTesting()

	root := gradexpath.Root()

	if root != "./tmp-delete-me" {
		t.Errorf("test root set up wrong %s", root)
	}

	// don't use GetRoot() here
	// JUST in case we kill a whole working installation
	os.RemoveAll("./tmp-delete-me")

	EnsureDirectoryStructure()

	testfiles, err := gradexpath.GetFileList("./test")

	assert.NoError(t, err)

	//fmt.Println(testfiles)

	for _, file := range testfiles {
		destination := filepath.Join(gradexpath.Ingest(), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)

	}

	//fmt.Println(gradexpath.Ingest())
	ingestfiles, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	assert.True(t, gradexpath.CopyIsComplete(testfiles, ingestfiles))

	StageFromIngest()

	expectedRejects, err := gradexpath.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := gradexpath.GetFileList(gradexpath.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, gradexpath.CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := gradexpath.GetFileList("./expected/temp-txt")
	assert.NoError(t, err)

	actualTxt, err := gradexpath.GetFileList(gradexpath.TempTxt())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gradexpath.CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := gradexpath.GetFileList("./expected/temp-pdf")
	assert.NoError(t, err)

	actualPdf, err := gradexpath.GetFileList(gradexpath.TempPdf())
	assert.NoError(t, err)

	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gradexpath.CopyIsComplete(expectedPdf, actualPdf))

	// now we test validate!

	assert.NoError(t, ValidateNewPapers())

	exam := "Practice Exam Drop Box"

	actualPdf, err = gradexpath.GetFileList(gradexpath.AcceptedPapers(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedPdf) == len(actualPdf))
	assert.True(t, gradexpath.CopyIsComplete(expectedPdf, actualPdf))

	actualTxt, err = gradexpath.GetFileList(gradexpath.AcceptedReceipts(exam))
	assert.NoError(t, err)
	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, gradexpath.CopyIsComplete(expectedTxt, actualTxt))

	tempPdf, err := gradexpath.GetFileList(gradexpath.TempPdf())
	assert.NoError(t, err)
	assert.Equal(t, len(tempPdf), 0)

	tempTxt, err := gradexpath.GetFileList(gradexpath.TempTxt())
	assert.NoError(t, err)
	assert.Equal(t, len(tempTxt), 0)

}
*/
