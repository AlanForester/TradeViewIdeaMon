FROM ruby:2.5.3
RUN apt-get update -qq && apt-get install -y nodejs postgresql-client
WORKDIR /www
COPY Gemfile /www/Gemfile
COPY Gemfile.lock /www/Gemfile.lock
RUN /usr/local/bin/bundle install
COPY . /www

# Add a script to be executed every time the container starts.
COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]
EXPOSE 3000

# Configure the main process to run when running the image
CMD ["rails", "server", "-b", "0.0.0.0"]