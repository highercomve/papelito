FROM fluent/fluentd:v1.7-1 as base

# Use root account to use apk
USER root

# below RUN includes plugin as examples elasticsearch is not required
# you may customize including plugins as you wish
RUN apk add --no-cache --update --virtual .build-deps \
        sudo build-base ruby-dev \
 && sudo gem install fluent-plugin-elasticsearch \
 && sudo gem install fluent-plugin-kafka \
 && sudo gem sources --clear-all \
 && apk del .build-deps \
 && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

COPY fluent.conf /fluentd/etc/

USER fluent


FROM base

USER root

RUN apk add --no-cache --update --virtual .build-deps \
    sudo build-base ruby-dev \
    # cutomize following instruction as you wish
    && sudo gem install fluent-plugin-elasticsearch --no-document \
    && sudo gem sources --clear-all \
    && apk del .build-deps \
    && rm -rf /home/fluent/.gem/ruby/2.5.0/cache/*.gem
