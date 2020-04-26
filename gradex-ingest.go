package ingester

/*
func somestuff() {

	// do this every time we start - it's idempotent
	ensureDirectoryStructure()

	a := app.New()

	rand.Seed(time.Now().UnixNano())

	w := a.NewWindow("Gradex-Ingest")

	if len(os.Args) < 2 {
		fmt.Printf("Requires two arguments: output_dir learn_receipt[s]\n")
		fmt.Printf("Usage: gradex-ingest ./renamed *.txt\n")
		os.Exit(0)
	}

	outputDir := os.Args[1]
	err := ensureDir(outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var inputPath []string

	inputPath = os.Args[2:]

	fmt.Printf("Receipts: %d\n", len(inputPath))

	suffix := filepath.Ext(inputPath[0])

	// sanity check
	if strings.ToLower(suffix) != ".txt" {
		fmt.Printf("Error: input path must be a .txt\n")
		os.Exit(1)
	}

	err = ensureDir(outputDir)
	if err != nil {
		fmt.Printf("Can't create output dir %s\n", outputDir)
		os.Exit(1)
	}

	var submissions []parselearn.Submission

	// ingest report
	for _, receipt := range inputPath {

		submission, err := parselearn.ParseLearnReceipt(receipt)
		if err == nil {
			submissions = append(submissions, submission)
		} else {
			fmt.Printf("Error with %s: %v\n", receipt, err)
		}

	}

	// TODO cross platform file join
	ingestReportPath := fmt.Sprintf("%s/ingest-report-%d.csv", outputDir, time.Now().UnixNano())
	parselearn.WriteSubmissionsToCSV(submissions, ingestReportPath)

	destinations := make(map[string]bool)

	// copy file to renamed dir, and rename it

	for _, sub := range submissions {

		//pretend everything has been coverted to pdf
		ext := filepath.Ext(sub.Filename)
		if ext != ".pdf" {
			filename := strings.TrimSuffix(sub.Filename, ext)
			sub.Filename = filename + ".pdf"
		}

		extOriginal := filepath.Ext(sub.OriginalFilename)
		if extOriginal != ".pdf" {
			filename := strings.TrimSuffix(sub.OriginalFilename, ext)
			sub.OriginalFilename = filename + ".pdf"
		}

		source := sub.Filename
		destination := fmt.Sprintf("%s/%s", outputDir, sub.OriginalFilename)
		if _, ok := destinations[destination]; !ok {
			destinations[destination] = true

		} else {
			destination = fmt.Sprintf("%s/%s", outputDir, sub.Filename)
			if _, ok := destinations[destination]; !ok {
				destinations[destination] = true
			} else {
				fmt.Printf("collision:%s-X-%s\n", sub.Filename, destination)
			}
		}
		fmt.Printf("%s->%s\n", source, destination)
		err = copy(source, destination, 32768)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}
	}

}
*/
