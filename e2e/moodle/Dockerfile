FROM php:7.4
ARG MOODLE_VERSION=3.10.2

RUN apt update && apt install -y ghostscript unoconv libpq-dev libzip-dev libicu-dev libpng-dev libxslt-dev
RUN docker-php-ext-install pgsql zip intl gd xsl xmlrpc soap
RUN docker-php-ext-enable opcache
RUN curl -LON https://github.com/moodle/moodle/archive/refs/tags/v${MOODLE_VERSION}.tar.gz
RUN tar xvzf v3.10.2.tar.gz && mv moodle-${MOODLE_VERSION} /var/moodle
RUN mkdir /var/moodledata
COPY config.php /var/moodle/
WORKDIR /var/moodle
ENTRYPOINT ["php", "-S", "0.0.0.0:80"]
EXPOSE 80
