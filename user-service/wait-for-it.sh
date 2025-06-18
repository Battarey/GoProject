#!/usr/bin/env sh
# wait-for-it.sh: ждет, пока TCP-сервис не станет доступен
# Источник: https://github.com/vishnubob/wait-for-it

set -e

host="$1"
shift

until nc -z $host; do
  >&2 echo "Ожидание $host..."
  sleep 1
done

>&2 echo "$host доступен, продолжаем"
exec "$@"
