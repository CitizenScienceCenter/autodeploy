package modules

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// InitRepo clones or updates a repo based on the branch info coming from Travis
func InitRepo(n string, b string, ad AutoDeploy) {
	r, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:               "https://github.com/citizensciencecenter/" + n,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	fmt.Println(err)
	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen("/tmp/foo")
		fmt.Println("Repo opened")
	} else {
		fmt.Printf("Repo checked out")
	}
	// TODO allow this to be configured?
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	fmt.Println("Fetched remotes")
	branches, _ := r.References()
	fmt.Println("Searching references")
	var target plumbing.ReferenceName
	for {
		v, err := branches.Next()
		ErrHandler(err)
		if strings.Contains(v.Name().String(), b) {
			target = v.Name()
			fmt.Println(target)
			break
		}
	}
	fmt.Println("Found branch")
	w, err := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{
		Branch: target,
		Force:  true,
	})
	//s, err := w.Submodules()
	RunCommand("git submodule init", &ad, "/tmp/foo", []string{}, "Submodules initialised")
	RunCommand("git submodule update --remote --recursive", &ad, "/tmp/foo", []string{}, "Submodules updated")
	/*s.Init()
	s.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		NoFetch:           false,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})*/
	fmt.Println("Updated submodules")
	err = w.Pull(&git.PullOptions{RemoteName: "origin", RecurseSubmodules: git.DefaultSubmoduleRecursionDepth})
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	fmt.Println(commit)
	hash := ref.Hash().String()
	// dockerURL := viper.GetString("docker.registry")
	branchFmt := strings.ReplaceAll(b, "/", "_")
	ad.Hash = hash
	tag := fmt.Sprintf("%s:%s%s", n, branchFmt, hash)
	ad.HookBody.Stage = "Git Pull"
	ad.HookBody.Status = "SUCCESS"
	Notify(ad)
	go dockerBuild(tag, ad)
}
