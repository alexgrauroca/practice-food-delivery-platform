ADRTOOLS := $(shell command -v adr 2> /dev/null)

adr-new: | install-adr-tools
ifndef dir
	$(error dir is required)
endif

ifndef title
	$(error title is required)
endif
	@# Adding the template into the target folder
	@mkdir -p $(dir)/templates
	@cp templates/template.md $(dir)/templates/template.md

	@# Setting up the adr directory to generate the new file
	@echo "$(dir)" > .adr-dir

	@# Creating the new adr
ifndef super
	@adr new $(title)
else
	@adr new -s $(super) $(title)
endif

	@# Removing the temporary templates folder
	@rm -rf $(dir)/templates

install-adr-tools:
ifeq ("$(wildcard $(ADRTOOLS))", "")
	sudo apt update
	sudo apt install adr-tools
endif