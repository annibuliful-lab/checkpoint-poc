import { seedPermissions } from './permission';
import { seedUsers } from './user';

async function main() {
  await seedPermissions();
  await seedUsers();
}

main();
