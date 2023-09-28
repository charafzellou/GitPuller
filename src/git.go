package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

func setURL(mode string, pseudo string, token string) string {
	if token == "" {
		fmt.Printf("No Personal Access Token was provided, so only public repositories will be cloned.\n")
		url := fmt.Sprintf("https://api.github.com/%s/%s/repos?type=public&sort=updated&per_page=2", mode, pseudo)
		fmt.Printf("URL for Public Repositories : %s\n", url)
		return url
	} else {
		fmt.Printf("A Personal Access Token was provided, so only private repositories will be cloned.\n")
		url := fmt.Sprintf("https://api.github.com/%s/%s/repos?type=private&sort=updated&per_page=2", mode, pseudo)
		fmt.Printf("URL for Private Repositories : %s\n", url)
		return url
	}
}

func getRepositories(url string, token string) []Repository {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("[getRepositories] http.NewRequest() failed due to the following reason: ", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("[getRepositories] client.Do() failed due to the following reason: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("[getRepositories] HTTP error: %d", resp.StatusCode)
	}

	// Decode the JSON response into a slice of Repository structs
	var repositories []Repository
	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		log.Fatal(err)
	}

	// If no repositories are found, exit the program
	if len(repositories) == 0 {
		log.Fatalf("[getRepositories] No repositories found for %s", url)
	}

	// Number of repositories found for the user/organization
	fmt.Printf("[getRepositories] Number of repositories found: %d\n", len(repositories))
	return repositories
}

func getMostRecentCommit(pseudo string, repo string, token string) *Commit {
	client := &http.Client{}

	fmt.Printf("Getting most recent commit for %s\n", repo)
	branch_url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", pseudo, repo)
	fmt.Printf("Branch URL: %s\n", branch_url)
	req, err := http.NewRequest("GET", branch_url, nil)
	if err != nil {
		log.Fatalf("[getMostRecentCommit] Error creating request for %s: %s", repo, err.Error())
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("[getMostRecentBranch] Error getting branches for %s: %s", repo, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("[getMostRecentBranch] HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[getMostRecentBranch] Error reading response body: %s", err.Error())
	}

	var branches []Branch
	json_err := json.Unmarshal(body, &branches)
	if json_err != nil {
		log.Fatalf("[getMostRecentBranch] Error unmarshalling response body: %s", err.Error())
	}

	var mostRecentCommit *Commit
	var mostRecentDate time.Time
	for _, branch := range branches {
		commit_url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", pseudo, repo, branch.Name)
		fmt.Printf("Commit URL: %s\n", commit_url)
		req, err := http.NewRequest("GET", commit_url, nil)
		if err != nil {
			log.Fatalf("[getMostRecentCommit] Error creating request for %s: %s", branch.Name, err.Error())
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("[getMostRecentCommit] Error getting commit for %s: %s", branch.Name, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("[getMostRecentCommit] HTTP error: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("[getMostRecentCommit] Error reading response body: %s", err.Error())
		}

		var commit Commit
		json.Unmarshal(body, &commit)

		commitDate, err := time.Parse(time.RFC3339, commit.Commit.Author.Date)
		if err != nil {
			log.Fatalf("[getMostRecentCommit] Error parsing commit date: %s", err.Error())
		}

		if mostRecentCommit == nil || commitDate.After(mostRecentDate) {
			mostRecentCommit = &commit
			mostRecentDate = commitDate
		}
	}

	return mostRecentCommit
}

func cloneAndPullRepo(pseudo string, repo Repository, token string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Repository Name: %s\n", repo.Name)

	// Clone the repository
	path := fmt.Sprintf("../assets/%s/%s", pseudo, repo.Name)
	fmt.Printf("%s repository path locally: %s\n", repo.Name, path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		cloneCmd := exec.Command("git", "clone", repo.URL, path)
		fmt.Printf("Executing \"git clone %s %s\" command\n", repo.URL, path)
		if repo.Private {
			cloneCmd.Env = append(cloneCmd.Env, "GIT_ASKPASS=echo")
		}
		if err := cloneCmd.Run(); err != nil {
			log.Fatalf("Error cloning %s: %s", repo.Name, err.Error())
		}
	} else if err != nil {
		log.Fatalf("Error checking if %s exists: %s", path, err.Error())
	}

	// Execute a "git pull --all" command in the repository's directory
	fmt.Printf("Executing \"git pull --all\" command in %s\n", repo.Name)
	pullCmd := exec.Command("git", "pull", "--all")
	pullCmd.Dir = path
	if err := pullCmd.Run(); err != nil {
		log.Fatalf("Error pulling changes for %s: %s", repo.Name, err.Error())
	}

	// Get the most recent commit
	mostRecentCommit := getMostRecentCommit(pseudo, repo.Name, token).SHA
	fmt.Printf("Most recent commit for %s: %s\n", repo.Name, mostRecentCommit)
	commitCmd := exec.Command("git", "checkout", mostRecentCommit)
	commitCmd.Dir = path
	if err := commitCmd.Run(); err != nil {
		log.Fatalf("Error checking out commit %s for %s: %s", mostRecentCommit, repo.Name, err.Error())
	}

	fmt.Printf("Pulled all changes for %s\n", repo.Name)
	Separator()
}
