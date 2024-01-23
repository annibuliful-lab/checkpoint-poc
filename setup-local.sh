docker-compose up -d

cp .env.sample .env

pnpm run db-primary:migrate

pnpm run db-primary:seed