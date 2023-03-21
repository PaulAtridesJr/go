curl -X "POST" "http://localhost:3333/task" -H "Content-Type: application/json" -d "{\"text\":\"task first\",\"tags\":\[\"todo\", \"life\"\], \"due\":\"2016-01-02T15:04:05+00:00\"}"


serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server two: %s\n", err)
		}
		cancelCtx()
	}()


	func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

    fmt.Printf("'%s'\n", r.URL.Path)
	fmt.Printf("'%s'\n", body) // json
	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)
	//fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))

	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
}

const keyServerAddr = "serverAddr"
	//mux.HandleFunc("/aa/bb/5/hh", getRoot)
	//mux.HandleFunc("/hello", getHello)
