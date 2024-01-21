import { PermissionAction } from '@prisma/client';
import { client } from './client';
const actions = Object.values(PermissionAction);

const subjects = ['project'];

export async function seedPermissions() {
  for (const subject of subjects) {
    for (const action of actions) {
      await client.permission.upsert({
        update: {},
        create: {
          subject,
          action,
        },
        where: {
          subject_action: {
            subject,
            action,
          },
        },
      });
    }
  }
}
