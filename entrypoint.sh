set -e

echo "Checking if superuser exists..."
# call CLI inside the container
./cli-tool register \
  --login="$INIT_SUPERUSER_LOGIN" \
  --password="$INIT_SUPERUSER_PASSWORD" || echo "Skipping superuser creation"

echo "Starting Protego server..."
exec ./protego
