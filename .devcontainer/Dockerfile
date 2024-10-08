FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

USER vscode
RUN go install github.com/jackc/tern/v2@latest
RUN sudo apt-get update && sudo apt-get install telnet
USER root

RUN echo '[ ! -z "$TERM" -a -r /etc/motd ] && cat /etc/motd' \
    >> /etc/bash.bashrc \
    ; echo "\
===================================================================\n\
= Welcome to dev containers                                       =\n\
===================================================================\n\
\n\
Here are some useful commands to get you started:\n\
\n\
Purge all data from the database (publications/datasets only):\n\
    go run main.go reset --force\n\
\n\
Seed the app with dummy organizations:\n\
    go run main.go seed organizations < .devcontainer/people/organizations.json\n\
\n\
Seed the app with dummy people:\n\
    go run main.go seed people < .devcontainer/people/people.json\n\
\n\
Seed the app with dummy projects:\n\
    go run main.go seed projects < .devcontainer/projects/projects.json\n\
\n\
Seed the app with dummy candidate records:\n\
    go run main.go seed candidate-records < .devcontainer/candidate_records/candidate_records.json\n\
\n\
Start the app:\n\
    make dev\n\
\n\
Happy hacking!\n"\
    > /etc/motd

# [Optional] Uncomment this line to install global node packages.
RUN su vscode -c "source /usr/local/share/nvm/nvm.sh" 2>&1