.PHONY: help
help: # Show help for each of the Makefile target.
	@grep -E '^[a-zA-Z0-9 _]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

run_python: # Run algorithm using python script
	python3 repositoryRank.py

run_golang: # Run algorithm using golang script
	cd golang && go run main.go