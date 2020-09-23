FROM php:fpm-alpine

COPY wait-for-it.sh /usr/bin/wait-for-it

RUN chmod +x /usr/bin/wait-for-it

RUN apk --update --no-cache add git

RUN set -e -u -x \
\
    # Install packages
\
    && apk add --no-cache --no-progress --virtual BUILD_DEPS ${PHPIZE_DEPS} \
    && apk add --no-cache --no-progress bash ca-certificates git openssh-client shadow

COPY --from=composer /usr/bin/composer /usr/bin/composer

WORKDIR /var/www

CMD composer install;  php-fpm

EXPOSE 9000
