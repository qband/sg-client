package main

import (
  "src.sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
  "fmt"
  "net/url"
  "gopkg.in/inconshreveable/log15.v2"
  "log"
  "sourcegraph.com/sqs/pbtypes"
  "os"
  "golang.org/x/oauth2"
  "golang.org/x/net/context"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
)

func main() {
  endpointURL := &url.URL{Scheme: "http", Host: "xxxx:8082", Path: "/"}
  ctx := context.Background()
  key := "xxxx"
  repoURI := os.Args[1]
  if repoURI == "" {
    fmt.Println("repoURI is empty")
  }
  cloneURI := os.Args[2]
  if cloneURI == "" {
    fmt.Println("cloneURI is empty")
  }

  // set endpoint url
  ctx = sourcegraph.WithGRPCEndpoint(ctx, endpointURL)
  //ctx = fed.NewRemoteContext(ctx, endpointURL)
  fmt.Println("========================", "endpoint:", sourcegraph.GRPCEndpoint(ctx))

  // set credential
  ctx = sourcegraph.WithCredentials(ctx, oauth2.StaticTokenSource(&oauth2.Token{
    TokenType: "Bearer",
    AccessToken: key}),
  )
  ctx = metadata.NewContext(ctx, metadata.MD{"want-access-token": []string{"Bearer " + key}})

  // get api client
  cl, err := sourcegraph.NewClientFromContext(ctx)
  if err != nil {
    log15.Error("Failed to verify saved auth credentials for %s", "endpointURL", endpointURL, "error", err)
    os.Exit(-1)
  }
  _, err = cl.Auth.Identify(ctx, &pbtypes.Void{})
  if err != nil {
    log.Printf("# Failed to verify saved auth credentials for %s", endpointURL)
    os.Exit(-1)
  }
  fmt.Println("========================", "client:", cl.Conn)

  // invoke api
  // if exist then delete
  if _, err := cl.Repos.Get(ctx, &sourcegraph.RepoSpec{URI: repoURI}); grpc.Code(err) != codes.NotFound {
    switch err {
    case nil:
      log15.Warn("repo already exists", "repoURI", repoURI)
      fmt.Errorf("Repo %s already exists", repoURI)
      fmt.Errorf("Delete Repo %s first, then clone", repoURI)
      // Delete the repo.
      _, err = cl.Repos.Delete(ctx, &sourcegraph.RepoSpec{
        URI: repoURI,
      })
      if err != nil {
        fmt.Println(err)
        os.Exit(-1)
      }

    default:
      log15.Warn("problem fetching repository", "error", err)
      fmt.Errorf("Problem fetching repository: %s", err)
      os.Exit(-1)
    }
  }

  // mirror the remote git repository
  repo, err := cl.Repos.Create(ctx, &sourcegraph.ReposCreateOp{
    URI: repoURI,
    CloneURL: cloneURI,
    VCS: "git",
    Mirror: true,
  })
  if err != nil {
    log15.Error("failed to create repo", "error", err)
    os.Exit(-1)
  }
  fmt.Println("========================", "ReposCreateOp:", repo)
}
