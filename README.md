# go-gistlog

Golang Library to use Github Gists as a simple NoSQL database, mainly for low-traffic logging purposes.

Each Gist created can have multiple files. So think in the gist as a **Database**, and in the file inside the gist as a **Table**.

Gists API calls are restricted by [rate-limiting](https://docs.github.com/en/developers/apps/rate-limits-for-github-apps), you need to be aware of that.

```golang
log := gistlog.NewLog("<PUT GIST ID HERE>", func() string {
	return "<PUT VALID GITHUB TOKEN HERE>"
})

//Inserting data into a new/existing file in the specified gist id above
err := log.Insert("aNewFileInTheGist", []string{
	"val1",
	"val2",
	"val3",
})
fmt.Println(err)

// Inserting data asynchronously into a new/existing file in the specified gist id above, in such case errors are ignored
log.InsertAsync("yetAnotherNewFile", []string{
	"val1",
	"val2",
	"val3",
})

// Read data from gist by filename, return a slice of slices
fmt.Println(log.Read("aNewFileInTheGist"))
```
