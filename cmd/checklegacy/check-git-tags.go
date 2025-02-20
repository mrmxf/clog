//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License         https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Check  # check the repo tags for production consistency

package checklegacy

// func checkGitTags(cmd *cobra.Command, args []string) {
// 	// Get the local repository
// 	repo, err := git.PlainOpen(".")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Get the list of tags in the local repository
// 	localTags, err := repo.TagObjects()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Get the list of tags in the origin repository
// 	opts := git.FetchOptions{
// 		RemoteName: "origin",
// 		Progress:   nil,
// 	}
// 	origin, err := repo.Remote("origin")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = origin.Fetch(&opts)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Get the list of tags in the local repository
// 	originTags, err := repo.TagObjects()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("originTags")
// 	fmt.Println(originTags)

// 	// Get the list of tags in the upstream repository
// 	upstream, err := repo.Remote("upstream")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = upstream.Fetch(&opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Get the list of tags in the local repository
// 	upstreamTags, err := repo.TagObjects()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("upstreamTags")
// 	fmt.Println(upstreamTags)

// 	// ... just iterates over the tags
// 	err = localTags.ForEach(func(t *object.Tag) error {
// 		fmt.Printf("    local:%v (%v)\n", t.Name, t.Message)
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	// ... just iterates over the tags
// 	err = originTags.ForEach(func(t *object.Tag) error {
// 		fmt.Printf("    local:%v (%v)\n", t.Name, t.Message)
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	// ... just iterates over the tags
// 	err = upstreamTags.ForEach(func(t *object.Tag) error {
// 		fmt.Printf("    local:%v (%v)\n", t.Name, t.Message)
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
