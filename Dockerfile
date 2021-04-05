FROM debian:buster
LABEL maintainer="Asaf Ohayon <asaf@sysbind.co.il>"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN apt update \
    && apt install --no-install-recommends -y \
        apt-utils \
        ca-certificates \
        lsb-release \
        pigz \
        curl \
	gnupg  \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && cat /etc/apt/sources.list.d/pgdg.list \
	&& curl --silent https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
    && apt update \
	&& apt install --no-install-recommends -y  \
        postgresql-client-13

COPY pgmesh /usr/bin/pgmesh

ENTRYPOINT ["/usr/bin/pgmesh"]
