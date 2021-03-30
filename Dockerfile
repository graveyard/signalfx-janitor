FROM debian

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update && apt-get install -y ca-certificates curl

COPY ./bin/sfncli /usr/bin/sfncli
COPY ./bin/signalfx-janitor /usr/bin/signalfx-janitor
CMD ["sfncli", "--cmd", "/usr/bin/signalfx-janitor", "--activityname", "${_DEPLOY_ENV}--${_APP_NAME}", "--region", "us-west-2", "--cloudwatchregion", "${_POD_REGION}", "--workername", "MAGIC_ECS_TASK_ID"]
