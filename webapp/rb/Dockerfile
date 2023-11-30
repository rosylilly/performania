FROM ruby:3.2.2

ENV GEM_HOME=/gem
ENV BUNDLE_PATH=/gem
ENV BUNDLE_BIN=/gem/bin
ENV PATH="/gem/bin:${PATH}"

EXPOSE 9292

WORKDIR /app
RUN bundle config set --local deployment 'true'

COPY Gemfile Gemfile.lock /app/
RUN bundle install

COPY . /app

CMD [ "bundle", "exec", "puma", "--bind", "tcp://0.0.0.0:9292" ]
