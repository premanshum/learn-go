package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func main() {

	// wd, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }
	// fmt.Println("WD:", wd)
	// branch := "main"

	// token := "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"

	// ansibleTempFolder := filepath.Join(wd, "temp")

	// _, err = git.PlainClone(ansibleTempFolder, false, &git.CloneOptions{
	// 	URL: "https://github.com/Maersk-Global/capella-automation.git",
	// 	Auth: &http.BasicAuth{
	// 		Username: "random-string", // anything except an empty string
	// 		Password: token,
	// 	},
	// 	ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
	// })
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }

	// destinationDir := "./ansible"

	// if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
	// 	//fmt.Println("Directory does not exist. Creating directory:", dirName)
	// 	// Create a new directory
	// 	err := os.Mkdir(destinationDir, 0755)
	// 	if err != nil {
	// 		fmt.Println("Error creating directory:", err)
	// 	}
	// 	fmt.Println("Directory created successfully:", destinationDir)
	// }

	// cmd := exec.Command("cp", "--recursive", "./temp/ansible/landingzone/roles", destinationDir)
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }

	//getCommonAnsibleRoles("capella-automation", "/etc/ansible", "main", "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")

}

/*
Copies the common ansible roles, like SSHKeyDistribution, InstallGrafana etc,
from the GitHub repo to /etc/ansible folder in the Management-Agent server
Steps:

 1. Create a temp folder that contain Ansible Common-roles (eg: capella-automation/ansible/common/roles). The folder will be deleted after the job is done.

 2. Clone the GitHub repository to the temp folder.

 3. Create a folder, ansible, in the server, inside `etc`. (/etc/ansible)

 4. Copy the Roles folder from the temp directory to the local directory

    That's all
*/
func getCommonAnsibleRoles(commonAnsibleRolesRepo string, localDirectory string, branch string, githubPAT string) error {

	if len(commonAnsibleRolesRepo) == 0 || len(githubPAT) == 0 || len(localDirectory) == 0 {
		fmt.Println("Error: one of the parameter is missing: CommonAnsibleRolesRepo, localDirectory or githubPAT")
		return fmt.Errorf("Error: one of the parameter is missing: CommonAnsibleRolesRepo, localDirectory or githubPAT")
	}

	if branch == "" {
		branch = "main"
	}

	tempCommonAnsibleRoles := "tempCommonAnsibleRoles"

	// 1. Create a temp folder that contain Ansible Common-roles
	err := CreateTempDir(tempCommonAnsibleRoles)
	if err != nil {
		fmt.Println("Process failed; Error:", err.Error())
		return fmt.Errorf("Process failed; Error:" + err.Error())
	}

	defer TempFolderCleanup(tempCommonAnsibleRoles)

	err = TempFolderCleanup(tempCommonAnsibleRoles)
	if err != nil {
		fmt.Println("Process failed; clean-up temp directory; Error:", err.Error())
		return fmt.Errorf("Process failed; clean-up temp directory; Error:" + err.Error())
	}

	// 2. Clone the GitHub repository to the temp folder.
	commonAnsibleRolesURL := fmt.Sprintf("https://github.com/Maersk-Global/%s", commonAnsibleRolesRepo)

	_, err = git.PlainClone(tempCommonAnsibleRoles, false, &git.CloneOptions{
		URL: commonAnsibleRolesURL,
		Auth: &http.BasicAuth{
			Username: "random-string", // anything except an empty string
			Password: githubPAT,
		},
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
	})
	if err != nil {
		fmt.Println("Process failed; Cloning GitHub repo; Error:", err.Error())
		return fmt.Errorf("Process failed; Cloning GitHub repo; Error:" + err.Error())
	}

	// 3. Create a folder, ansible, in the server, inside `etc`. (/etc/ansible)
	err = CreateTempDir(localDirectory)
	if err != nil {
		fmt.Println("Process failed; Error:", err.Error())
		return fmt.Errorf("Process failed; Error:" + err.Error())
	}

	// 4. Copy the Roles folder from the temp directory to the local directory
	cmd := exec.Command("cp", "--recursive", tempCommonAnsibleRoles+"/ansible/common/roles", localDirectory)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Process failed; Copying roles; Error:", err.Error())
		return fmt.Errorf("Process failed; Copying roles; Error:" + err.Error())
	}

	fmt.Println("Process success; Common Ansible roles are copied to the directory:" + localDirectory)

	return nil
}

func CreateTempDir(dirName string) error {
	//dirName := "secrets" to store windows server secret files which will be used by ansible to pass ssh keys to windows server
	//dirName := "ansible" to download the playbook to this local temp dir

	//check if temp dir exist
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		//fmt.Println("Directory does not exist. Creating directory:", dirName)
		// Create a new directory
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			//fmt.Println("Error creating directory:", err)
			return fmt.Errorf("error creating %s temp directory: %v", dirName, err)
		}
		//fmt.Println("Directory created successfully:", dirName)
		return nil
	}
	return nil
}

func TempFolderCleanup(folderName string) error {

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error cleaning up %s temp directory: %v", folderName, err)
	}

	tempFolder := filepath.Join(wd, folderName)

	err = os.RemoveAll(tempFolder)
	if err != nil {
		return fmt.Errorf("error removing %s temp directory: %v", folderName, err)
	}
	return nil
}
