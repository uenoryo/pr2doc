## 試し方

#### クローン

```
git clone https://github.com/uenoryo/pr2doc
cd pr2doc
```

#### コードを書く (exampleを使う)

```
vi example/example.go
```

以下を書き換えます

```go

const (
    repoOwner = "owner"
    repoName  = "repo-name"
    token     = "github-access-token"
)

```

#### ビルド&実行

```
go build -o pr2doc ./example/example.go

./pr2doc hashhashhash

```
