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
	fmt.Println(ad.Dir)
	//RunCommand("bash -c eval `ssh-agent`", &ad, ad.Dir, []string{}, "")
	//RunCommand("bash -c ssh-add", &ad, ad.Dir, []string{}, "")
	repoURL := "https://github.com/citizensciencecenter/" + n
	fmt.Println(repoURL)
	r, err := git.PlainClone(ad.Dir, false, &git.CloneOptions{
		URL:               "https://github.com/citizensciencecenter/" + n,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	//ErrHandler(err)
	fmt.Println(n)
	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen(ad.Dir)
		fmt.Println("Repo opened")
	} else {
		ErrHandler(err)
		fmt.Println("Repo checked out")
	}
	// TODO allow this to be configured?
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	//ErrHandler(err)
	fmt.Println("Fetched remotes")
	branches, err := r.References()
	ErrHandler(err)
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
	RunCommand("git submodule init", &ad, ad.Dir, []string{}, "Submodules initialised")
	/*s.Init()
	s.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		NoFetch:           false,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})*/
	fmt.Println("Updated submodules")
	err = w.Pull(&git.PullOptions{RemoteName: "origin", RecurseSubmodules: git.DefaultSubmoduleRecursionDepth})
	RunCommand("git submodule update --recursive", &ad, ad.Dir, []string{}, "Submodules updated")
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
