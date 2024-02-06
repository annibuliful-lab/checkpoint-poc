import { generate } from '@genql/cli';
import { readFileSync } from 'fs';
import path from 'path';

async function main() {
  const typeDefs = readFileSync('./tools/genql/generated.graphql', {
    encoding: 'utf-8',
  });

  await generate({
    schema: typeDefs,
    output: path.join(
      __dirname,
      '../../backend-testing/graphql/generated'
    ),
  });
}

main();
