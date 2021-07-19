package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/cpanato/github-sync/pkg/config"
	"github.com/cpanato/github-sync/pkg/github"
	"github.com/cpanato/github-sync/pkg/reconciler"
)

type options struct {
	dryRun     bool
	config     string
	authConfig string
}

func parseOptions() options {
	var o options
	flag.BoolVar(&o.dryRun, "dry-run", true, "does nothing if true (which is the default)")
	flag.StringVar(&o.config, "config", "", "path to a configuration file, or directory of files")
	flag.StringVar(&o.authConfig, "auth", "auth.json", "path to github auth config")
	flag.Parse()
	return o
}

func main() {
	log.Println("Starting GitHub GitOps reconciler")
	o := parseOptions()

	c, err := github.LoadConfig(o.authConfig)
	if err != nil {
		log.Fatalf("Failed to load github auth config: %v.\n", err)
	}

	client, err := github.NewClient(c.AuthToken)
	if err != nil {
		log.Fatalf("failed to init client: %v\n", err)
	}

	stat, err := os.Stat(o.config)
	if err != nil {
		log.Fatalf("Failed to stat %s: %v\n", o.config, err)
	}
	p := config.NewParser()

	if stat.IsDir() {
		err = p.ParseDir(o.config)
	} else {
		err = p.ParseFile(o.config, path.Dir(o.config))
	}
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	strDict := map[string]bool{}
	for _, value := range p.Config.Users {
		if _, exist := strDict[value.Username]; exist {
			log.Fatalf("Reconciliation failed: Username duplicate in config file. Username must be unique unique: username = %s", value.Username)
		} else {
			strDict[value.Username] = true
		}
	}

	strDictRepo := map[string]bool{}
	for _, repo := range p.Config.Repositories {
		if _, exist := strDictRepo[repo.Name]; exist {
			log.Fatalf("Reconciliation failed: Repository duplicate in config file. Repository must be unique unique: Repository = %s", repo.Name)
		} else {
			strDictRepo[repo.Name] = true
		}

		strDictCollab := map[string]bool{}
		for _, value := range repo.Collaborators {
			if _, exist := strDictCollab[value.Username]; exist {
				log.Fatalf("Reconciliation failed: Username duplicate in config file. Username must be unique unique: username = %s", value.Username)
			} else {
				strDictCollab[value.Username] = true
			}
		}
	}

	r := reconciler.New(client, p.Config, c.Org)
	if err := r.Reconcile(o.dryRun); err != nil {
		log.Fatalf("Reconciliation failed: %v\n", err)
	}
}
