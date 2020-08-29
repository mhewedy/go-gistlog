# go-gistlog

Uses Github Gists as a log database.

Each Gist created can have multiple files. so think in the gist as a database, and in the file inside the gist as a Table.

Gists resemble a NoSQL databases.

Gists are restricted by [rate-limiting](https://docs.github.com/en/developers/apps/rate-limits-for-github-apps), you need to be aware of that.

```golang
log := gistlog.NewLog("<PUT GIST ID HERE>", func() string {
	return "<PUT VALID GITHUB TOKEN HERE>"
})

fmt.Println("Inserting data into a new/existing file in the specified gist id")
err := log.Insert("aNewFileInTheGist", []string{"val1", "val3"})
fmt.Println(err)

fmt.Println("Inserting data asynchronously into a new/existing file in the specified gist id, " +
	"in such case errors are ignored")
log.InsertAsync("yetAnotherNewFile", []string{"val1", "val3"})

fmt.Println("Read data from gist by file name")
fmt.Println(log.Read("aNewFileInTheGist"))

time.Sleep(30 * time.Second) // wait for the async to be written
```
