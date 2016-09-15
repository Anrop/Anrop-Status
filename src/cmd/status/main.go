package main

import "status"
import "github.com/sourcegraph/checkup"

func main() {
	c := checkup.Checkup{
		Checkers: []checkup.Checker{
			checkup.HTTPChecker{Name: "anrop.se (HTTP)", URL: "http://anrop.se/start.php", Attempts: 5},
			checkup.HTTPChecker{Name: "www.anrop.se (HTTP)", URL: "http://www.anrop.se/start.php", Attempts: 5},
			checkup.HTTPChecker{Name: "www.anrop.se (HTTPS)", URL: "https://www.anrop.se/start.php", Attempts: 5},
			checkup.HTTPChecker{Name: "arma3.anrop.se (HTTP)", URL: "http://arma3.anrop.se", Attempts: 5},
			//checkup.HTTPChecker{Name: "arma3sync.anrop.se (HTTP)", URL: "http://arma3sync.anrop.se/.a3s/autoconfig ", Attempts: 5},
			checkup.HTTPChecker{Name: "jenkins.anrop.se (HTTP)", URL: "http://jenkins.anrop.se", Attempts: 5},
			status.TeamspeakChecker{Name: "ts.anrop.se (TeamSpeak)", URL: "ts.anrop.se:10011", Attempts: 5},
		},
		Storage: checkup.FS{
			Dir: "status",
		},
	}

	c.CheckAndStore()
}
